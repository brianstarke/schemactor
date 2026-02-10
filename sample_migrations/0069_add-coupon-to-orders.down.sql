DROP INDEX IF EXISTS idx_orders_coupon;

ALTER TABLE orders
DROP COLUMN IF EXISTS coupon_code,
DROP COLUMN IF EXISTS discount_amount_cents;
