# API 仕様

このファイルは、現在 `SetupHandlers()` から実際に登録されているルートのみを対象にしています。

## 認証

### `POST /login`

- 実装: [handler/authHandler.go](../handler/authHandler.go)
- 認証成功時、JWT を Cookie `access_token` に保存
- 成功レスポンス:
  - `200 OK`
  - JSON: `{ "redirect": "/mock/5" }`
- 失敗レスポンス:
  - `401 Unauthorized` または `500 Internal Server Error`

リクエスト JSON:

```json
{
  "email": "user@example.com",
  "pass": "plain-text-password"
}
```

## ページ

- `GET /login`
  - `login.html` を描画
- `GET /`
  - `index.html` を描画
- `GET /mock/:id`
  - `mock{id}.html` を描画
  - JWT 必須

実装: [handler/pageHandler.go](../handler/pageHandler.go)

## SSE

### `GET /sse`

- 実装: [handler/sseHandler.go](../handler/sseHandler.go)
- `text/event-stream` を返す
- 主なイベント:
  - `time-tick`
  - `question`
  - `user`
  - `tag`
  - `notice`
  - `error`

配信元の実装: [handler/eventHub.go](../handler/eventHub.go)

## REST API

REST API はすべて `/api/v1` 配下で、JWT Cookie 認証が必要です。

### User

- `POST /api/v1/user`
  - 新規ユーザー登録
  - 実装: `registerUser`
- `PUT /api/v1/user`
  - ユーザー更新
  - 実装: `updateUserByID`
- `DELETE /api/v1/user/:id`
  - ユーザー削除
  - 実装: `deleteUserByID`

`UserForm` の主な入力:

- `id`: 数値
- `name`: 文字列
- `email`: 文字列
- `roleId`: 文字列
- `password`: 文字列
- `role`: 任意

レスポンス:

- `POST`: 空ボディに近い成功レスポンス
- `PUT`: `200 OK` + `null`
- `DELETE`: `200 OK` + `null`

### Question

- `POST /api/v1/question`
  - 質問登録
  - 実装: `registerQuestion`
- `GET /api/v1/question/:id`
  - 質問取得
  - 実装: `getQuestionByID`
- `PUT /api/v1/question`
  - 質問更新
  - 実装: `updateQuestionByID`
- `DELETE /api/v1/question/:id`
  - 質問削除
  - 実装: `deleteQuestionByID`

`QuestionForm` の主な入力:

- `id`: 数値
- `originQuestionId`: 文字列ポインタ
- `answerId`: 数値ポインタ
- `supportId`: 数値ポインタ
- `title`: 文字列
- `content`: 文字列
- `deleted`: 真偽値
- `due`: ISO 日時文字列ポインタ
- `createdAt`: ISO 日時文字列ポインタ
- `support`: `SupportForm`
- `answer`: `AnswerForm`
- `memos`: `MemoForm[]`
- `tags`: `TagForm[]`
- `notices`: `NoticeForm[]`

`AnswerForm` の主な入力:

- `id`
- `userId`
- `content`
- `answeredAt`
- `createdAt`
- `refers`: `ReferForm[]`

`MemoForm` の主な入力:

- `id`
- `questionId`
- `userId`
- `content`

レスポンス:

- `POST`: 登録成功時は DB 登録のみで JSON を返さない
- `GET`: `QuestionForm`
- `PUT`: `200 OK` + `null`
- `DELETE`: `200 OK` + `""`

### Tag

- `POST /api/v1/tag`
  - タグ登録
  - 実装: `registerTag`
- `PUT /api/v1/tag`
  - タグ更新
  - 実装: `updateTag`
- `DELETE /api/v1/tag/:id`
  - タグ削除
  - 実装: `deleteTagByID`

`TagForm` の主な入力:

- `id`
- `title`
- `usage`
- `categoryId`
- `category`
- `questions`

レスポンス:

- `POST`: DB 登録のみ
- `PUT`: `200 OK` + `null`
- `DELETE`: `200 OK` + `null`

## 変換とバリデーション上の注意

- Form 型は [_mac_infrastructure/form.go](../_mac_infrastructure/form.go) にあり、外部 API の JSON 形です
- 多くの FK は JSON 上は文字列で受け取り、[_mac_infrastructure/converter.go](../_mac_infrastructure/converter.go) で `int64` に変換します
- `createdAt` や `due` などの日時は ISO 文字列から `time.Time` に変換します
- `QuestionToEntity` では `createdAt` が空のとき `time.Now()` を補完します
- 子要素の空データはサーバ側で落とすことがあります
  - 例: `tags[].id == 0` は保存対象にしない
  - 例: `memos[].userId == 0` または空文字コンテンツは保存対象にしない

## 実装上の補足

- `getUsers()` は関数として存在しますが、現時点ではルート登録されていません
- 通知取得用の DB 関数は存在しますが、通知専用の REST ルートは現在登録されていません

最終更新: 自動生成 2026-04-24
