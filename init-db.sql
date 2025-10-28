-- Create databases for each service
CREATE DATABASE service1_db;
CREATE DATABASE service2_db;
CREATE DATABASE service3_db;
CREATE DATABASE saga_db;

-- Grant privileges (optional, as postgres user already has all privileges)
GRANT ALL PRIVILEGES ON DATABASE service1_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE service2_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE service3_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE saga_db TO postgres;
