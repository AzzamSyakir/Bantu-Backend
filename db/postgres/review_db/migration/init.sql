CREATE TABLE
  IF NOT EXISTS review (
    id uuid NOT NULL PRIMARY KEY,
    reviewer_id uuid NOT NULL,
    reviewee_id uuid NOT NULL,
    rating int NOT NULL, -- on scale 1 to 5 
    comment text NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
  );