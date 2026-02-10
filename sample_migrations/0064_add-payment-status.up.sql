CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'refunded'
);

ALTER TABLE payments
ADD COLUMN status payment_status DEFAULT 'pending' NOT NULL;

CREATE INDEX idx_payments_status ON payments (status);
