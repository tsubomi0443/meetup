# 修正・機能追加ガイドライン

## 基本方針

- 実装事実に合わせて最小変更で直す
- 既存の `handler` と `_mac_infrastructure` の責務分割を崩さない
- UI の都合で無理に DB 層をゆがめない。必要な補完はサーバ側で行う
- ドキュメントを更新する場合は `doc/` 配下へ追加し、必要なら `doc/diff_YYYYMMDD.md` に変更内容を残す

## 配置ルール

- HTTP ルートとレスポンス制御: [handler/](../handler/)
- Form と Entity の変換: [_mac_infrastructure/converter.go](../_mac_infrastructure/converter.go)
- DB アクセス、トランザクション、関連同期: [_mac_infrastructure/db.go](../_mac_infrastructure/db.go)
- GORM モデル: [_mac_infrastructure/entity.go](../_mac_infrastructure/entity.go)
- UI/API 用 JSON 形: [_mac_infrastructure/form.go](../_mac_infrastructure/form.go)

## GORM 方針

- 新規コードは `gorm.G[T](db)` を優先する
- 単一レコード更新は [`UpdateByID`](../_mac_infrastructure/db.go) を優先する
- `Model(&m).Updates(&m)` のような伝統 API は、WHERE 条件の暗黙推論や関連名 `Omit` の扱いで不安定になりやすいため避ける
- 取得系は必要な `Preload` だけを明示する

## 関連更新の方針

質問更新のように親子をまとめて扱う処理は、単純な `UpdateByID` ではなく専用トランザクション関数を作るか、既存の [`UpdateQuestionInTransaction`](../_mac_infrastructure/db.go) を拡張する。

### NOT NULL の子テーブル

次のようなテーブルでは `Association.Replace` / `Association.Clear` を使わない。

- `tag_managers`
- `memos`
- `refer_managers`

理由:

- has-many の `Replace/Clear` は既存子レコードの FK を `NULL` にしようとする
- 上記テーブルの FK は NOT NULL のため、`23502` を起こしやすい

推奨パターン:

1. サーバ側で FK を補完する
2. 空行をスキップする
3. `WHERE parent_id = ?` で既存を DELETE
4. 残った行だけ `Create` で INSERT

## サーバ側補完

- Form から来る FK 文字列は converter で `int64` に変換する
- 子要素の `QuestionID`, `AnswerID` などは、親 ID が確定しているならサーバ側で上書きする
- `createdAt` が空なら `time.Now()` を使う
- `ID == 0` や空文字だけの要素は保存しない

例:

- `tags[].id == 0` は無視
- `memos[].userId == 0` または空コンテンツは無視
- `refers[].id == 0` や `refer_managers[].referId == 0` は無視対象にする

## Form と Entity の扱い

- UI から受け取る JSON の形は `Form` に閉じ込める
- DB 永続化の都合は `Entity` と DB 層に閉じ込める
- 画面が未対応でも、サーバ側の Entity 保存を先に整備するのは可
- ただし、不要な Form 拡張や無関係な converter 変更は避ける

## ハンドラ追加手順

1. `handler/*.go` にハンドラ関数を追加する
2. 必要な middleware を決める
3. [handler/handlerManager.go](../handler/handlerManager.go) の登録フローに組み込む
4. DB 読み書きが必要なら `_mac_infrastructure/db.go` に関数を追加する
5. 入出力型が必要なら `form.go` と `converter.go` を更新する

## ドキュメント更新の目安

次の変更では `doc/` も合わせて更新する。

- エンドポイント追加・削除
- Form/Entity の構造変更
- 更新パターンの変更
- 開発手順や起動手順の変更

最終更新: 自動生成 2026-04-24
