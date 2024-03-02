CREATE TABLE IF NOT EXISTS employee (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    surname VARCHAR(255),
    password BYTEA,
    is_admin BOOLEAN,
    phone_number VARCHAR(20),
    enrolled TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    category_id VARCHAR(255),
    price INT,
    description TEXT,
    amount INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
