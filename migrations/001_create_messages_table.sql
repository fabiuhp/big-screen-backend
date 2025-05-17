CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    content TEXT NOT NULL,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL
); 