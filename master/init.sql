-- Master database initialization script

-- Create a table for key-value storage
CREATE TABLE IF NOT EXISTS data (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- Insert sample data
INSERT INTO data (key, value) VALUES
('master_key1', 'master_value1'),
('master_key2', 'master_value2');

-- Create a replication role
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replicatorpassword';

-- Grant read, write, and modify access to the replicator
GRANT CONNECT ON DATABASE mydatabase TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO replicator;

-- Ensure future tables grant permissions automatically
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO replicator;

