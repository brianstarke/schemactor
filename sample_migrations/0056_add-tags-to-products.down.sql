DROP INDEX IF EXISTS idx_products_tags;

ALTER TABLE products
DROP COLUMN IF EXISTS tags;
