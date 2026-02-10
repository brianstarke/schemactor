ALTER TABLE notifications
DROP COLUMN IF EXISTS clicked_at,
DROP COLUMN IF EXISTS action_url;
