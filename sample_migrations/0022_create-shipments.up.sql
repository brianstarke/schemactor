CREATE TABLE shipments (
    id bigserial PRIMARY KEY,
    order_id bigint NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    tracking_number varchar(100),
    carrier varchar(100),
    status shipment_status DEFAULT 'preparing' NOT NULL,
    shipped_at timestamptz,
    delivered_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_shipments_order ON shipments (order_id);
CREATE INDEX idx_shipments_tracking ON shipments (tracking_number);
CREATE INDEX idx_shipments_status ON shipments (status);
