ALTER TABLE products
DROP COLUMN IF EXISTS primary_image_url,
DROP COLUMN IF EXISTS image_urls;
