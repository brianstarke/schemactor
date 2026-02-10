DROP INDEX IF EXISTS idx_payments_status;

ALTER TABLE payments
DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS payment_status;
