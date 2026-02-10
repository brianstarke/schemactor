DROP INDEX IF EXISTS idx_categories_parent;

ALTER TABLE categories
DROP COLUMN IF EXISTS parent_id,
DROP COLUMN IF EXISTS sort_order;
