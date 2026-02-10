ALTER TABLE inventory
DROP COLUMN IF EXISTS low_stock_threshold,
DROP COLUMN IF EXISTS out_of_stock_notified_at;
