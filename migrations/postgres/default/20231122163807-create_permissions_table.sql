-- +migrate Up

CREATE TABLE IF NOT EXISTS permissions
(
    id         bigserial    NOT NULL,
    key        varchar(128) NOT NULL DEFAULT '',
    name       varchar(128) NOT NULL DEFAULT '',
    "desc"     varchar(255) NOT NULL DEFAULT '',
    parent_id  int          NOT NULL DEFAULT 0,
    created_at bigint       NOT NULL DEFAULT 0,
    updated_at bigint       NOT NULL DEFAULT 0,
    deleted_at bigint       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (key)
);

CREATE INDEX ON permissions (deleted_at);

COMMENT ON COLUMN permissions.key IS '权限标识';
COMMENT ON COLUMN permissions.name IS '权限名称';
COMMENT ON COLUMN permissions.desc IS '权限描述';
COMMENT ON COLUMN permissions.parent_id IS '父级权限 id';

COMMENT ON TABLE permissions IS '权限表';

-- +migrate Down

DROP TABLE IF EXISTS permissions;