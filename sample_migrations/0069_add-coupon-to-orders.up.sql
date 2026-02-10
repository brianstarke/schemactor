ALTER TABLE orders
ADD COLUMN coupon_code varchar(50),
ADD COLUMN discount_amount_cents integer DEFAULT 0 NOT NULL;

CREATE INDEX idx_orders_coupon ON orders (coupon_code) WHERE coupon_code IS NOT NULL;
