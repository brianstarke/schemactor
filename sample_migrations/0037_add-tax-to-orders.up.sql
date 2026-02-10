ALTER TABLE orders
ADD COLUMN subtotal_cents integer,
ADD COLUMN tax_cents integer,
ADD COLUMN shipping_cents integer;
