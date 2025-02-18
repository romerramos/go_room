-- Drop currency columns from bills
PRAGMA foreign_keys=off;
BEGIN TRANSACTION;

CREATE TABLE bills_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    due_date DATETIME NOT NULL,
    paid BOOLEAN DEFAULT FALSE,
    issuer_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (issuer_id) REFERENCES issuers(id),
    FOREIGN KEY (receiver_id) REFERENCES receivers(id)
);

INSERT INTO bills_backup 
SELECT id, due_date, paid, issuer_id, receiver_id, created_at, updated_at
FROM bills;

DROP TABLE bills;

ALTER TABLE bills_backup RENAME TO bills;

-- Drop currency columns from bill_item_assignments
CREATE TABLE bill_item_assignments_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bill_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price REAL NOT NULL,
    subtotal REAL NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES bill_items(id) ON DELETE RESTRICT
);

INSERT INTO bill_item_assignments_backup 
SELECT id, bill_id, item_id, quantity, unit_price, subtotal, created_at, updated_at
FROM bill_item_assignments;

DROP TABLE bill_item_assignments;

ALTER TABLE bill_item_assignments_backup RENAME TO bill_item_assignments;

-- Drop currency from bill_items
CREATE TABLE bill_items_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

INSERT INTO bill_items_backup 
SELECT id, name, price, created_at, updated_at
FROM bill_items;

DROP TABLE bill_items;

ALTER TABLE bill_items_backup RENAME TO bill_items;

COMMIT;
PRAGMA foreign_keys=on;

-- Drop exchange rates table
DROP TABLE IF EXISTS exchange_rates; 