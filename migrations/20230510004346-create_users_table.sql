-- +migrate Up

CREATE TABLE IF NOT EXISTS `users`
(
    `id`         int unsigned    NOT NULL AUTO_INCREMENT,
    `name`       varchar(64)     NOT NULL DEFAULT '' COMMENT '名称',
    `age`        tinyint         NOT NULL DEFAULT 0 COMMENT '年龄',
    `phone`      varchar(11)     NOT NULL DEFAULT '' COMMENT '电话',
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

-- +migrate Down

DROP TABLE IF EXISTS `users`;