ALTER TABLE products
ADD COLUMN is_featured boolean DEFAULT false NOT NULL,
ADD COLUMN featured_priority integer;

CREATE INDEX idx_products_featured ON products (is_featured, featured_priority) WHERE is_featured = true;
