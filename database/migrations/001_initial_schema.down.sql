-- Drop triggers
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_caregivers_updated_at ON caregivers;
DROP TRIGGER IF EXISTS update_kids_updated_at ON kids;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_kid_type;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_kid_id;

DROP INDEX IF EXISTS idx_caregivers_created_at;
DROP INDEX IF EXISTS idx_caregivers_relationship;
DROP INDEX IF EXISTS idx_caregivers_email;
DROP INDEX IF EXISTS idx_caregivers_name;

DROP INDEX IF EXISTS idx_kids_created_at;
DROP INDEX IF EXISTS idx_kids_age;
DROP INDEX IF EXISTS idx_kids_name;

-- Drop tables (in reverse order due to foreign key constraints)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS caregivers;
DROP TABLE IF EXISTS kids;

-- Drop enum types
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS relationship_type;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";