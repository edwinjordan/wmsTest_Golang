-- +goose Up
-- Create products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(15,2) DEFAULT 0.00 NOT NULL,
    weight DECIMAL(8,2) DEFAULT 0.00 NOT NULL,
    dimensions VARCHAR(50),
    category VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    quantity INTEGER NOT NULL ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for products table
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_created_at ON products(created_at);

-- Create full-text search index for product search
CREATE INDEX idx_products_search ON products USING gin(to_tsvector('english', name || ' ' || sku || ' ' || category));

-- +goose Down
-- Drop indexes first
DROP INDEX IF EXISTS idx_products_search;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_sku;

-- Drop products table
DROP TABLE IF EXISTS products;