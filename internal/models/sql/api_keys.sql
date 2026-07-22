CREATE TABLE api_keys (
    api_key BINARY(16) NOT NULL,

    project_id BINARY(16) NOT NULL,

    name VARCHAR(255) NOT NULL,

    api_secret_hash BINARY(32) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL,

    PRIMARY KEY (api_key),

    INDEX idx_project (project_id)
);


CREATE TABLE api_key_permissions (
    api_key BINARY(16) NOT NULL,
    permission VARCHAR(50) NOT NULL,

    PRIMARY KEY (api_key, permission)
);