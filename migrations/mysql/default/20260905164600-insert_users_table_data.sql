-- +migrate Up

START TRANSACTION;

INSERT INTO users (username, password, nickname, mobile, salt, created_at, updated_at, deleted_at)
VALUES ('admin', MD5('Admin123!'), 'SuperAdmin', '', 'f521a082-de1c-4519-9829-a12ef133f02c', unix_timestamp(), unix_timestamp(), 0);

COMMIT;

-- +migrate Down

DELETE FROM users WHERE username = 'admin';
