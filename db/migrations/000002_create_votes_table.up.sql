CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    vote_hashed_cpf VARCHAR(64)
); 