-- +migrate Up
CREATE TABLE IF NOT EXISTS customers (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number  VARCHAR(20)  NOT NULL,
    email         VARCHAR(255),
    full_name     VARCHAR(120) NOT NULL,
    pin_hash      TEXT         NOT NULL,
    kyc_status    VARCHAR(20)  NOT NULL DEFAULT 'UNVERIFIED',
    status        VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_customers_phone_number UNIQUE (phone_number),
    CONSTRAINT chk_customers_phone_e164  CHECK (phone_number ~ '^\+[1-9][0-9]{7,14}$'),
    CONSTRAINT chk_customers_kyc_status  CHECK (kyc_status IN ('UNVERIFIED', 'PENDING', 'VERIFIED', 'REJECTED')),
    CONSTRAINT chk_customers_status      CHECK (status IN ('ACTIVE', 'SUSPENDED', 'CLOSED'))
);

-- Case-insensitive unique email (Rahim@x.com == rahim@x.com)
CREATE UNIQUE INDEX uq_customers_email_lower
    ON customers (LOWER(email))
    WHERE email IS NOT NULL;

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER trg_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();