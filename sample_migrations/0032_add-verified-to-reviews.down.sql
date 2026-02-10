DROP INDEX IF EXISTS idx_reviews_verified;

ALTER TABLE reviews
DROP COLUMN IF EXISTS verified_purchase,
DROP COLUMN IF EXISTS helpful_count;
