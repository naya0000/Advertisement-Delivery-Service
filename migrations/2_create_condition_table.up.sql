CREATE TABLE IF NOT EXISTS condition (
    id SERIAL PRIMARY KEY,
    ad_id INT NOT NULL,
    age_start INT DEFAULT 1,
    age_end INT DEFAULT 100,
    gender VARCHAR(1) DEFAULT NULL, -- "M" or "F"
    country VARCHAR(255)[] DEFAULT NULL, 
    platform VARCHAR(255)[] DEFAULT NULL,
    FOREIGN KEY (ad_id) REFERENCES advertisement(id)
);