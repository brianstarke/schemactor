ALTER TABLE orders
DROP COLUMN IF EXISTS shipping_address_line1,
DROP COLUMN IF EXISTS shipping_address_line2,
DROP COLUMN IF EXISTS shipping_city,
DROP COLUMN IF EXISTS shipping_state,
DROP COLUMN IF EXISTS shipping_postal_code,
DROP COLUMN IF EXISTS shipping_country;
