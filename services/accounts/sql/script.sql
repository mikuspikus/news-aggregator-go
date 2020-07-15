CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    uid UUID PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    created_at timestamp with time zone not null,
    edited_at timestamp with time zone not null,
    password_hash CHAR(60) NOT NULL,
    is_admin boolean default false
);

