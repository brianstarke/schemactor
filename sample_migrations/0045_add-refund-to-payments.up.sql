ALTER TABLE payments
ADD COLUMN refunded_amount_cents integer DEFAULT 0 NOT NULL,
ADD COLUMN refunded_at timestamptz;
