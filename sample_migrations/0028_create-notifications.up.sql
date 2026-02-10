CREATE TABLE notifications (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type notification_type NOT NULL,
    title varchar(200) NOT NULL,
    message text NOT NULL,
    read_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_notifications_user ON notifications (user_id);
CREATE INDEX idx_notifications_read ON notifications (read_at);
CREATE INDEX idx_notifications_type ON notifications (type);
