CREATE TYPE order_status AS ENUM (
    'pending',
    'confirmed',
    'shipped',
    'delivered',
    'cancelled',
    'refunded',
    'processing'
);
CREATE TABLE orders (
    id bigserial,
    user_id bigint NOT NULL,
    status order_status DEFAULT 'pending' NOT NULL,
    total_cents integer NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    shipping_address_line1 varchar(255),
    shipping_address_line2 varchar(255),
    shipping_city varchar(100),
    shipping_state varchar(100),
    shipping_postal_code varchar(20),
    shipping_country varchar(2),
    customer_notes text,
    internal_notes text,
    subtotal_cents integer,
    tax_cents integer,
    shipping_cents integer,
    order_number varchar(50),
    currency currency,
    is_gift boolean,
    gift_message text,
    gift_wrap_requested boolean,
    coupon_code varchar(50),
    discount_amount_cents integer,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_user ON orders (user_id);

CREATE INDEX idx_orders_status ON orders (status);

CREATE INDEX idx_orders_created ON orders (created_at);

CREATE INDEX idx_orders_number ON orders (order_number);

CREATE INDEX idx_orders_coupon ON orders (coupon_code) WHERE coupon_code IS NOT NULL;
