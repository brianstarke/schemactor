ALTER TABLE inventory
ADD COLUMN low_stock_threshold integer DEFAULT 10,
ADD COLUMN out_of_stock_notified_at timestamptz;
