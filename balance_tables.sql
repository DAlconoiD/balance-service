CREATE TABLE accounts (
    account_id INT PRIMARY KEY,
    balance NUMERIC(18, 2) CONSTRAINT non_negative_balance CHECK (balance >= 0) NOT NULL
);

CREATE TABLE transactions (
    transaction_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id INT REFERENCES accounts ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delta NUMERIC(18, 2) NOT NULL,
    remaining NUMERIC(18, 2) NOT NULL,
    message TEXT NOT NULL
);