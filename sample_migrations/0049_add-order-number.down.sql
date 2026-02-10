DROP INDEX IF EXISTS idx_orders_number;

ALTER TABLE orders
DROP COLUMN IF EXISTS order_number;
