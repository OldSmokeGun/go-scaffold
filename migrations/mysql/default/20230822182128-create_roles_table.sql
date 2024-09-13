-- +migrate Up

CREATE TABLE IF NOT EXISTS `roles`
(
    `id`         int unsigned    NOT NULL AUTO_INCREMENT,
    `name`       varchar(32)     NOT NULL DEFAULT '' COMMENT '角色名称',
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE `name` (`name`),
    KEY `deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='角色表';

-- +migrate Down

DROP TABLE IF EXISTS `roles`;