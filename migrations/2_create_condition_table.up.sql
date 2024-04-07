-- -- Create the advertisement_condition table
-- CREATE TABLE IF NOT EXISTS ad_condition (
--     id SERIAL PRIMARY KEY,
--     ad_id INT NOT NULL,
--     attribute VARCHAR(255) NOT NULL,
--     value VARCHAR(255) NOT NULL
--     -- FOREIGN KEY (ad_id) REFERENCES advertisement(id)
-- ) PARTITION BY RANGE (ad_id);

-- -- Create the partitions based on ranges of end_at
-- CREATE TABLE IF NOT EXISTS advertisement_partition_1 PARTITION OF advertisement FOR VALUES FROM ('2024-03-01T00:00:00Z') TO ('2024-03-31T23:59:59Z');

-- CREATE TABLE IF NOT EXISTS advertisement_partition_2 PARTITION OF advertisement FOR VALUES FROM ('2024-04-01T00:00:00Z') TO ('2024-04-30T23:59:59Z');

-- CREATE INDEX ON advertisement (end_at);