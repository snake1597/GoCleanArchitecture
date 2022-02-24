CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    account VARCHAR (100) UNIQUE NOT NULL,
    password VARCHAR (30) NOT NULL,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    birthday DATE NOT NULL,
    refresh_token VARCHAR (30) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);