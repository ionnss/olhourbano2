-- Add status to reports
ALTER TABLE reports ADD COLUMN status VARCHAR(255) NOT NULL DEFAULT 'pending';