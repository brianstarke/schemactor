CREATE TYPE notification_type AS ENUM (
    'order_update',
    'payment_received',
    'shipment_update',
    'review_posted',
    'promotion'
);


CREATE TYPE notification_priority AS ENUM (
    'low',
    'normal',
    'high',
    'urgent'
);
CREATE TABLE notifications (
    id bigserial,
    user_id bigint NOT NULL,
    type notification_type NOT NULL,
    title varchar(200) NOT NULL,
    message text NOT NULL,
    read_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    clicked_at timestamptz,
    action_url text,
    priority notification_priority,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_user ON notifications (user_id);

CREATE INDEX idx_notifications_read ON notifications (read_at);

CREATE INDEX idx_notifications_type ON notifications (type);

CREATE INDEX idx_notifications_priority ON notifications (priority);
