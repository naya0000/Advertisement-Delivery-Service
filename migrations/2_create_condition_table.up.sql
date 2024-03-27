CREATE TABLE IF NOT EXISTS Condition (
    id SERIAL PRIMARY KEY,
    ad_id INT NOT NULL,
    age_start INT,
    age_end INT,
    gender VARCHAR(1)[] DEFAULT NULL, -- ["M", "F"]
    country VARCHAR(255)[] DEFAULT NULL, 
    platform VARCHAR(255)[] DEFAULT NULL,
    FOREIGN KEY (ad_id) REFERENCES Advertisement(id)
);