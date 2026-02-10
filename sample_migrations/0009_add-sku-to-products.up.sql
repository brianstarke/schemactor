ALTER TABLE products
ADD COLUMN sku varchar(50) UNIQUE;

CREATE INDEX idx_products_sku ON products (sku) WHERE sku IS NOT NULL;
