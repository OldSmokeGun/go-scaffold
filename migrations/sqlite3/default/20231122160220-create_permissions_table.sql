-- +migrate Up

CREATE TABLE IF NOT EXISTS `permissions`
(
    `id`         integer PRIMARY KEY AUTOINCREMENT,
    `key`        varchar(128) UNIQUE NOT NULL DEFAULT '', -- 权限标识
    `name`       varchar(128)        NOT NULL DEFAULT '', -- 权限名称
    `desc`       varchar(255)        NOT NULL DEFAULT '', -- 权限描述
    `parent_id`  int unsigned        NOT NULL DEFAULT 0,  -- 父级权限 id
    `created_at` bigint              NOT NULL DEFAULT 0,
    `updated_at` bigint              NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned     NOT NULL DEFAULT 0
);

CREATE INDEX permissions_deleted_at ON permissions (deleted_at);

-- +migrate Down

DROP TABLE IF EXISTS `permissions`;