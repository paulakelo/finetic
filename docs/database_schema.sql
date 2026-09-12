-- STREAMING_CHUNK: Enabling UUID generation for secure IDs...
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUIF PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- STREAMING_CHUNK: Creating category and transaction ledgers...
-- Transaction Categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Transactions Table
CREATE TABLE transactions
(
    id                   UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id              UUID REFERENCES users (id) ON DELETE CASCADE,
    category_id          INT                      REFERENCES categories (id) ON DELETE SET NULL,
    amount               BIGINT                   NOT NULL, -- stored in lowest denomination (e.g., cents)
    type                 VARCHAR(20)              NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    mpesa_receipt_number VARCHAR(50) UNIQUE,
    sender_or_recipient  VARCHAR(255),
    raw_sms_data         TEXT,
    transaction_date     TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

-- STREAMING_CHUNK: Creating AI insights table and optimizing queries with indexes...
-- AI Insights Table
CREATE TABLE ai_insights (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    insight_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);