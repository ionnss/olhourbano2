-- Migration 003: Add performance indexes for reports table

-- Index for filtering by status (pending, approved, rejected)
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);

-- Index for filtering by problem type
CREATE INDEX IF NOT EXISTS idx_reports_problem_type ON reports(problem_type);

-- Index for ordering by creation date (newest first)
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at DESC);

-- Index for ordering by vote count (most popular first)
CREATE INDEX IF NOT EXISTS idx_reports_vote_count ON reports(vote_count DESC);

-- Composite index for common queries: status + creation date
CREATE INDEX IF NOT EXISTS idx_reports_status_created_at ON reports(status, created_at DESC);

-- Composite index for feed queries: status + vote count + creation date
CREATE INDEX IF NOT EXISTS idx_reports_feed_optimization ON reports(status, vote_count DESC, created_at DESC);

-- Index for geospatial queries (latitude, longitude for location-based searches)
CREATE INDEX IF NOT EXISTS idx_reports_location ON reports(latitude, longitude);

-- Index for user's own reports lookup
CREATE INDEX IF NOT EXISTS idx_reports_hashed_cpf ON reports(hashed_cpf); 