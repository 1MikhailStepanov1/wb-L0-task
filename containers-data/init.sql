CREATE DATABASE order_db;
CREATE USER order_service_user WITH PASSWORD 'Passw0rd';
GRANT ALL PRIVILEGES ON DATABASE "order_db" to order_service_user;