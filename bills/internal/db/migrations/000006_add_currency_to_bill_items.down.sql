-- Remove currency-related columns from bills
ALTER TABLE bills DROP COLUMN currency;
ALTER TABLE bills DROP COLUMN original_total;
ALTER TABLE bills DROP COLUMN eur_total;

-- Remove currency-related columns from bill_item_assignments
ALTER TABLE bill_item_assignments DROP COLUMN currency;
ALTER TABLE bill_item_assignments DROP COLUMN exchange_rate;
ALTER TABLE bill_item_assignments DROP COLUMN original_amount;
ALTER TABLE bill_item_assignments DROP COLUMN eur_amount;

-- Remove currency column from bill_items
ALTER TABLE bill_items DROP COLUMN currency; 