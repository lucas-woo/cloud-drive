CREATE TABLE api_keys (
    id BINARY(16) NOT NULL,

    project_id BINARY(16) NOT NULL,

    name VARCHAR(255) NOT NULL,

    api_key VARCHAR(255) NOT NULL UNIQUE,

    api_secret VARCHAR(255) NOT NULL,

    created_at TIMESTAMP NOT NULL

    PRIMARY KEY (id),

    INDEX idx_project (project_id)
);


CREATE TABLE api_key_permissions (
    api_key_id BINARY(16) NOT NULL,
    permission VARCHAR(50) NOT NULL,

    PRIMARY KEY (api_key_id, permission)

    INDEX idx_api_key (api_key_id)
);