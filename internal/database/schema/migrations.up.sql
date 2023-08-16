CREATE TABLE favorite_terminals (
    user_id INTEGER,
    terminal_id INTEGER[]
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE terminals(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(255) NOT NULL
);
