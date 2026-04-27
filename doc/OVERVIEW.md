# アプリ概要

このアプリは、社内向けの質問管理を中心にした Web アプリです。質問を登録し、回答・サポート担当・タグ・メモ・通知と紐づけて管理します。画面は `html/template` で配信し、クライアント側では `htmx` / `Alpine.js` / SSE を使って更新を受け取ります。

## 主なユースケース

- 質問の登録、取得、更新、削除
- ユーザーの登録、更新、削除
- タグの登録、更新、削除
- 回答、サポート、メモ、通知の質問への関連付け
- 回答期限が近い質問に対する通知の自動生成
- SSE による質問、ユーザー、タグ、通知の定期配信

## 技術スタック

- バックエンド: Go, [Echo v5](../main.go), [GORM](../_mac_infrastructure/db.go), PostgreSQL
- 認証: JWT Cookie ベースの簡易認証（[handler/authHandler.go](../handler/authHandler.go)）
- フロントエンド: Bun, Tailwind CSS, `htmx`, `Alpine.js`, `htmx-ext-sse`（[frontend/package.json](../frontend/package.json)）
- テンプレート: `html/template`（[main.go](../main.go)）
- 開発起動: `air` + Docker Compose（[.air.toml](../.air.toml), [docker-compose.yml](../docker-compose.yml)）

## ディレクトリ入口

- [main.go](../main.go): アプリ起動、DB 接続、Echo 初期化、ハンドラ登録
- [handler/](../handler/): 認証、ページ、SSE、API、通知ポーリング
- [_mac_infrastructure/](../_mac_infrastructure/): Entity、Form、Converter、DB アクセス
- [templates/](../templates/): サーバレンダリング用テンプレート
- [static/](../static/): 配信される CSS / JS / vendor 資産
- [frontend/](../frontend/): Tailwind と vendor ファイルのビルド元
- [env/](../env/): 環境変数の読み出し
- [doc/db/INIT.sql](./db/INIT.sql): DB 初期化 SQL
- [doc/db/er.md](./db/er.md): ER 図

## 併読推奨

- 設計とデータフロー: [ARCHITECTURE.md](./ARCHITECTURE.md)
- API と入出力: [API.md](./API.md)
- 開発手順: [DEV_GUIDE.md](./DEV_GUIDE.md)
- 修正・追加時の規約: [GUIDELINES.md](./GUIDELINES.md)

最終更新: 自動生成 2026-04-24
