-- Add city column to reports table
ALTER TABLE reports ADD COLUMN city VARCHAR(100);

-- Create index for city filtering
CREATE INDEX idx_reports_city ON reports(city);

-- Add comment to explain the column
COMMENT ON COLUMN reports.city IS 'Extracted city name for efficient filtering';
