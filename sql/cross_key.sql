-- create new table for parse id from new table to old table

CREATE TABLE id_mapping (
    source_table VARCHAR(100),
    source_column VARCHAR(100),
    source_id_value VARCHAR,
    target_table VARCHAR(100),
    target_column VARCHAR(100),
    target_id_value VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (source_table, source_id_value)
)

CREATE INDEX idx_id_mapping_source_id ON id_mapping(source_id_value);

CREATE INDEX idx_id_mapping_source_table_and_id ON id_mapping(source_table, source_id_value);


