ALTER TABLE users
ADD COLUMN email_verified boolean DEFAULT false NOT NULL,
ADD COLUMN email_verification_token varchar(255),
ADD COLUMN email_verified_at timestamptz;

CREATE INDEX idx_users_verification_token ON users (email_verification_token) WHERE email_verification_token IS NOT NULL;
