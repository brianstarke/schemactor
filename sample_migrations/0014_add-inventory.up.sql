CREATE TABLE inventory (
    id bigserial PRIMARY KEY,
    product_id bigint NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity integer NOT NULL DEFAULT 0,
    reserved_quantity integer NOT NULL DEFAULT 0,
    warehouse_location varchar(50),
    updated_at timestamptz DEFAULT now() NOT NULL,
    UNIQUE (product_id)
);

CREATE INDEX idx_inventory_product ON inventory (product_id);
