ALTER TABLE users
ADD COLUMN preferences jsonb DEFAULT '{}' NOT NULL,
ADD COLUMN newsletter_subscribed boolean DEFAULT false NOT NULL;
