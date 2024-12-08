CREATE DATABASE IF NOT EXISTS learned_vocabulary;

USE learned_vocabulary;

-- Create vocabulary table
CREATE TABLE IF NOT EXISTS vocabulary (
    word VARCHAR(255),
    frequency INT,
    error_count INT,
    PRIMARY KEY (word)
);

-- Create schedule table
CREATE TABLE IF NOT EXISTS schedule (
    id BIGINT,
    time DATE NOT NULL,
    word VARCHAR(255) UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS vocabulary_statistics (
    id INT AUTO_INCREMENT,
    time DATE,
    word_learned INT,
    wrong INT,
    PRIMARY KEY (id)
)
