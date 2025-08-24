CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    merchant_id VARCHAR(50) NOT NULL,
    amount_cents BIGINT NOT NULL,
    fee_cents BIGINT NOT NULL,
    status VARCHAR(255) NOT NULL,
    paid_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_merchant_id ON transactions (merchant_id);
CREATE INDEX idx_transactions_status ON transactions (status);
CREATE INDEX idx_transactions_paid_at ON transactions (paid_at);

CREATE TABLE settlements (
    id BIGSERIAL PRIMARY KEY,
    merchant_id VARCHAR(50) NOT NULL,
    settlement_date DATE NOT NULL,
    gross_cents BIGINT NOT NULL DEFAULT 0,
    fee_cents BIGINT NOT NULL DEFAULT 0,
    net_cents BIGINT NOT NULL DEFAULT 0,
    txn_count BIGINT NOT NULL DEFAULT 0,
    generated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    unique_job_id TEXT NOT NULL,
    UNIQUE (merchant_id, settlement_date)
);

CREATE TABLE jobs (
    job_id TEXT PRIMARY KEY,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    status VARCHAR(255) NOT NULL,
    progress NUMERIC(5, 2) NOT NULL DEFAULT 0,
    processed_count BIGINT NOT NULL DEFAULT 0,
    total_count BIGINT NOT NULL DEFAULT 0,
    result_path TEXT
);

-- Indexes to improve query performance on both tables.
CREATE INDEX idx_settlements_merchant_id ON settlements (merchant_id);
CREATE INDEX idx_jobs_status ON jobs (status);

INSERT INTO transactions (merchant_id, amount_cents, fee_cents, status, paid_at)
SELECT 
    'merchant_' || (RANDOM() * 100)::INTEGER,
    (RANDOM() * 100000 + 1000)::INTEGER,
    (RANDOM() * 1000 + 50)::INTEGER,
    (ARRAY['QUEUED','RUNNING','COMPLETED','FAILED','CANCELLED'])[floor(random() * 5 + 1)],
    TIMESTAMP '2025-01-01' + (RANDOM() * INTERVAL '30 days')
FROM generate_series(1, 1000000);

