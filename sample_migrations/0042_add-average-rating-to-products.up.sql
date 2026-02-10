ALTER TABLE products
ADD COLUMN average_rating numeric(3,2),
ADD COLUMN review_count integer DEFAULT 0 NOT NULL;

CREATE INDEX idx_products_rating ON products (average_rating DESC NULLS LAST);
