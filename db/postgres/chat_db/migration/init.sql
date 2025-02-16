CREATE TABLE
    IF NOT EXISTS chat (
        id uuid NOT NULL PRIMARY KEY,
        sender_id uuid NOT NULL,
        receiver_id uuid NOT NULL,
        message text NOT NULL,
        sent_at timestamp NOT NULL
    );