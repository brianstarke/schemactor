DROP INDEX IF EXISTS idx_categories_slug;

ALTER TABLE categories
DROP COLUMN IF EXISTS slug;
