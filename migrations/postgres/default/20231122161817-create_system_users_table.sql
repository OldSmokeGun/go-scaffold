-- +migrate Up

CREATE TABLE IF NOT EXISTS system_users
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

CREATE INDEX ON system_users (phone);
CREATE INDEX ON system_users (deleted_at);

COMMENT ON COLUMN system_users.username IS '用户名';
COMMENT ON COLUMN system_users.password IS '密码';
COMMENT ON COLUMN system_users.nickname IS '昵称';
COMMENT ON COLUMN system_users.phone IS '电话';
COMMENT ON COLUMN system_users.salt IS '盐值';

COMMENT ON TABLE system_users IS '用户表';

-- +migrate Down

DROP TABLE IF EXISTS system_users;