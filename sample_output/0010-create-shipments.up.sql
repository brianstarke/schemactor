CREATE TYPE shipment_status AS ENUM (
    'preparing',
    'in_transit',
    'delivered',
    'returned',
    'expedited'
);
CREATE TABLE shipments (
    id bigserial,
    order_id bigint NOT NULL,
    tracking_number varchar(100),
    carrier varchar(100),
    status shipment_status DEFAULT 'preparing' NOT NULL,
    shipped_at timestamptz,
    delivered_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    tracking_events jsonb,
    estimated_delivery_date date,
    actual_delivery_date date,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE INDEX idx_shipments_order ON shipments (order_id);

CREATE INDEX idx_shipments_tracking ON shipments (tracking_number);

CREATE INDEX idx_shipments_status ON shipments (status);
