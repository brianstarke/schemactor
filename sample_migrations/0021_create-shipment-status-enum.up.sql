CREATE TYPE shipment_status AS ENUM (
    'preparing',
    'in_transit',
    'delivered',
    'returned'
);
