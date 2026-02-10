DROP VIEW IF EXISTS products_with_stock;

CREATE OR REPLACE VIEW products_with_stock AS
SELECT 
    p.*,
    COALESCE(i.quantity, 0) as stock_quantity,
    COALESCE(i.reserved_quantity, 0) as reserved_quantity,
    COALESCE(i.quantity - i.reserved_quantity, 0) as available_quantity
FROM products p
LEFT JOIN inventory i ON p.id = i.product_id;
