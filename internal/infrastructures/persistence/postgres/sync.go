package postgres

import (
	"errors"
	"fmt"
	"strings"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// postgresTableCoalesceMaxID はホワイトリスト上のテーブルで COALESCE(MAX(id), 0) を返す。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - table string: テーブル名（memos, tag_managers 等）
//
// return:
//   - int64: 最大 ID（行なしは 0）
//   - error: 未対応テーブル・SQL エラー
func postgresTableCoalesceMaxID(tx *gorm.DB, table string) (int64, error) {
	switch table {
	case "memos", "tag_managers", "related_questions", "answers", "refer_managers":
	default:
		return 0, fmt.Errorf("postgresTableCoalesceMaxID: unsupported table %q", table)
	}
	var max int64
	q := fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)
	if err := tx.Raw(q).Scan(&max).Error; err != nil {
		return 0, err
	}
	return max, nil
}

// assignBulkInsertZeros は一括 INSERT 前に ID==0 の行へ連番の主キーを割り当てる。
// 混在 INSERT による *_pkey 重複を防ぐ。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - rows []T: 挿入行
//   - table string: テーブル名
//   - id func(*T) *int64: 行の ID フィールド参照
//
// return:
//   - error: ID 採番・DB エラー
func assignBulkInsertZeros[T any](tx *gorm.DB, rows []T, table string, id func(*T) *int64) error {
	hasZero := false
	batchMax := int64(0)
	for i := range rows {
		v := *id(&rows[i])
		if v == 0 {
			hasZero = true
		} else if v > batchMax {
			batchMax = v
		}
	}
	if !hasZero {
		return nil
	}
	dbMax, err := postgresTableCoalesceMaxID(tx, table)
	if err != nil {
		return err
	}
	next := batchMax
	if dbMax > next {
		next = dbMax
	}
	for i := range rows {
		if *id(&rows[i]) == 0 {
			next++
			*id(&rows[i]) = next
		}
	}
	return nil
}

// syncChildrenByKey は親 FK 配下の子行をキーで突き合わせ INSERT / UPDATE / DELETE する。
// want は in-place で主キーが埋め戻される。softDelete=true は論理削除、false は物理削除。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - table string: 子テーブル名
//   - parentColumn string: 親 FK 列名
//   - parentID int64: 親 ID
//   - want []T: 望ましい子行の集合
//   - keyFn func(*T) K: 自然キー抽出
//   - pkFn func(*T) *int64: 主キー参照
//   - applyUpdate func(tx *gorm.DB, prev T, next *T) error: 更新時コールバック（nil 可）
//   - softDelete bool: 論理削除かどうか
//
// return:
//   - []int64: 削除した子行の ID
//   - error: DB エラー
func syncChildrenByKey[T any, K comparable](
	tx *gorm.DB,
	table string,
	parentColumn string,
	parentID int64,
	want []T,
	keyFn func(*T) K,
	pkFn func(*T) *int64,
	applyUpdate func(tx *gorm.DB, prev T, next *T) error,
	softDelete bool,
) (deletedIDs []int64, err error) {
	var existing []T
	if err := tx.Where(parentColumn+" = ?", parentID).Find(&existing).Error; err != nil {
		return nil, err
	}

	var zeroK K
	byKey := make(map[K]T, len(existing))
	for _, row := range existing {
		byKey[keyFn(&row)] = row
	}

	var insertIdx []int
	for i := range want {
		k := keyFn(&want[i])
		if k == zeroK {
			insertIdx = append(insertIdx, i)
			continue
		}
		prev, ok := byKey[k]
		if !ok {
			insertIdx = append(insertIdx, i)
			continue
		}
		*pkFn(&want[i]) = *pkFn(&prev)
		if applyUpdate != nil {
			if err := applyUpdate(tx, prev, &want[i]); err != nil {
				return nil, err
			}
		}
		delete(byKey, k)
	}

	var toDelete []int64
	for _, row := range byKey {
		toDelete = append(toDelete, *pkFn(&row))
	}
	if len(toDelete) > 0 {
		var model T
		if softDelete {
			if err := tx.Where("id IN ?", toDelete).Delete(&model).Error; err != nil {
				return nil, err
			}
		} else {
			if err := tx.Unscoped().Where("id IN ?", toDelete).Delete(&model).Error; err != nil {
				return nil, err
			}
		}
		deletedIDs = append(deletedIDs, toDelete...)
	}

	if len(insertIdx) == 0 {
		return deletedIDs, nil
	}

	toInsert := make([]T, len(insertIdx))
	for j, i := range insertIdx {
		toInsert[j] = want[i]
	}
	if err := assignBulkInsertZeros(tx, toInsert, table, pkFn); err != nil {
		return deletedIDs, err
	}
	if err := tx.Create(&toInsert).Error; err != nil {
		return deletedIDs, err
	}
	for j, i := range insertIdx {
		*pkFn(&want[i]) = *pkFn(&toInsert[j])
	}
	return deletedIDs, nil
}

// resolveReferIDForSync は ReferManager から参照資料 ID を解決する。
// ReferID / Refer.ID が無い場合は title+url で既存行を検索し、なければ新規作成する。
func resolveReferIDForSync(tx *gorm.DB, rm entity.ReferManager) (int64, error) {
	if rm.Refer.ID != 0 {
		return rm.Refer.ID, nil
	}
	title := strings.TrimSpace(rm.Refer.Title)
	url := strings.TrimSpace(rm.Refer.URL)
	if title == "" || url == "" {
		return 0, nil
	}
	var existing entity.Refer
	err := tx.Where("title = ? AND url = ?", title, url).First(&existing).Error
	if err == nil {
		return existing.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	newRefer := entity.Refer{Title: title, URL: url}
	if err := tx.Create(&newRefer).Error; err != nil {
		return 0, err
	}
	return newRefer.ID, nil
}
