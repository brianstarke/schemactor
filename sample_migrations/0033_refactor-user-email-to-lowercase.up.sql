-- Migrate existing emails to lowercase
-- In production, you'd run: UPDATE users SET email = LOWER(email);

-- Add constraint to ensure emails are lowercase
ALTER TABLE users
ADD CONSTRAINT users_email_lowercase CHECK (email = LOWER(email));
