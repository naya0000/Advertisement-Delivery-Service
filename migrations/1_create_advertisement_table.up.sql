CREATE TABLE IF NOT EXISTS advertisement (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_start_at ON advertisement (start_at);
CREATE INDEX idx_end_at ON advertisement (end_at);