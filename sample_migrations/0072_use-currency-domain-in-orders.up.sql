-- Change orders.currency from varchar(3) to currency domain
ALTER TABLE orders
ALTER COLUMN currency TYPE currency USING currency::currency;
