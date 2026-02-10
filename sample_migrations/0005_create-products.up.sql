CREATE TABLE products (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    price_cents integer NOT NULL,
    category_id bigint REFERENCES categories(id),
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_products_category ON products (category_id);
CREATE INDEX idx_products_price ON products (price_cents);
