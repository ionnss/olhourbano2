-- Migration 004: Add performance indexes for votes table

-- Index for vote counting by report (if not automatically created by foreign key)
CREATE INDEX IF NOT EXISTS idx_votes_report_id ON votes(report_id);

-- Index for preventing duplicate votes and user vote lookup
CREATE INDEX IF NOT EXISTS idx_votes_hashed_cpf ON votes(vote_hashed_cpf);

-- Unique composite index to prevent duplicate votes from same user on same report
CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_unique_user_report ON votes(vote_hashed_cpf, report_id);

-- Index for ordering votes by creation time
CREATE INDEX IF NOT EXISTS idx_votes_created_at ON votes(created_at DESC);

-- Composite index for user vote history queries
CREATE INDEX IF NOT EXISTS idx_votes_user_time ON votes(vote_hashed_cpf, created_at DESC); 