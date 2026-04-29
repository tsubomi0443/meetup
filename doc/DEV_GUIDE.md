# 開発ガイド

## 前提

- Go
- Docker / Docker Compose
- Bun
- PostgreSQL を利用できる環境

## 環境変数

[.env.example](../.env.example) を元に `.env` を用意します。

主な変数:

- `PORT`
- `PSQL_HOST`
- `PSQL_PORT`
- `PSQL_USER`
- `PSQL_PASSWORD`
- `PSQL_DBNAME`
- `PSQL_SSLMODE`
- `JWK_KEY`

DB 接続文字列は [env/env.go](../env/env.go) の `GetDSN()` で構築されます。

## 起動方法

### Docker Compose

[docker-compose.yml](../docker-compose.yml) では次を起動します。

- `app`
  - Go アプリ本体
  - `command: ["air"]`
- `postgresql`
  - PostgreSQL
- `dynamodb`
  - DynamoDB Local

基本的な起動:

```bash
docker compose up --build
```

アプリコンテナは `.env` を読み込み、ホスト側の `${PORT}` をコンテナ内 `1323` に公開します。

### Air によるホットリロード

[.air.toml](../.air.toml) では、Go の再ビルド前に次を実行します。

```bash
bun run --cwd frontend build
```

その後、Go バイナリを `./tmp/main` にビルドして起動します。

## フロントエンドビルド

フロント資産の元は [frontend/](../frontend/) にあります。依存管理とスクリプト実行は Bun 前提です。

主要スクリプト:

- `bun install`
- `bun run --cwd frontend build`
- `bun run --cwd frontend dev`

[frontend/package.json](../frontend/package.json) の `build` は以下を行います。

- Tailwind CSS を `frontend/global.css` から `static/css/output.css` に出力
- vendor ライブラリを `static/vendor/` へコピー

使われる主なフロントライブラリ:

- `htmx.org`
- `htmx-ext-sse`
- `alpinejs`
- `lucide`
- `es-toolkit`

## テンプレートと静的配信

- テンプレート: [templates/](../templates/)
- 静的ファイル: [static/](../static/)

[main.go](../main.go) で `templates/**/*.html` をまとめて読み込みます。静的ファイルは [handler/pageHandler.go](../handler/pageHandler.go) で `/static` にマウントされています。

## DB

- 初期化 SQL: [doc/db/INIT.sql](./db/INIT.sql)
- ER 図: [doc/db/er.md](./db/er.md)

テーブル定義の主な特徴:

- `questions` を中心に `answers`, `supports`, `memos`, `tag_managers`, `notices` がぶら下がる
- `refer_managers`, `tag_managers`, `memos` には NOT NULL の FK がある
- `notices.question_id` は nullable

## 認証

認証は [handler/authHandler.go](../handler/authHandler.go) の JWT Cookie ベースです。

- Cookie 名: `access_token`
- 未認証時、JWT ミドルウェアは `/login` へリダイレクト
- `Secure` 属性は `MODE` によって切り替わる

## ログとエラー

- Echo 側で `RequestLogger` と `Recover` を利用
- アプリ起動失敗は `slog.Error`
- 一部の DB 参照は `logger.Silent` を使って不要ログを抑制

## 参考

- 全体像: [OVERVIEW.md](./OVERVIEW.md)
- 設計: [ARCHITECTURE.md](./ARCHITECTURE.md)
- API: [API.md](./API.md)
- 実装規約: [GUIDELINES.md](./GUIDELINES.md)

最終更新: 自動生成 2026-04-24
