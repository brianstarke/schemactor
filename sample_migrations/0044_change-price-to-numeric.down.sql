ALTER TABLE products
RENAME COLUMN price_amount TO price_cents;

ALTER TABLE products
ALTER COLUMN price_cents TYPE integer USING (price_cents * 100)::integer;
