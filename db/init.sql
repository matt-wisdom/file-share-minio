CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE files (
    file_id SERIAL PRIMARY,
    object_name VARCHAR(64) NOT NULL,
    file_name VARCHAR(32) NOT NULL,
    owner_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(), 
    FOREIGN KEY owner_id REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE file_shares (
    file_id INT NOT NULL,
    shared_with INT NOT NULL,
    shared_at TIMESTAMP DEFAULT NOW(),
    received_at TIMESTAMP,
    PRIMARY KEY (file_id, shared_with),
    FOREIGN KEY (file_id) REFERENCES files(file_id) ON DELETE CASCADE,
    FOREIGN KEY (shared_with) REFERENCES users(user_id) ON DELETE CASCADE
);