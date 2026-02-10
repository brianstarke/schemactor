DROP INDEX IF EXISTS idx_products_featured;

ALTER TABLE products
DROP COLUMN IF EXISTS is_featured,
DROP COLUMN IF EXISTS featured_priority;
