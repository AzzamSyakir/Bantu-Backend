CREATE TABLE
    IF NOT EXISTS jobs (
        id uuid NOT NULL PRIMARY KEY,
        company_id uuid, --harus ada yng di input dari salah satu field ini 
        user uuid, --harus ada yng di input dari salah satu field ini 
        title text NOT NULL,
        description text NOT NULL,
        budget numeric NOT NULL,
        deadline timestamp NOT NULL,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL,
    );