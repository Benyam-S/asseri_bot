CREATE TABLE temp_users (
    telegram_id VARCHAR PRIMARY KEY UNIQUE NOT NULL,
    user_name VARCHAR,
    category VARCHAR,
    phone_number VARCHAR UNIQUE NOT NULL,
    created_at DATETIME,
    updated_at DATETIME
);