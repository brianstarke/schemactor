ALTER TABLE orders
DROP COLUMN IF EXISTS customer_notes,
DROP COLUMN IF EXISTS internal_notes;
