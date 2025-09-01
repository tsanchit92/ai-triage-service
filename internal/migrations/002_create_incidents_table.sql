CREATE TABLE incidents (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    affected_service VARCHAR(255),
    ai_severity VARCHAR(20),
    ai_category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
