BEGIN;

-- ============================================
-- DROP（子→親）
-- ============================================

DROP TABLE IF EXISTS refer_managers;
DROP TABLE IF EXISTS memo_managers;
DROP TABLE IF EXISTS tag_managers;
DROP TABLE IF EXISTS escalations;
DROP TABLE IF EXISTS answers;

DROP TABLE IF EXISTS memos;
DROP TABLE IF EXISTS refers;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS questions;

DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS categories;



-- ============================================
-- CREATE（親→子）
-- ============================================

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL
);

CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT,
    origin_question_id BIGINT,
    title VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    due TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (origin_question_id)
        REFERENCES questions(id)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    passwordd VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role_id BIGINT NOT NULL,
    FOREIGN KEY (role_id)
        REFERENCES roles(id)
);

CREATE TABLE memos (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL
);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    usage INTEGER NOT NULL DEFAULT 0,
    category_id BIGINT NOT NULL,
    FOREIGN KEY (category_id)
        REFERENCES categories(id)
);

CREATE TABLE refers (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE answers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    answered_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id),
    FOREIGN KEY (question_id)
        REFERENCES questions(id)
);

CREATE TABLE memo_managers (
    id BIGSERIAL PRIMARY KEY,
    memo_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    FOREIGN KEY (memo_id)
        REFERENCES memos(id),
    FOREIGN KEY (question_id)
        REFERENCES questions(id)
);

CREATE TABLE tag_managers (
    id BIGSERIAL PRIMARY KEY,
    tag_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    FOREIGN KEY (tag_id)
        REFERENCES tags(id),
    FOREIGN KEY (question_id)
        REFERENCES questions(id)
);

CREATE TABLE refer_managers (
    id BIGSERIAL PRIMARY KEY,
    answer_id BIGINT NOT NULL,
    refer_id BIGINT NOT NULL,
    FOREIGN KEY (answer_id)
        REFERENCES answers(id),
    FOREIGN KEY (refer_id)
        REFERENCES refers(id)
);

CREATE TABLE escalations (
    id BIGSERIAL PRIMARY KEY,
    from_question_id BIGINT NOT NULL,
    to_question_id BIGINT NOT NULL,
    escalated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (from_question_id)
        REFERENCES questions(id),
    FOREIGN KEY (to_question_id)
        REFERENCES questions(id)
);



-- ============================================
-- INSERT
-- ============================================

INSERT INTO roles (role_name) VALUES
    ('Manager'),
    ('Creator');

INSERT INTO categories (category_name) VALUES
    ('労務'),
    ('総務');

INSERT INTO tags (title, usage, category_id) VALUES
    ('諸手当',   0, 1),
    ('休暇',     0, 2),
    ('規程',     0, 1),
    ('健康診断', 0, 2);

INSERT INTO users (name, passwordd, email, role_id) VALUES
    ('自分',      'placeholder', 'jinji.taro@example.com', 1),
    ('鈴木 一郎', 'placeholder', 'suzuki@example.com',     2),
    ('田中 花子', 'placeholder', 'tanaka@example.com',     2);

INSERT INTO questions (title, status, due, created_at) VALUES
    ('通勤手当の経路変更について', 1, '2026-04-10', '2026-04-08 09:30'),
    ('育児休業の延長申請',         2, '2026-04-14', '2026-04-04 14:00'),
    ('健康診断の受診日変更',       3, '2026-04-09', '2026-04-01 11:15'),
    ('慶弔休暇の適用範囲',         1, '2026-04-11', '2026-04-07 16:45');

INSERT INTO tag_managers (tag_id, question_id) VALUES
    (1, 1),
    (2, 2),
    (4, 3),
    (3, 4);

COMMIT;