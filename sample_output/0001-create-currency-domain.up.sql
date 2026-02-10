CREATE DOMAIN currency AS varchar(3) DEFAULT 'USD' CHECK (value ~ '^[A-Z]{3}$');
