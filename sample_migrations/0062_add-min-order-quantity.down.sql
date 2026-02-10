ALTER TABLE products
DROP COLUMN IF EXISTS min_order_quantity,
DROP COLUMN IF EXISTS max_order_quantity;
