DROP INDEX IF EXISTS idx_products_supplier;

ALTER TABLE products
DROP COLUMN IF EXISTS supplier_id;
