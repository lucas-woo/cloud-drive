CREATE TABLE project_objects (
    project_id BINARY(16) NOT NULL,
    object_id VARCHAR(255) NOT NULL,

    folder VARCHAR(255) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,

    PRIMARY KEY (project_id, object_id)
);