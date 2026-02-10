ALTER TABLE orders
ADD COLUMN shipping_address_line1 varchar(255),
ADD COLUMN shipping_address_line2 varchar(255),
ADD COLUMN shipping_city varchar(100),
ADD COLUMN shipping_state varchar(100),
ADD COLUMN shipping_postal_code varchar(20),
ADD COLUMN shipping_country varchar(2);
