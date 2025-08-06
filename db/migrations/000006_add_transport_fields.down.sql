-- Remove transport-specific fields from reports table
DROP INDEX IF EXISTS idx_reports_transport_data;
DROP INDEX IF EXISTS idx_reports_transport_type;
ALTER TABLE reports DROP COLUMN IF EXISTS transport_data;
ALTER TABLE reports DROP COLUMN IF EXISTS transport_type; 