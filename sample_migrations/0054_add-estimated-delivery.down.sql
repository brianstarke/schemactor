ALTER TABLE shipments
DROP COLUMN IF EXISTS estimated_delivery_date,
DROP COLUMN IF EXISTS actual_delivery_date;
