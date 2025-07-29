CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    problem_type VARCHAR(100),
    hashed_cpf VARCHAR(64),
    email VARCHAR(100),
    location TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    description TEXT,
    photo_path TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    vote_count INTEGER DEFAULT 0,
    status VARCHAR(100) NOT NULL DEFAULT 'pending'
); 