ALTER TABLE users
ADD COLUMN two_factor_enabled boolean DEFAULT false NOT NULL,
ADD COLUMN two_factor_secret varchar(255),
ADD COLUMN backup_codes text[];
