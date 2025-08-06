-- Add transport-specific fields to reports table
ALTER TABLE reports ADD COLUMN transport_type VARCHAR(50);
ALTER TABLE reports ADD COLUMN transport_data JSONB;

-- Create index on transport_type for faster queries
CREATE INDEX idx_reports_transport_type ON reports(transport_type);

-- Create GIN index on transport_data for JSON queries
CREATE INDEX idx_reports_transport_data ON reports USING GIN(transport_data); 