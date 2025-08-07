-- Migration 008: Create comments table for reports
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    hashed_cpf VARCHAR(64) NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 500),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_comments_report_id ON comments(report_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_hashed_cpf ON comments(hashed_cpf);

-- Add comment count column to reports table
ALTER TABLE reports ADD COLUMN comment_count INTEGER DEFAULT 0;

-- Create index for comment count
CREATE INDEX IF NOT EXISTS idx_reports_comment_count ON reports(comment_count DESC);
