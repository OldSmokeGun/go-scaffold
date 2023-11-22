-- +migrate Up

CREATE TABLE IF NOT EXISTS roles
(
    id         bigserial   NOT NULL,
    name       varchar(32) NOT NULL DEFAULT '',
    created_at bigint      NOT NULL DEFAULT 0,
    updated_at bigint      NOT NULL DEFAULT 0,
    deleted_at bigint      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE INDEX ON roles (deleted_at);

COMMENT ON COLUMN roles.name IS '角色名称';

COMMENT ON TABLE roles IS '角色表';

-- +migrate Down

DROP TABLE IF EXISTS roles;