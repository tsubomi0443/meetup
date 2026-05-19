```mermaid
erDiagram
%% 全テーブルはcreated_at, updated_at, deleted_atを持つものとする

%% ユーザデータの格納先（SAML認証の予定）
USERS {
    number id PK
    string name
    string email
    string password

    number department_id FK
    %% 権限ID、Ownerは全てのログへのアクセス、BossはDepartment以下またはBossRelations参照、GeneralはPeeringのリレーション参照
    number authority_id FK
}
DEPARTMENTS ||--|{ USERS: "ユーザは必ずどこかの部署に所属する"
AUTHORITIES ||--o{ USERS: ""

%% MentorユーザとMenteeユーザの紐づけテーブル
PEERINGS {
    number id PK
    number mentor_id FK
    number mentee_id FK
}
USERS ||--o{ PEERINGS: "Mentor"
USERS ||--|| PEERINGS: "Mentee"

%% 部署名を設定するためのテーブル
DEPARTMENTS {
    number id PK
    string name
    number boss_id FK
}
USERS ||--o{ DEPARTMENTS: "部署は部長を１人を設定できる。"

%% チーム単位の関係を作成するテーブル
%% SeniorManagerの下に複数のMemberユーザが所属する
TEAMS {
    number id PK
    string name "チーム名"
    number leader_id FK "管理者(部長など)のユーザID"
    number department_id FK "部署ID"
}
USERS ||--o{ TEAMS: "チームにはリーダーを１人設定できる"
DEPARTMENTS ||--o{ TEAMS: "部署には０以上のチームが存在する"

%% チームメンバーを格納するテーブル
TEAM_MEMBERS {
    number id PK
    number team_id FK
    number user_id FK
}
USERS ||--o{ TEAM_MEMBERS: "管理者もメンバーの１人として扱う想定のため必ず１人以上となる"
TEAMS ||--|{ TEAM_MEMBERS: "チームには１人以上が所属する（member_id）"

%% 権限テーブル。機能制限などに用いる。ユーザの回答の確認範囲などの制御。
AUTHORITIES {
    number id PK
    number level "Owner=100, Boss=50, General=10"
    string name "Owner | Boss | General"
}

%% 1on1のテンプレートを表すテーブル
TEMPLATES {
    number id PK
    string title "テンプレート名"
    text description "テンプレート説明"
}

%% テンプレートに記載されるトピックを格納するテーブル
TOPICS {
    number id PK
    number topic_number "トピック番号 1からの連番"
    text content "トピック内容"
    number template_id FK "テンプレートID"
}
TEMPLATES ||--|{ TOPICS: ""

%% 1on1ミーティングを現すテーブル
MEETINGS {
    number id PK
    number template_id FK "使用するテンプレートのID"
    number peering_id FK "ペアID"
    number status "面談実施状況。0(未実施)|10(実施済み)"

    datetime execute_schedule "面談予定日"
    datetime executed_at "面談実施日"
}
TEMPLATES ||--|{ MEETINGS: "テンプレートを用いて複数のミーティングが設定される"
PEERINGS ||--o{ MEETINGS: "ペア１つに大して複数のミーティングが存在する（月次のため）"

%% 1on1ミーティングに利用するトピックのテーブル。ミーティングとトピックを紐づける中間テーブル
MEETING_TOPICS {
    number id PK
    number meeting_id FK
    number topic_id FK
}
MEETINGS ||--|{ MEETING_TOPICS: "ミーティングには１つ以上のトピックが設定される"
TOPICS ||--o{ MEETING_TOPICS: "トピック１つに対して、ミーティングはユーザ単位で作成されるためHasMany"

%% トピックに対して（主に）Menteeが記入した内容を格納するテーブル
ANSWERS {
    number id PK
    text content
    datetime submitted_at

    number author_id FK "記入者のユーザID"
    number meeting_topic_id FK
}
USERS ||--o{ ANSWERS: "回答にはユーザ１人が設定される"
MEETING_TOPICS ||--o| ANSWERS: "トピックにはメンティー１人の回答が登録される（初期は無回答のため０も許容）"

%% コメント格納テーブル。記入済みトピックに対してMentorが面談後にコメントを残し、ミーティング自体にも複数のユーザがコメントを残すことができる。
COMMENTS {
    number id PK
    text content

    number comment_type_id FK
    number answer_id FK null "meeting_idとは排他関係"
    number meeting_id FK    null "answer_idとは排他関係"
    number author_id FK "コメント記入者のユーザID。Mentor以外の参照権限を持つユーザが記入できるため、こちらにも記入者IDを持つ"
}
COMMENT_TYPES ||--o{ COMMENTS: ""
ANSWERS ||--o| COMMENTS: "feedback, 回答に対して０か１つのコメントが設定される（メンターからのコメントが回答１につき０，１で設定）。"
MEETINGS ||--o{ COMMENTS: "summary, ミーティング結果にフィードバックのコメントが設定できる"
USERS ||--o{ COMMENTS: "記入者ユーザID。ユーザは複数のコメントを回答、ミーティングに残すことができる"

%% コメントのタイプ（アンサー単位かミーティング自体）のマスタテーブル
COMMENT_TYPES {
    number id PK
    string name "アンサー|ミーティング"
}

```