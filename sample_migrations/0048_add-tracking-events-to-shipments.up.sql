ALTER TABLE shipments
ADD COLUMN tracking_events jsonb DEFAULT '[]' NOT NULL;
