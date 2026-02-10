CREATE TABLE suppliers (
    id bigserial,
    name varchar(255) NOT NULL,
    contact_email varchar(255),
    contact_phone varchar(20),
    address text,
    created_at timestamptz DEFAULT now() NOT NULL,
    website varchar(255),
    tax_id varchar(50),
    payment_terms varchar(100),
    rating numeric(3,2),
    total_orders integer,
    PRIMARY KEY (id)
);
