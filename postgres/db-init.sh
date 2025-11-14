#!/bin/bash
set -e

# Используем переменные по умолчанию, если не заданы через env
POSTGRES_USER=${POSTGRES_USER:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
PGPASSWORD=$POSTGRES_PASSWORD
export PGPASSWORD

echo "Создание баз данных и пользователя program..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER program WITH PASSWORD 'program';
    CREATE DATABASE tickets;
    GRANT ALL PRIVILEGES ON DATABASE tickets TO program;
    CREATE DATABASE flights;
    GRANT ALL PRIVILEGES ON DATABASE flights TO program;
    CREATE DATABASE privileges;
    GRANT ALL PRIVILEGES ON DATABASE privileges TO program;
EOSQL

echo "Создание таблиц и вставка данных в tickets..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "tickets" <<-EOSQL
    CREATE TABLE IF NOT EXISTS ticket (
        id SERIAL PRIMARY KEY,
        ticket_uid UUID UNIQUE NOT NULL,
        username VARCHAR(80) NOT NULL,
        flight_number VARCHAR(20) NOT NULL,
        price INT NOT NULL,
        status VARCHAR(20) NOT NULL CHECK (status IN ('PAID', 'CANCELED'))
    );
EOSQL

echo "Создание таблиц и вставка данных в flights..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "flights" <<-EOSQL
    CREATE TABLE IF NOT EXISTS airport (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255),
        city VARCHAR(255),
        country VARCHAR(255)
    );

    CREATE TABLE IF NOT EXISTS flight (
        id SERIAL PRIMARY KEY,
        flight_number VARCHAR(20) NOT NULL,
        datetime TIMESTAMP WITH TIME ZONE NOT NULL,
        from_airport_id INT REFERENCES airport(id),
        to_airport_id INT REFERENCES airport(id),
        price INT NOT NULL
    );

    INSERT INTO airport (id, name, city, country) VALUES
    (1, 'Шереметьево', 'Москва', 'Россия'),
    (2, 'Пулково', 'Санкт-Петербург', 'Россия')
    ON CONFLICT DO NOTHING;

    INSERT INTO flight (id, flight_number, datetime, from_airport_id, to_airport_id, price) VALUES
    (1, 'AFL031', '2021-10-08 20:00', 2, 1, 1500)
    ON CONFLICT DO NOTHING;
EOSQL

echo "Создание таблиц и вставка данных в privileges..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "privileges" <<-EOSQL
    CREATE TABLE IF NOT EXISTS privilege (
        id SERIAL PRIMARY KEY,
        username VARCHAR(80) NOT NULL UNIQUE,
        status VARCHAR(80) NOT NULL DEFAULT 'BRONZE' CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
        balance INT
    );

    CREATE TABLE IF NOT EXISTS privilege_history (
        id SERIAL PRIMARY KEY,
        privilege_id INT REFERENCES privilege(id),
        ticket_uid UUID NOT NULL,
        datetime TIMESTAMP NOT NULL,
        balance_diff INT NOT NULL,
        operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('FILL_IN_BALANCE', 'DEBIT_THE_ACCOUNT'))
    );

    INSERT INTO privilege (username, status, balance)
    VALUES ('Test Max', 'GOLD', 1500)
    ON CONFLICT DO NOTHING;

    INSERT INTO privilege_history (privilege_id, ticket_uid, datetime, balance_diff, operation_type)
    VALUES (1, '049161bb-badd-4fa8-9d90-87c9a82b0668', '2021-10-08T19:59:19Z', 1500, 'FILL_IN_BALANCE')
    ON CONFLICT DO NOTHING;
EOSQL

echo "Инициализация баз данных завершена!"
