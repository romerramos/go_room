CREATE TABLE exchange_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    currency_from TEXT NOT NULL,
    currency_to TEXT NOT NULL,
    rate REAL NOT NULL,
    created_at DATETIME NOT NULL
); 