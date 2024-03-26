CREATE TABLE IF NOT EXISTS Condition (
    id SERIAL PRIMARY KEY,
    advertisement_id INT NOT NULL,
    age_start INT,
    age_end INT,
    gender VARCHAR(1),
    country VARCHAR(255),
    platform VARCHAR(255),
    FOREIGN KEY (advertisement_id) REFERENCES Advertisement(id)
);