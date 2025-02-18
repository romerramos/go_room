CREATE TABLE IF NOT EXISTS bills (
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