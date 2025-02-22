-- Drop existing tables if they exist
DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS proposals;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS regencies;
DROP TABLE IF EXISTS provinces;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS users;

-- Provinces Table
CREATE TABLE provinces (
    id SERIAL PRIMARY KEY,
    province_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Regencies Table
CREATE TABLE regencies (
    id SERIAL PRIMARY KEY,
    province_id INT NOT NULL,
    regency_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_regencies_province FOREIGN KEY (province_id) REFERENCES provinces (id) ON DELETE CASCADE
);

-- Users Table
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    role VARCHAR(20) CHECK (role IN ('client', 'freelancer', 'company')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admins Table
CREATE TABLE admins (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Jobs Table
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    price DECIMAL(10, 2),
    regency_id INT NOT NULL,
    province_id INT NOT NULL,
    posted_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_jobs_posted_by FOREIGN KEY (posted_by) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_jobs_regency FOREIGN KEY (regency_id) REFERENCES regencies (id) ON DELETE CASCADE,
    CONSTRAINT fk_jobs_province FOREIGN KEY (province_id) REFERENCES provinces (id) ON DELETE CASCADE
);

-- Proposals Table
CREATE TABLE proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    freelancer_id UUID NOT NULL,
    proposal_text TEXT,
    proposed_price DECIMAL(10, 2),
    status VARCHAR(20) CHECK (status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_proposals_job FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
    CONSTRAINT fk_proposals_freelancer FOREIGN KEY (freelancer_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Transactions Table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    freelancer_id UUID NOT NULL,
    company_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'failed')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions_job FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_freelancer FOREIGN KEY (freelancer_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_company FOREIGN KEY (company_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Chat Table
CREATE TABLE chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_chat_sender FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_chat_receiver FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Review Table
CREATE TABLE review (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_review_job FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
    CONSTRAINT fk_review_reviewer FOREIGN KEY (reviewer_id) REFERENCES users (id) ON DELETE CASCADE
);

COPY provinces (id, province_name)
FROM '/docker-entrypoint-initdb.d/data/provinces.csv' DELIMITER ',';

COPY regencies (id, province_id, regency_name)
FROM '/docker-entrypoint-initdb.d/data/regencies.csv' DELIMITER ',';
