CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(report_id, ip_address)
); 