BEGIN;


-- DROP


DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;

-- ============================================
-- CREATE TABLE
-- ============================================

-- ROLEのCREATE
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- USERのCREATE
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    memo TEXT DEFAULT '',
    role_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- SUPPORT_STATUSのCREATE
CREATE TABLE support_statuses (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- SUPPORTのCREATE
CREATE TABLE supports (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    support_status_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- QUESTIONのCREATE
CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    origin_question_id BIGINT,
    support_id BIGINT,
    talkroom_id VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    due TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- ANSWERのCREATE
CREATE TABLE answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    is_final BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- MEMOのCREATE
CREATE TABLE memos (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- CATEGORYのCREATE
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- TAGのCREATE
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    usage INTEGER NOT NULL DEFAULT 0,
    category_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- TAG_MANAGERのCREATE
CREATE TABLE tag_managers (
    id BIGSERIAL PRIMARY KEY,
    tag_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- REFERのCREATE
CREATE TABLE refers (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- REFER_MANAGERのCREATE
CREATE TABLE refer_managers (
    id BIGSERIAL PRIMARY KEY,
    answer_id BIGINT NOT NULL,
    refer_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- ESCALATIONのCREATE
CREATE TABLE escalations (
    id BIGSERIAL PRIMARY KEY,
    from_question_id BIGINT NOT NULL,
    to_question_id BIGINT NOT NULL,
    escalated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- NOTICE_TYPEのCREATE
CREATE TABLE notice_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- NOTICEのCREATE
CREATE TABLE notices (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL,
    question_id BIGINT,
    content TEXT,
    display_due TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- RELATED_QUESTIONのCREATE
CREATE TABLE related_questions (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL,
    related_question_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_no_self_reference
        CHECK (question_id <> related_question_id),
    CONSTRAINT uq_related_questions
        UNIQUE (question_id, related_question_id)
);

-- SENDERのCREATE
CREATE TABLE senders (
    id BIGSERIAL PRIMARY KEY,
    uid VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    department_name VARCHAR(255)
);

-- SENDER_TALKのCREATE
CREATE TABLE sender_talks (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    talkroom_id VARCHAR(255),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- ============================================
-- CONSTRAINT
-- ============================================

-- USERのCONSTRAINT
ALTER TABLE users
    ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id);

-- SUPPORTのCONSTRAINT
ALTER TABLE supports
    ADD CONSTRAINT fk_support_user FOREIGN KEY (user_id) REFERENCES users(id),
    ADD CONSTRAINT fk_support_status FOREIGN KEY (support_status_id) REFERENCES support_statuses(id);

-- QUESTIONのCONSTRAINT
ALTER TABLE questions
    ADD CONSTRAINT fk_question_origin FOREIGN KEY (origin_question_id) REFERENCES questions(id),
    ADD CONSTRAINT fk_question_support FOREIGN KEY (support_id) REFERENCES supports(id);

-- ANSWERのCONSTRAINT
ALTER TABLE answers
    ADD CONSTRAINT fk_answer_user FOREIGN KEY (user_id) REFERENCES users(id),
    ADD CONSTRAINT fk_answer_question FOREIGN KEY (question_id) REFERENCES questions(id);

-- MEMOのCONSTRAINT
ALTER TABLE memos
    ADD CONSTRAINT fk_memo_question FOREIGN KEY (question_id) REFERENCES questions(id),
    ADD CONSTRAINT fk_memo_user FOREIGN KEY (user_id) REFERENCES users(id);

-- CATEGORY/TAGのCONSTRAINT
ALTER TABLE tags
    ADD CONSTRAINT fk_tag_category FOREIGN KEY (category_id) REFERENCES categories(id);

-- TAG_MANAGERのCONSTRAINT
ALTER TABLE tag_managers
    ADD CONSTRAINT fk_tm_tag FOREIGN KEY (tag_id) REFERENCES tags(id),
    ADD CONSTRAINT fk_tm_question FOREIGN KEY (question_id) REFERENCES questions(id);

-- REFER_MANAGERのCONSTRAINT
ALTER TABLE refer_managers
    ADD CONSTRAINT fk_rm_answer FOREIGN KEY (answer_id) REFERENCES answers(id),
    ADD CONSTRAINT fk_rm_refer FOREIGN KEY (refer_id) REFERENCES refers(id);

-- ESCALATIONのCONSTRAINT
ALTER TABLE escalations
    ADD CONSTRAINT fk_es_from FOREIGN KEY (from_question_id) REFERENCES questions(id),
    ADD CONSTRAINT fk_es_to FOREIGN KEY (to_question_id) REFERENCES questions(id);

-- NOTICEのCONSTRAINT
ALTER TABLE notices
    ADD CONSTRAINT fk_notice_type FOREIGN KEY (type_id) REFERENCES notice_types(id),
    ADD CONSTRAINT fk_notice_question FOREIGN KEY (question_id) REFERENCES questions(id);

-- RELATED_QUESTIONのCONSTRAINT
ALTER TABLE related_questions
    ADD CONSTRAINT fk_rq_q FOREIGN KEY (question_id) REFERENCES questions(id),
    ADD CONSTRAINT fk_rq_related FOREIGN KEY (related_question_id) REFERENCES questions(id);

-- SENDER_TALKのCONSTRAINT
ALTER TABLE sender_talks
    ADD CONSTRAINT fk_sender_talk_sender FOREIGN KEY (sender_id) REFERENCES senders(id),
    ADD CONSTRAINT fk_sender_talk_question FOREIGN KEY (question_id) REFERENCES questions(id);

-- ============================================
-- INDEX
-- ============================================

-- USERのINDEX
CREATE INDEX idx_users_role ON users(role_id);

-- SUPPORTのINDEX
CREATE INDEX idx_support_user ON supports(user_id);
CREATE INDEX idx_support_status ON supports(support_status_id);

-- QUESTIONのINDEX
CREATE INDEX idx_question_origin ON questions(origin_question_id);
CREATE INDEX idx_question_support ON questions(support_id);

-- ANSWERのINDEX
CREATE INDEX idx_answer_question ON answers(question_id);
CREATE INDEX idx_answer_user ON answers(user_id);

-- MEMOのINDEX
CREATE INDEX idx_memo_question ON memos(question_id);
CREATE INDEX idx_memo_user ON memos(user_id);

-- TAGのINDEX
CREATE INDEX idx_tag_category ON tags(category_id);

-- TAG_MANAGERのINDEX
CREATE INDEX idx_tm_tag ON tag_managers(tag_id);
CREATE INDEX idx_tm_question ON tag_managers(question_id);

-- REFER_MANAGERのINDEX
CREATE INDEX idx_rm_answer ON refer_managers(answer_id);
CREATE INDEX idx_rm_refer ON refer_managers(refer_id);

-- ESCALATIONのINDEX
CREATE INDEX idx_es_from ON escalations(from_question_id);
CREATE INDEX idx_es_to ON escalations(to_question_id);

-- NOTICEのINDEX
CREATE INDEX idx_notice_type ON notices(type_id);
CREATE INDEX idx_notice_question ON notices(question_id);

-- RELATED_QUESTIONのINDEX
CREATE INDEX idx_rq_question ON related_questions(question_id);

-- SENDER_TALKのINDEX
CREATE INDEX idx_sender_talk_question ON sender_talks(question_id);
CREATE INDEX idx_sender_talk_sender ON sender_talks(sender_id);


-- INSERT
-- 外部キー依存順


-- ROLEのINSERT
INSERT INTO roles (name) VALUES
('Admin'),
('Manager'),
('Staff'),
('Employee');

-- CATEGORYのINSERT
INSERT INTO categories (name) VALUES
('総務'),
('人事'),
('その他');

-- SUPPORT_STATUSのINSERT
INSERT INTO support_statuses (name) VALUES
('未対応'),
('対応中'),
('完了');

-- USERのINSERT
-- Admin以外の初期パスワードは「password」、Adminは「admin」
INSERT INTO users (name, password, email, role_id) VALUES
('admin', '8d52e45b42babab443e7be6a211a4111172eb7efd6e8b9abb6f67039c1297415', 'admin', 1),
('Taro Yamada', '101dd96cf87d902a2c73d60979fa7627e693d34651dea00019649eb289a1d6b5', 'taro@example.com', 2),
('Hanako Suzuki', '101dd96cf87d902a2c73d60979fa7627e693d34651dea00019649eb289a1d6b5', 'hanako@example.com', 3),
('Jiro Tanaka', '101dd96cf87d902a2c73d60979fa7627e693d34651dea00019649eb289a1d6b5', 'jiro@example.com', 3),
('Sato Hiromichi', '101dd96cf87d902a2c73d60979fa7627e693d34651dea00019649eb289a1d6b5', 'sato@example.com', 4);

-- SUPPORTのINSERT
INSERT INTO supports (user_id, support_status_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- QUESTIONのINSERT
INSERT INTO questions (origin_question_id, support_id, talkroom_id, title, content, due) VALUES
(NULL, 1, 'room-1', 'First Question', 'First question body', '2026-04-30 12:00:00'),
(1, 2, 'room-2', 'Follow-up Question', 'Follow-up body', '2026-05-01 12:00:00'),
(NULL, 3, 'room-3', 'Health Check Schedule', 'Health check schedule details', '2026-05-10 09:00:00');

-- ANSWERのINSERT
INSERT INTO answers (user_id, question_id, content, is_final) VALUES
(1, 1, 'This is an answer to the first question', true),
(2, 2, 'Answer to follow-up question', true),
(3, 3, 'Health check date can be changed via internal form.', true);

-- MEMOのINSERT
INSERT INTO memos (question_id, user_id, content) VALUES
(1, 1, '途中メモ１'),
(1, 2, '途中メモ２'),
(2, 2, '／(^o^)＼'),
(3, 3, '健康診断の規程を参照する');

-- REFERのINSERT
INSERT INTO refers (title, url) VALUES
('PostgreSQL Documentation', 'https://www.postgresql.org/docs/'),
('GORM Official', 'https://gorm.io'),
('社内総務規程', 'https://intra.example.local/rules/general-affairs');

-- TAGのINSERT
INSERT INTO tags (name, usage, category_id) VALUES
('諸手当', 1, 1),
('休暇', 1, 2),
('規程', 1, 3),
('健康診断', 1, 1);

-- REFER_MANAGERのINSERT
INSERT INTO refer_managers (answer_id, refer_id) VALUES
(1, 1),
(1, 2),
(2, 1),
(3, 3);

-- TAG_MANAGERのINSERT
INSERT INTO tag_managers (tag_id, question_id) VALUES
(1, 1),
(2, 2),
(4, 3);

-- ESCALATIONのINSERT
INSERT INTO escalations (from_question_id, to_question_id, escalated_at) VALUES
(1, 2, CURRENT_TIMESTAMP);

-- NOTICE_TYPEのINSERT
INSERT INTO notice_types (name) VALUES
('SYSTEM'),
('ALERT'),
('QUESTION');

-- NOTICEのINSERT
INSERT INTO notices (type_id, question_id, content, display_due) VALUES
(2, 1, 'First Question の期限が近づいています', '2026-04-29 09:00:00'),
(1, NULL, 'システムメンテナンスのお知らせ', NULL);

-- RELATED_QUESTIONのINSERT
INSERT INTO related_questions (question_id, related_question_id) VALUES
(2, 1),
(3, 1),
(3, 2);

-- SENDERのINSERT
INSERT INTO senders (uid, name, department_name) VALUES
('lw-uid-001', '外部ユーザA', 'Sales'),
('lw-uid-002', '外部ユーザB', 'HR');

-- SENDER_TALKのINSERT
INSERT INTO sender_talks (sender_id, question_id, talkroom_id, content) VALUES
(1, 1, 'room-1', '質問を送信しました'),
(1, 1, 'room-1', '追加情報です'),
(2, 2, 'room-2', '別の質問です');

COMMIT;
