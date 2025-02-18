-- Add currency column to bill_items
ALTER TABLE bill_items ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';

-- Add currency-related columns to bill_item_assignments
ALTER TABLE bill_item_assignments ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';
ALTER TABLE bill_item_assignments ADD COLUMN exchange_rate REAL;
ALTER TABLE bill_item_assignments ADD COLUMN original_amount REAL;
ALTER TABLE bill_item_assignments ADD COLUMN eur_amount REAL;

-- Add total amounts in original currency and EUR to bills
ALTER TABLE bills ADD COLUMN currency TEXT NOT NULL DEFAULT 'EUR';
ALTER TABLE bills ADD COLUMN original_total REAL;
ALTER TABLE bills ADD COLUMN eur_total REAL; 