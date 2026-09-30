-- +migrate Up

CREATE TABLE IF NOT EXISTS `system_roles`
(
    `id`         integer PRIMARY KEY AUTOINCREMENT,
    `name`       varchar(32) UNIQUE NOT NULL DEFAULT '', -- 角色名称
    `created_at` bigint             NOT NULL DEFAULT 0,
    `updated_at` bigint             NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned    NOT NULL DEFAULT 0
);

CREATE INDEX system_roles_deleted_at ON system_roles (deleted_at);

-- +migrate Down

DROP TABLE IF EXISTS `system_roles`;