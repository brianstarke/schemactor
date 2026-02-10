CREATE DOMAIN currency AS varchar(3)
    DEFAULT 'USD'
    CHECK (value ~ '^[A-Z]{3}$');

COMMENT ON DOMAIN currency IS 'ISO 4217 three-letter currency code';
