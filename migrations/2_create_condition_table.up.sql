CREATE TABLE IF NOT EXISTS Condition (
    id SERIAL PRIMARY KEY,
    ad_id INT NOT NULL,
    age_start INT, -- 1
    age_end INT,   -- 100
    gender VARCHAR(1) DEFAULT NULL, -- "M" or "F"
    country VARCHAR(255)[] DEFAULT NULL, -- ["TW", "JP"]
    platform VARCHAR(255)[] DEFAULT NULL, -- ["android", "ios"]
    FOREIGN KEY (ad_id) REFERENCES Advertisement(id)
);