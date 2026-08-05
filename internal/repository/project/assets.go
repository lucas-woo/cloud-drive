package projectrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-woo/cloud-drive/internal/config"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	assetsmodels "github.com/lucas-woo/cloud-drive/internal/models/assets"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func (r *ProjectRepository) GetAssets(
	ctx context.Context,
	projectId uuid.UUID,
	cursor *dto.AssetCursor,
	limit int,
) ([]*dto.ProjectObject, error) {

	if limit <= 0 {
		limit = 40
	}

	query := `
		SELECT
			project_id,
			collection_id,
			folder_id,
			object_id,
			file_size,
			format,
			is_active,
			created_at,
			modified_at
		FROM project_objects
		WHERE project_id = ?
	`

	args := []any{projectId[:]}

	if cursor != nil {
		query += `
			AND (
				created_at < ?
				OR (created_at = ? AND object_id < ?)
			)
		`
		args = append(args,
			cursor.CreatedAt,
			cursor.CreatedAt,
			cursor.ObjectId[:],
		)
	}

	query += `
		ORDER BY created_at DESC, object_id DESC
		LIMIT ?
	`

	args = append(args, limit)

	rows, err := r.sqldb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*dto.ProjectObject

	for rows.Next() {
		var asset dto.ProjectObject

		err := rows.Scan(
			&asset.ProjectId,
			&asset.CollectionId,
			&asset.FolderId,
			&asset.ObjectId,
			&asset.FileSize,
			&asset.Format,
			&asset.IsActive,
			&asset.CreatedAt,
			&asset.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *ProjectRepository) GetAllFolders(
	ctx context.Context,
	projectId uuid.UUID,
) ([]*dto.ProjectFolder, error) {

	query := fmt.Sprintf(`
		SELECT
			folder_id,
			folder_name,
			folder_size,
			asset_count,
			last_upload,
			created_at,
			modified_at
		FROM %s
		WHERE project_id = ?
		ORDER BY folder_name ASC
	`, config.ProjectFoldersTable)

	rows, err := r.sqldb.QueryContext(ctx, query, projectId[:])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*dto.ProjectFolder

	for rows.Next() {
		var (
			folderIdBytes []byte
			folder        dto.ProjectFolder
		)

		err := rows.Scan(
			&folderIdBytes,
			&folder.FolderName,
			&folder.FolderSize,
			&folder.AssetCount,
			&folder.LastUpload,
			&folder.CreatedAt,
			&folder.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		folder.FolderId, err = uuid.FromBytes(folderIdBytes)
		if err != nil {
			return nil, err
		}

		folders = append(folders, &folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *ProjectRepository) GetAllCollections(
	ctx context.Context,
	projectId uuid.UUID,
) ([]*assetsmodels.CollectionModel, error) {

	filter := bson.M{
		"project_id": projectId,
	}

	cursor, err := r.mongodb.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	collections := make([]*assetsmodels.CollectionModel, 0)

	for cursor.Next(ctx) {
		var collection assetsmodels.CollectionModel

		if err := cursor.Decode(&collection); err != nil {
			return nil, err
		}

		collections = append(collections, &collection)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

func (r *ProjectRepository) GetAssetsInFolder(
	ctx context.Context,
	projectId uuid.UUID,
	folderId uuid.UUID,
	assetCursor *dto.AssetCursor,
	limit int,
) ([]*dto.ProjectObject, error) {
	var (
		query string
		args  []any
	)

	if assetCursor == nil {
		query = fmt.Sprintf(`
			SELECT
				project_id,
				collection_id,
				folder_id,
				object_id,
				file_size,
				format,
				is_active,
				created_at,
				modified_at
			FROM %s
			WHERE
				project_id = ?
				AND folder_id = ?
			ORDER BY created_at DESC, object_id DESC
			LIMIT ?
		`, config.ProjectObjectsTable)

		args = []any{
			projectId[:],
			folderId[:],
			limit,
		}
	} else {
		query = fmt.Sprintf(`
			SELECT
				project_id,
				collection_id,
				folder_id,
				object_id,
				file_size,
				format,
				is_active,
				created_at,
				modified_at
			FROM %s
			WHERE
				project_id = ?
				AND folder_id = ?
				AND (
					created_at < ?
					OR (
						created_at = ?
						AND object_id < ?
					)
				)
			ORDER BY created_at DESC, object_id DESC
			LIMIT ?
		`, config.ProjectObjectsTable)

		args = []any{
			projectId[:],
			folderId[:],
			assetCursor.CreatedAt,
			assetCursor.CreatedAt,
			assetCursor.ObjectId[:],
			limit,
		}
	}

	rows, err := r.sqldb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*dto.ProjectObject

	for rows.Next() {
		var asset dto.ProjectObject

		if err := rows.Scan(
			&asset.ProjectId,
			&asset.CollectionId,
			&asset.FolderId,
			&asset.ObjectId,
			&asset.FileSize,
			&asset.Format,
			&asset.IsActive,
			&asset.CreatedAt,
			&asset.ModifiedAt,
		); err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *ProjectRepository) GetAssetsInCollection(
	ctx context.Context,
	projectId uuid.UUID,
	collectionId uuid.UUID,
	assetCursor *dto.AssetCursor,
	limit int,
) ([]*dto.ProjectObject, error) {
	var (
		query string
		args  []any
	)

	if assetCursor == nil {
		query = fmt.Sprintf(`
			SELECT
				project_id,
				collection_id,
				folder_id,
				object_id,
				file_size,
				format,
				is_active,
				created_at,
				modified_at
			FROM %s
			WHERE
				project_id = ?
				AND collection_id = ?
			ORDER BY created_at DESC, object_id DESC
			LIMIT ?
		`, config.ProjectObjectsTable)

		args = []any{
			projectId[:],
			collectionId[:],
			limit,
		}
	} else {
		query = fmt.Sprintf(`
			SELECT
				project_id,
				collection_id,
				folder_id,
				object_id,
				file_size,
				format,
				is_active,
				created_at,
				modified_at
			FROM %s
			WHERE
				project_id = ?
				AND collection_id = ?
				AND (
					created_at < ?
					OR (
						created_at = ?
						AND object_id < ?
					)
				)
			ORDER BY created_at DESC, object_id DESC
			LIMIT ?
		`, config.ProjectObjectsTable)

		args = []any{
			projectId[:],
			collectionId[:],
			assetCursor.CreatedAt,
			assetCursor.CreatedAt,
			assetCursor.ObjectId[:],
			limit,
		}
	}

	rows, err := r.sqldb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*dto.ProjectObject

	for rows.Next() {
		var asset dto.ProjectObject

		if err := rows.Scan(
			&asset.ProjectId,
			&asset.CollectionId,
			&asset.FolderId,
			&asset.ObjectId,
			&asset.FileSize,
			&asset.Format,
			&asset.IsActive,
			&asset.CreatedAt,
			&asset.ModifiedAt,
		); err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *ProjectRepository) CreateNewFolder(
	ctx context.Context,
	projectId uuid.UUID,
	folderName string,
) (uuid.UUID, error) {

	folderId := uuid.New()
	now := time.Now().UTC()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			folder_id,
			project_id,
			folder_name,
			folder_size,
			asset_count,
			last_upload,
			created_at,
			modified_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, config.ProjectFoldersTable)

	_, err := r.sqldb.ExecContext(
		ctx,
		query,
		folderId[:],
		projectId[:],
		folderName,
		0,
		0,
		nil,
		now,
		now,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return folderId, nil
}

func (r *ProjectRepository) CreateNewCollection(
	ctx context.Context,
	projectId, creatorId uuid.UUID,
	name, description string,
) (uuid.UUID, error) {

	collectionId := uuid.New()
	now := time.Now().UTC()

	collection := assetsmodels.CollectionModel{
		CollectionId: collectionId,
		ProjectId:    projectId,
		CreatedBy:    creatorId,
		Name:         name,
		Description:  description,
		CreatedAt:    now,
		LastModified: now,
		IsPublic:     false,
	}

	_, err := r.mongodb.InsertOne(ctx, collection)
	if err != nil {
		return uuid.Nil, err
	}

	return collectionId, nil
}