ALTER TABLE products
ADD COLUMN min_order_quantity integer DEFAULT 1 NOT NULL,
ADD COLUMN max_order_quantity integer;
