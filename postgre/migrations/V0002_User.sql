CREATE TABLE users(
    user_id INT PRIMARY KEY,
    first_name  VARCHAR(128),
    last_name VARCHAR(128),
    user_name VARCHAR(128),
    language_code VARCHAR(16),
    clan VARCHAR(128),
    nickname VARCHAR(128),
    has_sudo bool
);

