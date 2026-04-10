-- ============================================
-- CREATE TABLE
-- ============================================

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL
);

CREATE TABLE support_statuses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    passwordd VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role_id BIGINT NOT NULL,
    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);

CREATE TABLE supports (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    support_status_id BIGINT NOT NULL,
    CONSTRAINT fk_supports_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),
    CONSTRAINT fk_supports_support_status
        FOREIGN KEY (support_status_id)
        REFERENCES support_statuses(id)
);

CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT,
    origin_question_id BIGINT,
    support_id BIGINT,
    title VARCHAR(255) NOT NULL,
    due TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_questions_origin
        FOREIGN KEY (origin_question_id)
        REFERENCES questions(id),
    CONSTRAINT fk_questions_support
        FOREIGN KEY (support_id)
        REFERENCES supports(id)
);

CREATE TABLE answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    answered_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_answers_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),
    CONSTRAINT fk_answers_question
        FOREIGN KEY (question_id)
        REFERENCES questions(id)
);

CREATE TABLE memos (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    CONSTRAINT fk_memos_question
        FOREIGN KEY (question_id)
        REFERENCES questions(id),
    CONSTRAINT fk_memos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

CREATE TABLE refers (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    usage INTEGER NOT NULL DEFAULT 0,
    category_id BIGINT NOT NULL,
    CONSTRAINT fk_tags_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
);

CREATE TABLE refer_managers (
    id BIGSERIAL PRIMARY KEY,
    answer_id BIGINT NOT NULL,
    refer_id BIGINT NOT NULL,
    CONSTRAINT fk_refer_managers_answer
        FOREIGN KEY (answer_id)
        REFERENCES answers(id),
    CONSTRAINT fk_refer_managers_refer
        FOREIGN KEY (refer_id)
        REFERENCES refers(id)
);

CREATE TABLE tag_managers (
    id BIGSERIAL PRIMARY KEY,
    tag_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    CONSTRAINT fk_tag_managers_tag
        FOREIGN KEY (tag_id)
        REFERENCES tags(id),
    CONSTRAINT fk_tag_managers_question
        FOREIGN KEY (question_id)
        REFERENCES questions(id)
);

CREATE TABLE escalations (
    id BIGSERIAL PRIMARY KEY,
    from_question_id BIGINT NOT NULL,
    to_question_id BIGINT NOT NULL,
    escalated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_escalations_from_question
        FOREIGN KEY (from_question_id)
        REFERENCES questions(id),
    CONSTRAINT fk_escalations_to_question
        FOREIGN KEY (to_question_id)
        REFERENCES questions(id)
);


-- ============================================
-- INDEX
-- ============================================

CREATE INDEX idx_users_role_id
ON users(role_id);

CREATE INDEX idx_supports_user_id
ON supports(user_id);

CREATE INDEX idx_supports_support_status_id
ON supports(support_status_id);

CREATE INDEX idx_questions_origin_question_id
ON questions(origin_question_id);

CREATE INDEX idx_questions_support_id
ON questions(support_id);

CREATE INDEX idx_answers_user_id
ON answers(user_id);

CREATE INDEX idx_answers_question_id
ON answers(question_id);

CREATE INDEX idx_memos_question_id
ON memos(question_id);

CREATE INDEX idx_memos_user_id
ON memos(user_id);

CREATE INDEX idx_tags_category_id
ON tags(category_id);

CREATE INDEX idx_refer_managers_answer_id
ON refer_managers(answer_id);

CREATE INDEX idx_refer_managers_refer_id
ON refer_managers(refer_id);

CREATE INDEX idx_tag_managers_tag_id
ON tag_managers(tag_id);

CREATE INDEX idx_tag_managers_question_id
ON tag_managers(question_id);

CREATE INDEX idx_escalations_from_question_id
ON escalations(from_question_id);

CREATE INDEX idx_escalations_to_question_id
ON escalations(to_question_id);


-- ============================================
-- INSERT INTO
-- 外部キー依存順に並べています
-- ============================================

-- roles
INSERT INTO roles (role_name) VALUES
('Admin'),
('Manager'),
('Staff'),
('Employee');

-- categories
INSERT INTO categories (category_name) VALUES
('総務'),
('人事'),
('その他');

-- support_statuses
INSERT INTO support_statuses (title) VALUES
('未対応'),
('対応中'),
('完了');

-- users
INSERT INTO users (name, passwordd, email, role_id) VALUES
('Taro Yamada', 'hashed_password_1', 'taro@example.com', 1),
('Hanako Suzuki', 'hashed_password_2', 'hanako@example.com', 2),
('Jiro Tanaka', 'hashed_password_3', 'jiro@example.com', 3),
('Sato Hiromichi', 'hashed_password_4', 'sato@example.com', 4);

-- supports
INSERT INTO supports (user_id, support_status_id) VALUES
(1, 1),
(2, 2),
(3, 3);

-- questions
-- origin_question_id は自己参照なので親→子の順で入れています
INSERT INTO questions (message_id, origin_question_id, support_id, title, due, created_at) VALUES
(1001, NULL, 1, 'First Question', '2026-04-30 12:00:00', CURRENT_TIMESTAMP),
(1002, 1,    2, 'Follow-up Question', '2026-05-01 12:00:00', CURRENT_TIMESTAMP),
(1003, NULL, 3, 'Health Check Schedule', '2026-05-10 09:00:00', CURRENT_TIMESTAMP);

-- answers
INSERT INTO answers (user_id, question_id, content, answered_at, created_at) VALUES
(1, 1, 'This is an answer to the first question', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 2, 'Answer to follow-up question', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 3, 'Health check date can be changed via internal form.', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- memos
INSERT INTO memos (question_id, user_id, content) VALUES
(1, 1, '途中メモ１'),
(1, 2, '途中メモ２'),
(2, 2, '／(^o^)＼'),
(3, 3, '健康診断の規程を参照する');

-- refers
INSERT INTO refers (title, url) VALUES
('PostgreSQL Documentation', 'https://www.postgresql.org/docs/'),
('GORM Official', 'https://gorm.io'),
('社内総務規程', 'https://intra.example.local/rules/general-affairs');

-- tags
INSERT INTO tags (title, usage, category_id) VALUES
('諸手当', 1, 1),
('休暇', 1, 2),
('規程', 1, 3),
('健康診断', 1, 1);

-- refer_managers
INSERT INTO refer_managers (answer_id, refer_id) VALUES
(1, 1),
(1, 2),
(2, 1),
(3, 3);

-- tag_managers
INSERT INTO tag_managers (tag_id, question_id) VALUES
(1, 1),
(2, 2),
(4, 3);

-- escalations
INSERT INTO escalations (from_question_id, to_question_id, escalated_at) VALUES
(1, 2, CURRENT_TIMESTAMP);