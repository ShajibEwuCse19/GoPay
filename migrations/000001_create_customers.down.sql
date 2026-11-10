-- +migrate Down
DROP TRIGGER IF EXISTS trg_customers_updated_at ON customers;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS customers CASCADE;
