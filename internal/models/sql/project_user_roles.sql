CREATE TABLE project_user_roles (
    user_id BINARY(16) NOT NULL,
    project_id BINARY(16) NOT NULL,
    role VARCHAR(50) NOT NULL,

    PRIMARY KEY (user_id, project_id, role),

    INDEX idx_project (project_id),
    INDEX idx_user (user_id)
);