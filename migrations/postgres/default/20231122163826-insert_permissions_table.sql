-- +migrate Up

INSERT INTO permissions (key, name, parent_id, created_at, updated_at) VALUES ('/users', '用户管理', 0, (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));
INSERT INTO permissions (key, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/users', '用户列表', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/user/:id', '用户详情', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/user', '用户新增', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('PUT /api/v1/user', '用户更新', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('DELETE /api/v1/user/:id', '用户删除', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/user/roles', '分配用户角色', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/user/roles', '获取用户角色', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/users') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));

INSERT INTO permissions (key, name, parent_id, created_at, updated_at) VALUES ('/roles', '角色管理', 0, (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));
INSERT INTO permissions (key, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/roles', '角色列表', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/role/:id', '角色详情', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/role', '角色新增', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('PUT /api/v1/role', '角色更新', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('DELETE /api/v1/role/:id', '角色删除', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/role/permissions', '授予角色权限', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/role/permissions', '获取角色权限', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/roles') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));

INSERT INTO permissions (key, name, parent_id, created_at, updated_at) VALUES ('/permissions', '权限管理', 0, (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));
INSERT INTO permissions (key, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/permissions', '权限列表', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/permissions') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/permission/:id', '权限详情', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/permissions') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/permission', '权限新增', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/permissions') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('PUT /api/v1/permission', '权限更新', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/permissions') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('DELETE /api/v1/permission/:id', '权限删除', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/permissions') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));

INSERT INTO permissions (key, name, parent_id, created_at, updated_at) VALUES ('/products', '产品管理', 0, (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));
INSERT INTO permissions (key, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/products', '产品列表', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/products') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('GET /api/v1/product/:id', '产品详情', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/products') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('POST /api/v1/product', '产品新增', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/products') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('PUT /api/v1/product', '产品更新', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/products') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0)))),
       ('DELETE /api/v1/product/:id', '产品删除', (SELECT id FROM (SELECT id FROM permissions WHERE key = '/products') AS t), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))), (SELECT EXTRACT(EPOCH FROM now()::timestamp(0))));

-- +migrate Down

TRUNCATE TABLE permissions;
