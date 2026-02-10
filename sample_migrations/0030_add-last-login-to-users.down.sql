DROP INDEX IF EXISTS idx_users_last_login;

ALTER TABLE users
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS last_login_ip;
