ALTER TABLE reviews
ADD COLUMN seller_response text,
ADD COLUMN seller_responded_at timestamptz;
