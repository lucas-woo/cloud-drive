CREATE TABLE project_objects (
    project_id BINARY(16) NOT NULL,
    collection_id BINARY(16) NULL,
    folder_id BINARY(16) NOT NULL,
    object_id BINARY(16) NOT NULL,

    file_size BIGINT UNSIGNED NULL,
    format VARCHAR(255) NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_pending BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NULL,
    modified_at TIMESTAMP NULL,

    PRIMARY KEY (object_id),

    INDEX idx_project_id (project_id),
    INDEX idx_collection_id (collection_id),
    INDEX idx_folder_id (folder_id),

    INDEX idx_project_created_object (
        project_id,
        created_at DESC,
        object_id DESC
    ),

    INDEX idx_project_folder_created_object (
        project_id,
        folder_id,
        created_at DESC,
        object_id DESC
    ),

    INDEX idx_project_collection_created_object (
        project_id,
        collection_id,
        created_at DESC,
        object_id DESC
    )
);