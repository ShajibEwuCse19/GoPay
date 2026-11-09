-- +migrate Down
DROP TABLE IF EXISTS customers;
DROP INDEX IF EXISTS uq_customers_email_lower;
DROP INDEX IF EXISTS idx_customers_phone_number;