ALTER TABLE payments
DROP COLUMN IF EXISTS refunded_amount_cents,
DROP COLUMN IF EXISTS refunded_at;
