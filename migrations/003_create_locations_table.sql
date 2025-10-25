-- +goose Up
-- Create locations table
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    zone VARCHAR(10) NOT NULL,
    aisle VARCHAR(10) NOT NULL,
    rack VARCHAR(10) NOT NULL,
    shelf VARCHAR(10) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    temperature DECIMAL(5,2),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for locations table
CREATE INDEX idx_locations_code ON locations(code);
CREATE INDEX idx_locations_zone ON locations(zone);
CREATE INDEX idx_locations_aisle ON locations(aisle);
CREATE INDEX idx_locations_rack ON locations(rack);
CREATE INDEX idx_locations_shelf ON locations(shelf);
CREATE INDEX idx_locations_is_active ON locations(is_active);
CREATE INDEX idx_locations_zone_aisle_rack_shelf ON locations(zone, aisle, rack, shelf);

-- +goose Down
-- Drop indexes first
DROP INDEX IF EXISTS idx_locations_zone_aisle_rack_shelf;
DROP INDEX IF EXISTS idx_locations_is_active;
DROP INDEX IF EXISTS idx_locations_shelf;
DROP INDEX IF EXISTS idx_locations_rack;
DROP INDEX IF EXISTS idx_locations_aisle;
DROP INDEX IF EXISTS idx_locations_zone;
DROP INDEX IF EXISTS idx_locations_code;

-- Drop locations table
DROP TABLE IF EXISTS locations;