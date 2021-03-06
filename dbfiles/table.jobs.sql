CREATE TABLE jobs (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    employer VARCHAR(255),
    intiator_id VARCHAR(255),
    title VARCHAR(255),
    description TEXT,
    type VARCHAR(255),
    sector VARCHAR(255),
    education_level VARCHAR(255),
    experience VARCHAR(255),
    gender VARCHAR(255),
    status VARCHAR(255),
    contact_types VARCHAR(255),
    contact_info VARCHAR(255), 
    link TEXT,
    due_date DATETIME,
    post_type VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME
);