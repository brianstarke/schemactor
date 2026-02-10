ALTER TABLE products
ADD COLUMN is_active boolean DEFAULT true NOT NULL;

CREATE INDEX idx_products_active ON products (is_active);
