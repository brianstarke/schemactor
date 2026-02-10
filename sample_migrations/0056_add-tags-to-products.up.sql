ALTER TABLE products
ADD COLUMN tags text[];

CREATE INDEX idx_products_tags ON products USING gin(tags);
