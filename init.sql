CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    monthly_income DECIMAL(12,2) NOT NULL,
    monthly_expenses DECIMAL(12,2) NOT NULL,
    marital_status VARCHAR(20) NOT NULL,
    agree_to_be_scored BOOLEAN NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    dependents INTEGER DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    bank_name VARCHAR(100) NOT NULL,
    monthly_payment_amount DECIMAL(12,2),
    total_repayment_amount DECIMAL(12,2),
    number_of_payments INTEGER,
    annual_percentage_rate DECIMAL(12,2),
    first_repayment_date VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bank_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    bank_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    bank_id VARCHAR(100),
    submitted_at TIMESTAMP,
    completed_at TIMESTAMP,
    error TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_offers_application_id ON offers(application_id);
CREATE INDEX IF NOT EXISTS idx_bank_submissions_application_id ON bank_submissions(application_id);
