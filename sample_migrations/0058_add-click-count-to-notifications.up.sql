ALTER TABLE notifications
ADD COLUMN clicked_at timestamptz,
ADD COLUMN action_url text;
