CREATE TABLE products (
    id bigserial,
    name varchar(255) NOT NULL,
    description text,
    price_cents numeric(10,2) NOT NULL,
    category_id bigint,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    sku varchar(50),
    supplier_id bigint,
    primary_image_url text,
    image_urls text[],
    discount_percentage integer,
    weight_grams integer,
    dimensions_cm varchar(50),
    is_featured boolean,
    featured_priority integer,
    slug varchar(255),
    is_active boolean,
    average_rating numeric(3,2),
    review_count integer,
    tags text[],
    warranty_months integer,
    warranty_description text,
    min_order_quantity integer,
    max_order_quantity integer,
    search_vector tsvector,
    currency currency,
    PRIMARY KEY (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE INDEX idx_products_category ON products (category_id);

CREATE INDEX idx_products_price ON products (price_cents);

CREATE INDEX idx_products_sku ON products (sku) WHERE sku IS NOT NULL;

CREATE INDEX idx_products_supplier ON products (supplier_id);

CREATE INDEX idx_products_featured ON products (is_featured, featured_priority) WHERE is_featured = true;

CREATE INDEX idx_products_slug ON products (slug);

CREATE INDEX idx_products_active ON products (is_active);

CREATE INDEX idx_products_rating ON products (average_rating);

CREATE INDEX idx_products_tags ON products (tags);

CREATE INDEX idx_products_search ON products (search_vector);
COMMENT ON COLUMN products.search_vector IS 'Full-text search vector for product name and description';
