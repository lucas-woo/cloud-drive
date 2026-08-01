CREATE TABLE project_folders (
    folder_id BINARY(16) NOT NULL,
    project_id BINARY(16) NOT NULL,

    folder_name VARCHAR(255) NOT NULL,

    folder_size BIGINT UNSIGNED NOT NULL DEFAULT 0,
    asset_count INT UNSIGNED NOT NULL DEFAULT 0,

    last_upload TIMESTAMP NULL,
    created_at TIMESTAMP NULL,
    modified_at TIMESTAMP NULL,

    PRIMARY KEY (folder_id),

    INDEX idx_project_id (project_id)
);