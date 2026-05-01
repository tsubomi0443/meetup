erDiagram
    USER {
        number id PK
        string name
        string password
        string email
        string memo
        number role_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    USER ||--|| ROLE: "ユーザーは一つの権限を持つ (Admin>Manager>Staff>Emproyee)"

    ROLE {
        number id PK
        string role_name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    QUESTION {
        number id PK
        number origin_question_id FK "元の質問ID (エスカレーション元の場合)"
        number answer_id FK "紐づく回答 (最大1件)"
        number support_id FK
        string title
        string content
        datetime due
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    QUESTION ||--o| ANSWER: "質問は一つの回答を持つ (0または1)"
    QUESTION ||--o| SUPPORT: "質問は一つのサポート情報に紐づく"
    QUESTION ||--o{ ESCALATION : "エスカレーション元 (from_question_id)"
    QUESTION ||--o{ ESCALATION : "エスカレーション先 (to_question_id)"
    QUESTION ||--o{ TAG_MANAGER: "質問は複数のタグを持つことができる"
    QUESTION ||--o{ MEMO: "質問には複数のメモが関連付けられる"
    QUESTION ||--o{ RELATED_QUESTION: "１つの質問に対して、複数の関連する質問を持つ（question_id 側）"
    QUESTION ||--o{ RELATED_QUESTION: "１つの質問に対して、複数の関連する質問を持つ（related_question_id 側）"

    ANSWER {
        number id PK
        number user_id FK
        string content
        datetime answered_at
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    ANSWER ||--|| USER: "回答は一つのユーザーによって作成される"
    ANSWER ||--o{ REFER_MANAGER: "回答は複数の参照情報を持つことができる"

    SUPPORT {
        number id PK
        number user_id FK
        number support_status_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    SUPPORT ||--|| USER: "サポートは一つのユーザーに紐づく"
    SUPPORT ||--|| SUPPORT_STATUS: "サポートは一つのステータスを持つ"

    SUPPORT_STATUS {
        number id PK
        string title
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    MEMO {
        number id PK
        number question_id FK
        number user_id FK
        string content
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    MEMO ||--|| USER: "メモは一つのユーザーによって作成される"

    REFER {
        number id PK
        string title
        string url
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    REFER ||--o{ REFER_MANAGER: "参照は複数の参照管理情報に紐づく"

    ESCALATION {
        number id PK
        number from_question_id FK "エスカレーション元の質問ID"
        number to_question_id FK "エスカレーション先の質問ID"
        datetime escalated_at
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    REFER_MANAGER {
        number id PK
        number answer_id FK
        number refer_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    TAG {
        number id PK
        string name
        number usage
        number category_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    TAG ||--|| CATEGORY: "タグは一つのカテゴリに属する"
    TAG ||--o{ TAG_MANAGER: "タグは複数のタグ管理情報に紐づく"

    TAG_MANAGER {
        number id PK
        number tag_id FK
        number question_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    CATEGORY {
        number id PK
        string name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    NOTICE {
        number id PK
        number type_id FK
        number question_id FK
        string content
        datetime displayDue
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
    NOTICE ||--|| NOTICE_TYPE: "通知は必ず１つの通知タイプを持つ"
    NOTICE ||--o| QUESTION: "質問期限通知の場合は質問を持っている場合もある"

    NOTICE_TYPE {
        number id PK
        string name "SYSTEM|ALERT|QUESTION"
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    RELATED_QUESTION {
        number id PK
        number question_id FK
        number related_question_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }
