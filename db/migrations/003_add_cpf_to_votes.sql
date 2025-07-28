-- Add CPF column to votes table and update constraints
ALTER TABLE votes ADD COLUMN cpf_hash VARCHAR(64);

-- Drop the old IP-based unique constraint
ALTER TABLE votes DROP CONSTRAINT IF EXISTS votes_report_id_ip_address_key;

-- Add new CPF-based unique constraint
ALTER TABLE votes ADD CONSTRAINT unique_report_cpf UNIQUE(report_id, cpf_hash);

-- Create index for better performance
CREATE INDEX idx_votes_cpf_hash ON votes(cpf_hash);

-- Add comment explaining the change
COMMENT ON COLUMN votes.cpf_hash IS 'SHA256 hash of CPF for vote verification'; 