CREATE TYPE user_status AS ENUM (
    'active',
    'inactive',
    'suspended',
    'deleted',
    'banned'
);
CREATE TABLE users (
    id bigserial,
    username varchar(50) NOT NULL,
    email varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    status user_status,
    first_name varchar(100),
    last_name varchar(100),
    date_of_birth date,
    email_verified boolean,
    email_verification_token varchar(255),
    email_verified_at timestamptz,
    preferences jsonb,
    newsletter_subscribed boolean,
    last_login_at timestamptz,
    last_login_ip inet,
    timezone varchar(50),
    avatar_url text,
    deleted_at timestamptz,
    two_factor_enabled boolean,
    two_factor_secret varchar(255),
    backup_codes text[],
    PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_verification_token ON users (email_verification_token) WHERE email_verification_token IS NOT NULL;

CREATE INDEX idx_users_last_login ON users (last_login_at);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);
