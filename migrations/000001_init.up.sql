CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'client',
    is_active BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS banks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS enterprises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_id INT NOT NULL REFERENCES banks(id) ON DELETE RESTRICT,
    account_number VARCHAR(20) NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);

CREATE TABLE IF NOT EXISTS deposits (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_id INT NOT NULL REFERENCES banks(id) ON DELETE RESTRICT,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    interest_rate NUMERIC(5, 2) NOT NULL, 
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deposits_user_id ON deposits(user_id);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INT REFERENCES accounts(id) ON DELETE SET NULL, 
    from_deposit_id INT REFERENCES deposits(id) ON DELETE SET NULL,

    to_account_id INT REFERENCES accounts(id) ON DELETE SET NULL,  
    to_deposit_id INT REFERENCES deposits(id) ON DELETE SET NULL,
    
    amount NUMERIC(15, 2) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, 
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_from_acc ON transactions(from_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_acc ON transactions(to_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_dep ON transactions(from_deposit_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_dep ON transactions(to_deposit_id);

CREATE TABLE IF NOT EXISTS enterprise_employees (
    enterprise_id INT NOT NULL REFERENCES enterprises(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (enterprise_id, user_id)
);

CREATE TABLE IF NOT EXISTS salary_applications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enterprise_id INT NOT NULL REFERENCES enterprises(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS action_logs (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id), 
    action_type VARCHAR(100) NOT NULL, 
    details JSONB,
    is_undone BOOLEAN NOT NULL DEFAULT false, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);