CREATE TYPE payment_method AS ENUM (
    'credit_card',
    'debit_card',
    'paypal',
    'bank_transfer',
    'cryptocurrency'
);


CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'refunded'
);
CREATE TABLE payments (
    id bigserial,
    order_id bigint NOT NULL,
    method payment_method NOT NULL,
    amount_cents integer NOT NULL,
    transaction_id varchar(255),
    processed_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    refunded_amount_cents integer,
    refunded_at timestamptz,
    status payment_status,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);

CREATE INDEX idx_payments_order ON payments (order_id);

CREATE INDEX idx_payments_transaction ON payments (transaction_id);

CREATE INDEX idx_payments_status ON payments (status);
