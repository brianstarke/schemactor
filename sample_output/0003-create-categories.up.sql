CREATE TABLE categories (
    id bigserial,
    name varchar(100) NOT NULL,
    description text,
    created_at timestamptz DEFAULT now() NOT NULL,
    parent_id bigint,
    sort_order integer,
    slug varchar(100),
    image_url text,
    is_active boolean,
    PRIMARY KEY (id)
);

CREATE INDEX idx_categories_parent ON categories (parent_id);

CREATE INDEX idx_categories_slug ON categories (slug);
