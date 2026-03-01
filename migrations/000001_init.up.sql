CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'client',
    is_active BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);