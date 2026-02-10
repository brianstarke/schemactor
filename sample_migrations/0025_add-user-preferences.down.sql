ALTER TABLE users
DROP COLUMN IF EXISTS preferences,
DROP COLUMN IF EXISTS newsletter_subscribed;
