-- Sample data for development and testing

-- Insert sample kids
INSERT INTO kids (name, birthdate) VALUES 
    ('Alice Johnson', '2015-03-15'),
    ('Bob Smith', '2012-07-22'),
    ('Emma Wilson', '2017-11-08'),
    ('Charlie Brown', '2013-09-10'),
    ('Sophia Davis', '2016-05-28');

-- Insert sample caregivers
INSERT INTO caregivers (name, email, relationship) VALUES 
    ('Sarah Johnson', 'sarah.johnson@example.com', 'parent'),
    ('Mike Smith', 'mike.smith@example.com', 'guardian'),
    ('Grace Wilson', 'grace.wilson@example.com', 'grandparent'),
    ('David Brown', 'david.brown@example.com', 'parent'),
    ('Lisa Davis', 'lisa.davis@example.com', 'relative');

-- Insert sample transactions
INSERT INTO transactions (kid_id, type, amount, description) VALUES 
    -- Alice's transactions
    (1, 'earn', 5, 'Completed homework perfectly'),
    (1, 'earn', 10, 'Cleaned room thoroughly'),
    (1, 'spend', 3, 'Bought sticker reward'),
    (1, 'earn', 8, 'Helped with dishes for a week'),
    
    -- Bob's transactions
    (2, 'earn', 15, 'Perfect behavior for a week'),
    (2, 'spend', 5, 'Bought comic book'),
    (2, 'earn', 12, 'Completed extra math exercises'),
    (2, 'spend', 8, 'Movie night reward'),
    
    -- Emma's transactions
    (3, 'earn', 2, 'Put toys away nicely'),
    (3, 'earn', 3, 'Brushed teeth without reminder'),
    (3, 'earn', 5, 'Shared toys with friend'),
    
    -- Charlie's transactions
    (4, 'earn', 20, 'Outstanding school report'),
    (4, 'spend', 10, 'Video game time'),
    (4, 'earn', 7, 'Helped neighbor with groceries'),
    
    -- Sophia's transactions
    (5, 'earn', 4, 'Read book to younger sibling'),
    (5, 'earn', 6, 'Kept room tidy for 3 days'),
    (5, 'spend', 2, 'Small toy reward');