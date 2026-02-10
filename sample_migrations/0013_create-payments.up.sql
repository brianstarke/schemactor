CREATE TABLE payments (
    id bigserial PRIMARY KEY,
    order_id bigint NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    method payment_method NOT NULL,
    amount_cents integer NOT NULL,
    transaction_id varchar(255),
    processed_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_payments_order ON payments (order_id);
CREATE INDEX idx_payments_transaction ON payments (transaction_id);
