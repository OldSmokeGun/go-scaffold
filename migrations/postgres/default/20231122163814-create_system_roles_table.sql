-- +migrate Up

CREATE TABLE IF NOT EXISTS system_roles
(
    id         bigserial   NOT NULL,
    name       varchar(32) NOT NULL DEFAULT '',
    created_at bigint      NOT NULL DEFAULT 0,
    updated_at bigint      NOT NULL DEFAULT 0,
    deleted_at bigint      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE INDEX ON system_roles (deleted_at);

COMMENT ON COLUMN system_roles.name IS '角色名称';

COMMENT ON TABLE system_roles IS '角色表';

-- +migrate Down

DROP TABLE IF EXISTS system_roles;