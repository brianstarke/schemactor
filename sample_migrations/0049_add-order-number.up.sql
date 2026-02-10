ALTER TABLE orders
ADD COLUMN order_number varchar(50) UNIQUE;

CREATE INDEX idx_orders_number ON orders (order_number);
