-- +goose Up
-- Create stock_movements table
CREATE TABLE stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    reference VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for stock_movements table
CREATE INDEX idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX idx_stock_movements_location_id ON stock_movements(location_id);
CREATE INDEX idx_stock_movements_user_id ON stock_movements(user_id);
CREATE INDEX idx_stock_movements_type ON stock_movements(type);
CREATE INDEX idx_stock_movements_created_at ON stock_movements(created_at);
CREATE INDEX idx_stock_movements_reference ON stock_movements(reference);

-- Create composite indexes for common queries
CREATE INDEX idx_stock_movements_product_created ON stock_movements(product_id, created_at DESC);
CREATE INDEX idx_stock_movements_location_created ON stock_movements(location_id, created_at DESC);
CREATE INDEX idx_stock_movements_user_created ON stock_movements(user_id, created_at DESC);

-- +goose Down
-- Drop indexes first
DROP INDEX IF EXISTS idx_stock_movements_user_created;
DROP INDEX IF EXISTS idx_stock_movements_location_created;
DROP INDEX IF EXISTS idx_stock_movements_product_created;
DROP INDEX IF EXISTS idx_stock_movements_reference;
DROP INDEX IF EXISTS idx_stock_movements_created_at;
DROP INDEX IF EXISTS idx_stock_movements_type;
DROP INDEX IF EXISTS idx_stock_movements_user_id;
DROP INDEX IF EXISTS idx_stock_movements_location_id;
DROP INDEX IF EXISTS idx_stock_movements_product_id;

-- Drop stock_movements table
DROP TABLE IF EXISTS stock_movements;