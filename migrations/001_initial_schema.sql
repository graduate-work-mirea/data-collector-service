-- Create raw_data table
CREATE TABLE raw_data (
    id SERIAL PRIMARY KEY,
    data_source_id INT DEFAULT 1, -- For MVP fixed source
    data JSONB NOT NULL,
    collected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
