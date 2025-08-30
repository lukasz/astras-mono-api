-- Initial schema migration for Astras API
-- Creates tables for Kids, Caregivers, and Transactions

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE relationship_type AS ENUM ('parent', 'guardian', 'grandparent', 'relative', 'caregiver');
CREATE TYPE transaction_type AS ENUM ('earn', 'spend');

-- Kids table
CREATE TABLE kids (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK (length(trim(name)) >= 2),
    birthdate DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Caregivers table
CREATE TABLE caregivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK (length(trim(name)) >= 2),
    email VARCHAR(255) NOT NULL UNIQUE,
    relationship relationship_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Star transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    kid_id INTEGER NOT NULL REFERENCES kids(id) ON DELETE CASCADE,
    type transaction_type NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 1 AND amount <= 100),
    description VARCHAR(255) NOT NULL CHECK (length(trim(description)) > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX idx_kids_name ON kids(name);
CREATE INDEX idx_kids_birthdate ON kids(birthdate);
CREATE INDEX idx_kids_created_at ON kids(created_at);

CREATE INDEX idx_caregivers_name ON caregivers(name);
CREATE INDEX idx_caregivers_email ON caregivers(email);
CREATE INDEX idx_caregivers_relationship ON caregivers(relationship);
CREATE INDEX idx_caregivers_created_at ON caregivers(created_at);

CREATE INDEX idx_transactions_kid_id ON transactions(kid_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_kid_type ON transactions(kid_id, type);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at
CREATE TRIGGER update_kids_updated_at 
    BEFORE UPDATE ON kids 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_caregivers_updated_at 
    BEFORE UPDATE ON caregivers 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at 
    BEFORE UPDATE ON transactions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();