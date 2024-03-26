CREATE TABLE IF NOT EXISTS Condition (
    id SERIAL PRIMARY KEY,
    ad_id INT NOT NULL,
    age_start INT,
    age_end INT,
    gender VARCHAR(1), -- M, F
    country VARCHAR(255), 
    platform VARCHAR(255),
    FOREIGN KEY (ad_id) REFERENCES Advertisement(id)
);