CREATE TABLE suppliers (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    contact_email varchar(255),
    contact_phone varchar(20),
    address text,
    created_at timestamptz DEFAULT now() NOT NULL
);
