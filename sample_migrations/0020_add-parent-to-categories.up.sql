ALTER TABLE categories
ADD COLUMN parent_id bigint REFERENCES categories(id) ON DELETE CASCADE,
ADD COLUMN sort_order integer DEFAULT 0;

CREATE INDEX idx_categories_parent ON categories (parent_id);
