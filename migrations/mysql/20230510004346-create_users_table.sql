-- +migrate Up

CREATE TABLE IF NOT EXISTS `users`
(
    `id`         int unsigned    NOT NULL AUTO_INCREMENT,
    `username`   varchar(32)     NOT NULL DEFAULT '' COMMENT '用户名',
    `password`   varchar(64)     NOT NULL DEFAULT '' COMMENT '密码',
    `nickname`   varchar(64)     NOT NULL DEFAULT '' COMMENT '昵称',
    `phone`      varchar(11)     NOT NULL DEFAULT '' COMMENT '电话',
    `salt`       varchar(64)     NOT NULL DEFAULT '' COMMENT '盐值',
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE `username` (`username`),
    KEY `phone` (`phone`),
    KEY `deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

-- +migrate Down

DROP TABLE IF EXISTS `users`;