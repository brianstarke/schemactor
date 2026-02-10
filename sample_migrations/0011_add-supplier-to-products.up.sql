ALTER TABLE products
ADD COLUMN supplier_id bigint REFERENCES suppliers(id);

CREATE INDEX idx_products_supplier ON products (supplier_id);
