CREATE TABLE IF NOT EXISTS Wallet (
    id SERIAL PRIMARY KEY,
    custom_id VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(18, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS Transactions (
    id               SERIAL PRIMARY KEY,
    from_wallet_id   VARCHAR(50)    NOT NULL,
    to_wallet_id     VARCHAR(50)    NOT NULL,
    amount           DECIMAL(18, 2) NOT NULL,
    transaction_time TIMESTAMP,
    FOREIGN KEY (from_wallet_id) REFERENCES Wallet (custom_id),
    FOREIGN KEY (to_wallet_id) REFERENCES Wallet (custom_id)
);