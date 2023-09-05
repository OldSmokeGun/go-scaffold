-- +migrate Up

CREATE TABLE IF NOT EXISTS `products`
(
    `id`         int unsigned    NOT NULL AUTO_INCREMENT,
    `name`       varchar(128)    NOT NULL DEFAULT '' COMMENT '名称',
    `desc`       varchar(255)    NOT NULL DEFAULT '' COMMENT '描述',
    `price`      int             NOT NULL DEFAULT 0 COMMENT '价格',
    `created_at` bigint          NOT NULL DEFAULT 0,
    `updated_at` bigint          NOT NULL DEFAULT 0,
    `deleted_at` bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `name` (`name`),
    KEY `deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='产品表';

-- +migrate Down

DROP TABLE IF EXISTS `products`;