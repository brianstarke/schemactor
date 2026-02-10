CREATE TYPE notification_priority AS ENUM (
    'low',
    'normal',
    'high',
    'urgent'
);

ALTER TABLE notifications
ADD COLUMN priority notification_priority DEFAULT 'normal' NOT NULL;

CREATE INDEX idx_notifications_priority ON notifications (priority);
