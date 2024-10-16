CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    frequency DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- SELECT * FROM reports