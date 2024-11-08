CREATE TABLE IF NOT EXISTS `sequence.sample_data` (
    date DATE,
    project_id STRING,
    num_transactions INT64,
    total_volume_usd FLOAT64
) OPTIONS (
    expiration_timestamp = TIMESTAMP '2024-11-15 00:00:00 UTC',
    description = 'sample data for sequence expire 2024-11-15',
);