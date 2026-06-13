CREATE TABLE project_admin (
    
    user_id BINARY(16) NOT NULL,
    project_id BINARY(16) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, project_id)
);