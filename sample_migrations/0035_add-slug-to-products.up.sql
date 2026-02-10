ALTER TABLE products
ADD COLUMN slug varchar(255) UNIQUE;

CREATE INDEX idx_products_slug ON products (slug);
