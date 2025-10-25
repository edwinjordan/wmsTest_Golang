-- +goose Up
-- Insert default admin user
-- Password: admin123 (bcrypt hashed)
INSERT INTO users (username, email, password, api_key, is_active) VALUES 
('admin', 'admin@wms.local', '$2a$10$rJ1P.3oKV1z8mZ7KpYcPt.mz5.P4U4bYz6I3F5l1rP.3oKV1z8mZ7K', 'wms_admin_default_api_key_change_in_production', true);

-- Insert sample products
INSERT INTO products (sku, name, description, price, weight, dimensions, category, quantity) VALUES 
('LAPTOP-001', 'MacBook Pro 14"', 'Apple MacBook Pro 14-inch with M2 chip', 1999.99, 1.6, '31.26x22.12x1.55', 'Electronics', 10),
('MOUSE-001', 'Wireless Mouse', 'Logitech MX Master 3 Wireless Mouse', 99.99, 0.14, '12.6x8.4x5.1', 'Electronics', 25),
('BOOK-001', 'Go Programming', 'The Go Programming Language Book', 49.99, 0.5, '23.5x19.1x2.5', 'Books', 100),
('CHAIR-001', 'Office Chair', 'Ergonomic Office Chair with Lumbar Support', 299.99, 15.0, '66x66x117', 'Furniture', 5);

-- Insert sample locations
INSERT INTO locations (code, name, zone, aisle, rack, shelf, capacity) VALUES 
('A-01-01-01', 'Zone A Aisle 1 Rack 1 Shelf 1', 'A', '01', '01', '01', 100),
('A-01-01-02', 'Zone A Aisle 1 Rack 1 Shelf 2', 'A', '01', '01', '02', 100),
('A-01-02-01', 'Zone A Aisle 1 Rack 2 Shelf 1', 'A', '01', '02', '01', 150),
('B-01-01-01', 'Zone B Aisle 1 Rack 1 Shelf 1', 'B', '01', '01', '01', 200),
('B-02-01-01', 'Zone B Aisle 2 Rack 1 Shelf 1', 'B', '02', '01', '01', 200);

-- Insert sample stock movements (history)
INSERT INTO stock_movements (product_id, location_id, user_id, type, quantity, reference, notes) VALUES 
(1, 1, 1, 'IN', 25, 'PO-2024-001', 'Initial stock from supplier'),
(1, 2, 1, 'IN', 15, 'PO-2024-001', 'Initial stock from supplier'),
(2, 1, 1, 'IN', 50, 'PO-2024-002', 'Mouse bulk purchase'),
(3, 3, 1, 'IN', 100, 'PO-2024-003', 'Book inventory restocking'),
(4, 4, 1, 'IN', 10, 'PO-2024-004', 'Office furniture delivery');

-- +goose Down
-- Delete sample data in reverse order (due to foreign key constraints)
DELETE FROM stock_movements WHERE user_id = 1;
DELETE FROM stocks WHERE id IN (1, 2, 3, 4, 5);
DELETE FROM locations WHERE code LIKE 'A-%' OR code LIKE 'B-%';
DELETE FROM products WHERE sku IN ('LAPTOP-001', 'MOUSE-001', 'BOOK-001', 'CHAIR-001');
DELETE FROM users WHERE username = 'admin';