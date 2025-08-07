-- Migration 008: Rollback comments table creation
DROP TABLE IF EXISTS comments;

-- Remove comment count column from reports table
ALTER TABLE reports DROP COLUMN IF EXISTS comment_count;
