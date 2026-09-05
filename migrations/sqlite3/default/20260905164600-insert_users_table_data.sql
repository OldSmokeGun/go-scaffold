-- +migrate Up

-- User.Create: password = md5(Admin123!); salt 与 Register 一样用 uuid
INSERT INTO users (`username`, `password`, `nickname`, `phone`, `salt`, `created_at`, `updated_at`, `deleted_at`)
VALUES ('admin', '2637a5c30af69a7bad877fdb65fbd78b', 'SuperAdmin', '13800138000', 'f521a082-de1c-4519-9829-a12ef133f02c', strftime('%s', 'now'), strftime('%s', 'now'), 0);

-- +migrate Down

DELETE FROM users WHERE username = 'admin';
