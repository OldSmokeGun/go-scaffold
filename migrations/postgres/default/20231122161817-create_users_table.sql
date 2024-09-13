-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
    id         bigserial   NOT NULL,
    username   varchar(32) NOT NULL DEFAULT '',
    password   varchar(64) NOT NULL DEFAULT '',
    nickname   varchar(64) NOT NULL DEFAULT '',
    phone      varchar(11) NOT NULL DEFAULT '',
    salt       varchar(64) NOT NULL DEFAULT '',
    created_at bigint      NOT NULL DEFAULT 0,
    updated_at bigint      NOT NULL DEFAULT 0,
    deleted_at bigint      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (username)
);

CREATE INDEX ON users (phone);
CREATE INDEX ON users (deleted_at);

COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password IS '密码';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.phone IS '电话';
COMMENT ON COLUMN users.salt IS '盐值';

COMMENT ON TABLE users IS '用户表';

-- +migrate Down

DROP TABLE IF EXISTS users;