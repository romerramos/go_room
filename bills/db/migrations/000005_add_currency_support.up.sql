-- Create exchange rates table
CREATE TABLE IF NOT EXISTS exchange_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    currency_from TEXT NOT NULL,
    currency_to TEXT NOT NULL,
    rate REAL NOT NULL,
    created_at DATETIME NOT NULL
);

-- Add currency to bill_items
ALTER TABLE bill_items ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';

-- Add currency info to bill_item_assignments
ALTER TABLE bill_item_assignments ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';
ALTER TABLE bill_item_assignments ADD COLUMN exchange_rate REAL NOT NULL DEFAULT 1.0;
ALTER TABLE bill_item_assignments ADD COLUMN original_amount REAL NOT NULL DEFAULT 0.0;
ALTER TABLE bill_item_assignments ADD COLUMN eur_amount REAL NOT NULL DEFAULT 0.0;

-- Add currency info to bills
ALTER TABLE bills ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';
ALTER TABLE bills ADD COLUMN original_total REAL NOT NULL DEFAULT 0.0;
ALTER TABLE bills ADD COLUMN eur_total REAL NOT NULL DEFAULT 0.0; 