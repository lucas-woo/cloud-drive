CREATE TABLE project_objects (
    project_id BINARY(16) NOT NULL,
    object_id BINARY(16) NOT NULL,

    folder VARCHAR(255) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_pending BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NULL,
    modified_at TIMESTAMP NULL,

    PRIMARY KEY (project_id, object_id)
);