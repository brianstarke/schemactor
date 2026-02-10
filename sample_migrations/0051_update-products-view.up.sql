DROP VIEW IF EXISTS products_with_stock;

CREATE OR REPLACE VIEW products_with_stock AS
SELECT 
    p.*,
    c.name as category_name,
    s.name as supplier_name,
    COALESCE(i.quantity, 0) as stock_quantity,
    COALESCE(i.reserved_quantity, 0) as reserved_quantity,
    COALESCE(i.quantity - i.reserved_quantity, 0) as available_quantity
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN suppliers s ON p.supplier_id = s.id
LEFT JOIN inventory i ON p.id = i.product_id
WHERE p.is_active = true;
