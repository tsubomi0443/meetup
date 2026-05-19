package postgres

import "gorm.io/gorm"

// commonPreloadBuilder は Preload 時に id 昇順で並べるビルダを返す。
//
// return:
//   - func(db gorm.PreloadBuilder) error: Preload コールバック
func commonPreloadBuilder() func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		db.Order("id")
		return nil
	}
}

// userPreloadBuilder はユーザー Preload 用の Select・Order を返す。
//
// args:
//   - includePassword bool: password 列を含めるか
//
// return:
//   - func(db gorm.PreloadBuilder) error: Preload コールバック
func userPreloadBuilder(includePassword bool) func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		if includePassword {
			db.Select("id, name, email, password, memo, role_id")
		} else {
			db.Select("id, name, email, memo, role_id")
		}
		db.Order("id ASC")
		return nil
	}
}
