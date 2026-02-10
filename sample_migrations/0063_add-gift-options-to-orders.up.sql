ALTER TABLE orders
ADD COLUMN is_gift boolean DEFAULT false NOT NULL,
ADD COLUMN gift_message text,
ADD COLUMN gift_wrap_requested boolean DEFAULT false;
