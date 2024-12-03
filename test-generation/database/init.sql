CREATE DATABASE IF NOT EXISTS learned_vocabulary;

USE learned_vocabulary;

-- Create vocabulary table
CREATE TABLE IF NOT EXISTS vocabulary (
    Word VARCHAR(255),
    frequency INT,
    error_count INT,
    PRIMARY KEY (Word)
);

-- Create schedule table
CREATE TABLE IF NOT EXISTS schedule (
    id BIGINT,
    time DateTime,
    word VARCHAR(255),
    PRIMARY KEY (id)
);

SELECT * FROM vocabulary WHERE frequency > 0