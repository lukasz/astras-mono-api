-- Astras API Database Schema
-- PostgreSQL database schema for Kid, Caregiver, and Star Transaction management

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

-- Sample data for development/testing
INSERT INTO kids (name, birthdate) VALUES 
    ('Alice Johnson', '2015-03-15'),
    ('Bob Smith', '2012-07-22'),
    ('Emma Wilson', '2017-11-08');

INSERT INTO caregivers (name, email, relationship) VALUES 
    ('Sarah Johnson', 'sarah.johnson@example.com', 'parent'),
    ('Mike Smith', 'mike.smith@example.com', 'guardian'),
    ('Grace Wilson', 'grace.wilson@example.com', 'grandparent');

INSERT INTO transactions (kid_id, type, amount, description) VALUES 
    (1, 'earn', 5, 'Completed homework perfectly'),
    (2, 'spend', 3, 'Bought sticker reward'),
    (1, 'earn', 10, 'Cleaned room thoroughly'),
    (3, 'earn', 2, 'Helped with dishes'),
    (2, 'earn', 15, 'Perfect behavior for a week');