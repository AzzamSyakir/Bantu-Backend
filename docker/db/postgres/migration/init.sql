DROP TABLE IF EXISTS review;
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS escrow;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS proposals;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS regencies;
DROP TABLE IF EXISTS provinces;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS users;

CREATE TABLE provinces (
  id SERIAL NOT NULL PRIMARY KEY,
  province_name VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE regencies (
  id SERIAL NOT NULL PRIMARY KEY,
  province_id INT NOT NULL,
  regency_name VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_regencies_province FOREIGN KEY (province_id) REFERENCES provinces (id) ON DELETE CASCADE
);

CREATE TABLE users (
  id uuid NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  balance  BIGINT NOT NULL DEFAULT 0,
  role VARCHAR(20) CHECK (role IN ('client', 'freelancer', 'company')) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE admins (
  id uuid NOT NULL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE jobs (
  id uuid NOT NULL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),
  price DECIMAL(10, 2),
  regency_id INT NOT NULL,
  province_id INT NOT NULL,
  posted_by uuid NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_jobs_posted_by FOREIGN KEY (posted_by) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_jobs_regency FOREIGN KEY (regency_id) REFERENCES regencies (id) ON DELETE CASCADE,
  CONSTRAINT fk_jobs_province FOREIGN KEY (province_id) REFERENCES provinces (id) ON DELETE CASCADE
);

CREATE TABLE proposals (
  id uuid PRIMARY KEY,
  job_id uuid NOT NULL,
  freelancer_id uuid NOT NULL,
  proposal_text TEXT,
  proposed_price DECIMAL(10, 2),
  status VARCHAR(20) CHECK (status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_proposals_job FOREIGN KEY (job_id) REFERENCES jobs (id) ON DELETE CASCADE,
  CONSTRAINT fk_proposals_freelancer FOREIGN KEY (freelancer_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE transactions (
  id uuid PRIMARY KEY,
  job_id uuid NOT NULL,
  proposal_id uuid NOT NULL,
  sender_id uuid NOT NULL,
  receiver_id uuid NOT NULL,
  amount   BIGINT NOT NULL DEFAULT 0,
  transaction_type VARCHAR(20) CHECK (transaction_type IN ('top_up', 'pay_freelancer', 'withdrawal')) NOT NULL,
  payment_method VARCHAR(20) CHECK (payment_method IN ('virtual_account', 'e_money', 'debit', 'credit', 'pay_later', 'qr', 'payment_link', 'internal_wallet')),
  status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'failed')) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE escrow (
  id uuid PRIMARY KEY,
  transaction_id uuid NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
  amount   BIGINT NOT NULL DEFAULT 0,
  status VARCHAR(20) CHECK (status IN ('pending', 'released', 'refunded')) DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat (
  id uuid PRIMARY KEY,
  sender_id uuid NOT NULL,
  receiver_id uuid NOT NULL,
  message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  read_at TIMESTAMP NULL DEFAULT NULL,
  CONSTRAINT fk_chat_sender FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_chat_receiver FOREIGN KEY (receiver_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE review (
  id uuid PRIMARY KEY,
  job_id uuid NOT NULL,
  reviewer_id uuid NOT NULL,
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
