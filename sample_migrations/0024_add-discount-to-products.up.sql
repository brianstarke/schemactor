ALTER TABLE products
ADD COLUMN discount_percentage integer DEFAULT 0 CHECK (discount_percentage >= 0 AND discount_percentage <= 100);
