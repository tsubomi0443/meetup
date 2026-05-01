# `PUT /api/v1/question` メモ連投の検証手順

## 背景

質問更新では `tag_managers` / `memos` / `related_questions` を **DELETE + 一括 INSERT** している。一方でクライアントは既存行に **明示 `id`**、新規メモなどは **`id: 0`** を送り、GORM は1つの INSERT 文に **明示 `id`** と **`DEFAULT`（シーケンス）** が混ざった VALUES を並べがちである。

PostgreSQL でこのとき `nextval` が返した値が、**同じ INSERT バッチ内**の明示 `id` と一致すると `duplicate key ... pkey` になる（インサート自体が失敗し、それ以前の **`INSERT 後に setval` だけでは防げない**）。

## 対策（実装）

`_mac_infrastructure/db.go` の `UpdateQuestionInTransaction` で、各バルク INSERT の **`tx.Create` の直前に**、`id == 0` の行へ

`max(テーブル全体の MAX(id), バッチ内の最大明示 id) + 1` からの連番

を代入する（`assignMemoBulkInsertZeros` / `assignTagManagerBulkInsertZeros` / `assignRelatedQuestionBulkInsertZeros`）。これにより一括 INSERT から `DEFAULT` が消え、上記競合が起きない。

## 手動確認（アプリ）

1. 任意の質問の詳細を開く。
2. メモ欄に文字を入れ、「メモ登録」を **連続 5 回以上** クリックする。
3. いずれも 500 にならず、コンソールに `duplicate key ... memos_pkey` が出ないこと。
4. 失敗時は従来どおり `pop()` され、アラート後に下書きは残ること（再クリックで再送可能）。

## 回帰の観点

- タグの付け外し、関連質問の変更も同じ DELETE + バルク INSERT のため、ID 事前割り当てを同様に適用している。
