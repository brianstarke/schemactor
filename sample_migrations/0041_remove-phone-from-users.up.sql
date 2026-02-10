-- Decided to move phone to separate table for better validation
ALTER TABLE users
DROP COLUMN IF EXISTS phone;
