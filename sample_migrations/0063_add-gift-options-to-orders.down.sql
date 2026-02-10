ALTER TABLE orders
DROP COLUMN IF EXISTS is_gift,
DROP COLUMN IF EXISTS gift_message,
DROP COLUMN IF EXISTS gift_wrap_requested;
