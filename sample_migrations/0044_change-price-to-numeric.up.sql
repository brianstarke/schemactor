-- Change from integer cents to numeric for better precision
ALTER TABLE products
ALTER COLUMN price_cents TYPE numeric(10,2) USING price_cents::numeric / 100;

ALTER TABLE products
RENAME COLUMN price_cents TO price_amount;
