ALTER TABLE users
ADD COLUMN last_login_at timestamptz,
ADD COLUMN last_login_ip inet;

CREATE INDEX idx_users_last_login ON users (last_login_at DESC);
