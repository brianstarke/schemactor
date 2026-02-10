CREATE TABLE categories (
    id bigserial PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE,
    description text,
    created_at timestamptz DEFAULT now() NOT NULL
);
