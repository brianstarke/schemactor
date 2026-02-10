ALTER TABLE categories
ADD COLUMN image_url text,
ADD COLUMN is_active boolean DEFAULT true NOT NULL;
