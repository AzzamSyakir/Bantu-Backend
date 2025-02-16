CREATE TABLE
    IF NOT EXISTS transactions (
        id uuid NOT NULL PRIMARY KEY,
        transaction_type text NOT NULL, -- eg: 'deposit', 'payment', 'withdrawal'
        amount numeric NOT NULL,
        status text NOT NULL, -- eg: 'pending', 'completed', 'cancelled'
        user_id uuid NOT NULL,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL,
    );