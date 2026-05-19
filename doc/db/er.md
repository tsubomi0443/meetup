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

    ROLE {
        number id PK
        string name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    ROLE ||--o{ USER : "ユーザーは一つの権限を持つ"

    QUESTION {
        number id PK
        number origin_question_id FK
        number support_id FK
        string talkroom_id "LINEWORKSのトークルームID"
        string title
        text content
        datetime due
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    ANSWER {
        number id PK
        number user_id FK
        number question_id FK
        text content
        boolean is_final
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    QUESTION ||--o{ ANSWER : "質問は複数の回答を持つ"
    USER ||--o{ ANSWER : "回答は一つのユーザーによって作成される"

    SUPPORT {
        number id PK
        number user_id FK
        number support_status_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    SUPPORT_STATUS {
        number id PK
        string name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    USER ||--o{ SUPPORT : "サポートは一つのユーザーに紐づく"
    SUPPORT_STATUS ||--o{ SUPPORT : "サポートは一つのステータスを持つ"
    QUESTION ||--o| SUPPORT : "質問は一つのサポートに紐づく（任意）"

    MEMO {
        number id PK
        number question_id FK
        number user_id FK
        text content
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    USER ||--o{ MEMO : "メモは一つのユーザーによって作成される"
    QUESTION ||--o{ MEMO : "質問には複数のメモが関連付けられる"

    CATEGORY {
        number id PK
        string name
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

    CATEGORY ||--o{ TAG : "タグは一つのカテゴリに属する"

    TAG_MANAGER {
        number id PK
        number tag_id FK
        number question_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    TAG ||--o{ TAG_MANAGER : "タグは複数のタグ管理に紐づく"
    QUESTION ||--o{ TAG_MANAGER : "質問は複数のタグを持つことができる"

    REFER {
        number id PK
        string title
        string url
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

    ANSWER ||--o{ REFER_MANAGER : "回答は複数の参照情報を持つことができる"
    REFER ||--o{ REFER_MANAGER : "参照は複数の参照管理に紐づく"

    ESCALATION {
        number id PK
        number from_question_id FK
        number to_question_id FK
        datetime escalated_at
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    QUESTION ||--o{ ESCALATION : "エスカレーション元"
    QUESTION ||--o{ ESCALATION : "エスカレーション先"

    NOTICE_TYPE {
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
        datetime display_due
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    NOTICE_TYPE ||--o{ NOTICE : "通知は一つの通知タイプを持つ"
    QUESTION ||--o{ NOTICE : "期限通知などで質問に紐づく場合がある"

    RELATED_QUESTION {
        number id PK
        number question_id FK
        number related_question_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    QUESTION ||--o{ RELATED_QUESTION : "関連質問（question_id 側）"
    QUESTION ||--o{ RELATED_QUESTION : "関連質問（related_question_id 側）"

    SENDER {
        number id PK
        string uid "LINEWORKSのユーザID（一意）"
        string name
        string department_name
    }

    SENDER_TALK {
        number id PK
        number sender_id FK
        number question_id FK
        string talkroom_id "LINEWORKSのトークルームID"
        text content
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    SENDER ||--o{ SENDER_TALK : "送信者のトーク履歴"
    QUESTION ||--o{ SENDER_TALK : "質問に紐づく外部チャット"
