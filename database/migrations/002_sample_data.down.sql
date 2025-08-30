-- Remove sample data (in reverse order due to foreign key constraints)

-- Clear transactions
DELETE FROM transactions;

-- Clear caregivers
DELETE FROM caregivers;

-- Clear kids
DELETE FROM kids;

-- Reset sequences to start from 1 again
ALTER SEQUENCE transactions_id_seq RESTART WITH 1;
ALTER SEQUENCE caregivers_id_seq RESTART WITH 1;
ALTER SEQUENCE kids_id_seq RESTART WITH 1;