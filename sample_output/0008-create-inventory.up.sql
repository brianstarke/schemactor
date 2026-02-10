CREATE TABLE inventory (
    id bigserial,
    product_id bigint NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    reserved_quantity integer DEFAULT 0 NOT NULL,
    warehouse_location varchar(50),
    updated_at timestamptz DEFAULT now() NOT NULL,
    low_stock_threshold integer,
    out_of_stock_notified_at timestamptz,
    PRIMARY KEY (id),
    UNIQUE (product_id),
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX idx_inventory_product ON inventory (product_id);
