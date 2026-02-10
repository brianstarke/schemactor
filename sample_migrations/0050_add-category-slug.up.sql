ALTER TABLE categories
ADD COLUMN slug varchar(100) UNIQUE;

CREATE INDEX idx_categories_slug ON categories (slug);
