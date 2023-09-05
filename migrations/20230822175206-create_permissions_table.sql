-- +migrate Up

CREATE TABLE IF NOT EXISTS `permissions`
(
    `id`         int unsigned    NOT NULL AUTO_INCREMENT,
    `key`        varchar(128)    NOT NULL DEFAULT '' COMMENT '权限标识',
    `name`       varchar(128)    NOT NULL DEFAULT '' COMMENT '权限名称',
    `desc`       varchar(255)    NOT NULL DEFAULT '' COMMENT '权限描述',
    `parent_id`  int unsigned    NOT NULL DEFAULT 0 COMMENT '父级权限 id',
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE `key` (`key`),
    KEY `deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='权限表';

-- +migrate Down

DROP TABLE IF EXISTS `permissions`;