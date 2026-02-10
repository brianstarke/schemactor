DROP INDEX IF EXISTS idx_products_rating;

ALTER TABLE products
DROP COLUMN IF EXISTS average_rating,
DROP COLUMN IF EXISTS review_count;
