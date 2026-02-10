DROP INDEX IF EXISTS idx_notifications_priority;

ALTER TABLE notifications
DROP COLUMN IF EXISTS priority;

DROP TYPE IF EXISTS notification_priority;
