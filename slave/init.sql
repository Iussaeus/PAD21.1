-- Slave database initialization script

-- Create the replication role to match the master's configuration
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replicatorpassword';

-- Grant the replicator role appropriate permissions
GRANT CONNECT ON DATABASE mydatabase TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;

-- Ensure future tables grant permissions automatically
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;

