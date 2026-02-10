ALTER TABLE suppliers
ADD COLUMN rating numeric(3,2),
ADD COLUMN total_orders integer DEFAULT 0 NOT NULL;
