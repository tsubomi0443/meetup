```mermaid
graph LR;

mentee ~~~ ANS_FUNC
mentor ~~~ ANS_FUNC
boss ~~~ ANS_FUNC
mentor~~~mentor
ANS_FUNC~~~MNG_FUNC
ANS_FUNC~~~admin
admin~~~MNG_FUNC

%% カッコ内は権限を現す 管理者>上長>一般のパワー感
mentee(("メンティー（一般）"))
%% メンターはペアリングで登録された人物の回答の閲覧と回答に対するコメント、回答を評価する権限が与えられる
mentor(("メンター（一般）"))
%% 上長は配下すべてに対してメンターと同様の権限を持つ
boss(("上司（上長）"))
%% オーナーはすべての権限を持つ
admin(("管理者（オーナー）"))

mentee-->input_topic
mentee-->temp_topic
mentee-->regist_topic
mentee-->|自分の回答|refer_topic

mentor-->|メンティーの回答|refer_topic
mentor-->comment_topic

boss-->|配下ユーザの回答|refer_topic
boss-->|配下ユーザへのコメント|comment_topic

subgraph ANS_FUNC[回答機能]
    direction TB;
    input_topic(回答を入力する)
    temp_topic(回答を仮保存する)
    regist_topic(回答を提出する)

    refer_topic(回答を閲覧する)
    comment_topic(回答に対してコメントする)
end

subgraph MNG_FUNC[管理機能]
    direction LR;
    crud_user("ユーザ<br/>登録/編集/閲覧/削除")
    crud_department("部署<br/>登録/編集/閲覧/削除")
    crud_template("テンプレート<br/>登録/編集/閲覧/削除")
    crud_peer("ペア<br/>登録/編集/閲覧/削除")
    check_progress(実施状況確認)
end

admin-->crud_peer
admin-->crud_user
admin-->crud_template
admin-->crud_department
admin-->check_progress
admin-->|全ユーザの回答|refer_topic
```

## HR Assist（本リポ mock5）の認可メモ

- ログイン後の JWT は typed claims（`user_id` / `email` / `role_id` / `name`（ロール名））で発行する。
- mock5 の SSR では `roles.id` が **1=Admin**、**2=Manager** のときだけユーザー管理ビュー・招待/編集モーダルを出す（それ以外は Staff 等）。
- mock5 ではタグ/ユーザー/詳細の一部 UI 向けに、DB マスタの **roles / categories / support_statuses** を SSR で渡し、フロントのハードコードを減らす。
