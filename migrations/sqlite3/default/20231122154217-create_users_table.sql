-- +migrate Up

CREATE TABLE IF NOT EXISTS `users`
(
    `id`         integer PRIMARY KEY AUTOINCREMENT,
    `username`   varchar(32) UNIQUE NOT NULL DEFAULT '', -- 用户名
    `password`   varchar(64)        NOT NULL DEFAULT '', -- 密码
    `nickname`   varchar(64)        NOT NULL DEFAULT '', -- 昵称
    `phone`      varchar(11)        NOT NULL DEFAULT '', -- 电话
    `salt`       varchar(64)        NOT NULL DEFAULT '', -- 盐值
    `created_at` bigint             NOT NULL DEFAULT 0,
    `updated_at` bigint             NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned    NOT NULL DEFAULT 0
);

CREATE INDEX users_phone ON users (phone);
CREATE INDEX users_deleted_at ON users (deleted_at);

-- +migrate Down

DROP TABLE IF EXISTS `users`;