-- Create the advertisement table
CREATE TABLE IF NOT EXISTS advertisement (
    id SERIAL NOT NULL,
    title VARCHAR(255) NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    conditions JSONB,
    PRIMARY KEY (id, end_at)
) PARTITION BY RANGE (end_at);

-- Create the partitions based on ranges of end_at
CREATE TABLE IF NOT EXISTS advertisement_partition_1 PARTITION OF advertisement FOR VALUES FROM ('2024-03-01T00:00:00Z') TO ('2024-03-31T23:59:59Z');

CREATE TABLE IF NOT EXISTS advertisement_partition_2 PARTITION OF advertisement FOR VALUES FROM ('2024-04-01T00:00:00Z') TO ('2024-04-30T23:59:59Z');


CREATE INDEX ON advertisement (end_at);

-- 可以增加 sub-partitioning for start_at

-- Create indexes on the partitions if needed
-- CREATE INDEX idx_advertisement_partition_1_end_at_start_at ON advertisement_partition_1 (end_at, start_at);
-- CREATE INDEX idx_advertisement_partition_2_end_at_start_at ON advertisement_partition_2 (end_at, start_at);
