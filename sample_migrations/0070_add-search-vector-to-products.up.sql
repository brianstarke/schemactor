ALTER TABLE products
ADD COLUMN search_vector tsvector;

CREATE INDEX idx_products_search ON products USING gin(search_vector);

COMMENT ON COLUMN products.search_vector IS 'Full-text search vector for product name and description';
