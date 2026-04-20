CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_login text,
    user_password text
);