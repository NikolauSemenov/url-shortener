CREATE TABLE IF NOT EXISTS urls(
    original_url VARCHAR(500) NOT NULL,
    alias VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (original_url, alias)
);