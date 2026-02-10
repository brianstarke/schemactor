ALTER TABLE orders
DROP COLUMN IF EXISTS subtotal_cents,
DROP COLUMN IF EXISTS tax_cents,
DROP COLUMN IF EXISTS shipping_cents;
