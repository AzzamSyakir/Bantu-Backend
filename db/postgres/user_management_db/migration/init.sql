create table
CREATE TABLE
    IF NOT EXISTS users (
        id uuid NOT NULL PRIMARY KEY,
        name text NOT NULL UNIQUE,
        password text NOT NULL,
        email text UNIQUE, --harus ada yng di input dari salah satu field ini 
        phone text UNIQUE, --harus ada yng di input dari salah satu field ini 
        balance int NOT NULL DEFAULT 0,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL,
    );

CREATE TABLE
    IF NOT EXISTS admin (
        id uuid NOT NULL PRIMARY KEY,
        name text NOT NULL UNIQUE,
        password text NOT NULL,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS company (
        id uuid NOT NULL PRIMARY KEY,
        name text NOT NULL UNIQUE,
        password text NOT NULL,
        email text UNIQUE, --harus ada yng di input dari salah satu field ini 
        phone text UNIQUE, --harus ada yng di input dari salah satu field ini 
        balance int NOT NULL DEFAULT 0,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL,
    );

CREATE TABLE
    IF NOT EXISTS freelancer (
        id uuid NOT NULL PRIMARY KEY,
        name text NOT NULL UNIQUE,
        password text NOT NULL,
        email text UNIQUE, --harus ada yng di input dari salah satu field ini 
        phone text UNIQUE, --harus ada yng di input dari salah satu field ini 
        balance int NOT NULL DEFAULT 0,
        created_at timestamp NOT NULL,
        updated_at timestamp NOT NULL,
    );