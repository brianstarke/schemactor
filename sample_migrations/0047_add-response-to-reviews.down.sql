ALTER TABLE reviews
DROP COLUMN IF EXISTS seller_response,
DROP COLUMN IF EXISTS seller_responded_at;
