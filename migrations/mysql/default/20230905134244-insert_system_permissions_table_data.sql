-- +migrate Up

START TRANSACTION;

INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at) VALUES ('/system-users', '用户管理', 0, unix_timestamp(), unix_timestamp());
INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/system-users', '用户列表', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/system-user/:id', '用户详情', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/system-user', '用户新增', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('PUT /api/v1/system-user', '用户更新', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('DELETE /api/v1/system-user/:id', '用户删除', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/system-user/roles', '分配用户角色', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/system-user/roles', '获取用户角色', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-users') AS t), unix_timestamp(), unix_timestamp());

INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at) VALUES ('/system-roles', '角色管理', 0, unix_timestamp(), unix_timestamp());
INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/system-roles', '角色列表', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/system-role/:id', '角色详情', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/system-role', '角色新增', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('PUT /api/v1/system-role', '角色更新', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('DELETE /api/v1/system-role/:id', '角色删除', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/system-role/permissions', '授予角色权限', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/system-role/permissions', '获取角色权限', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-roles') AS t), unix_timestamp(), unix_timestamp());

INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at) VALUES ('/system-permissions', '权限管理', 0, unix_timestamp(), unix_timestamp());
INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/system-permissions', '权限列表', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-permissions') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/system-permission/:id', '权限详情', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-permissions') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/system-permission', '权限新增', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-permissions') AS t), unix_timestamp(), unix_timestamp()),
       ('PUT /api/v1/system-permission', '权限更新', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-permissions') AS t), unix_timestamp(), unix_timestamp()),
       ('DELETE /api/v1/system-permission/:id', '权限删除', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/system-permissions') AS t), unix_timestamp(), unix_timestamp());

INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at) VALUES ('/products', '产品管理', 0, unix_timestamp(), unix_timestamp());
INSERT INTO system_permissions (`key`, name, parent_id, created_at, updated_at)
VALUES ('GET /api/v1/products', '产品列表', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/products') AS t), unix_timestamp(), unix_timestamp()),
       ('GET /api/v1/product/:id', '产品详情', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/products') AS t), unix_timestamp(), unix_timestamp()),
       ('POST /api/v1/product', '产品新增', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/products') AS t), unix_timestamp(), unix_timestamp()),
       ('PUT /api/v1/product', '产品更新', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/products') AS t), unix_timestamp(), unix_timestamp()),
       ('DELETE /api/v1/product/:id', '产品删除', (SELECT id FROM (SELECT id FROM system_permissions WHERE `key` = '/products') AS t), unix_timestamp(), unix_timestamp());

COMMIT;

-- +migrate Down

TRUNCATE TABLE system_permissions;
