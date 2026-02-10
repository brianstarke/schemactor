ALTER TABLE reviews
ADD COLUMN verified_purchase boolean DEFAULT false NOT NULL,
ADD COLUMN helpful_count integer DEFAULT 0 NOT NULL;

CREATE INDEX idx_reviews_verified ON reviews (verified_purchase);
