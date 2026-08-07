package utils

import (
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
)

type RestMapper struct {}



func (c *RestMapper) ConvertAssetCursor(cursor *mediav1.AssetCursor) *api.AssetCursor {
	return &api.AssetCursor{
		ObjectId: cursor.GetObjectId(),
		CreatedAt: cursor.CreatedAt.AsTime(),
	}	
}

func (c *RestMapper) ConvertProjectObjectSlice(objects []*mediav1.ProjectObject) []api.ProjectObject {
    if objects == nil {
        return nil
    }

    result := make([]api.ProjectObject, 0, len(objects))

    for _, obj := range objects {
        if obj == nil {
            continue
        }

        result = append(result, api.ProjectObject{
            ProjectId:    obj.GetProjectId(),
            CollectionId: obj.GetCollectionId(),
            FolderId:     obj.GetFolderId(),
            ObjectId:     obj.GetObjectId(),
            FileSize:     obj.GetFileSize(),
            Format:       obj.GetFormat(),
            IsActive:     obj.GetIsActive(),
            CreatedAt:    obj.GetCreatedAt().AsTime(),
            ModifiedAt:   obj.GetModifiedAt().AsTime(),
        })
    }

    return result
}

func (c *RestMapper) ConvertCollectionSlice(collections []*mediav1.ProjectCollection) []api.Collection {
	if collections == nil {
		return nil
	}

	result := make([]api.Collection, 0, len(collections))

	for _, collection := range collections {
		if collection == nil {
			continue
		}

		item := api.Collection{
			CollectionId: collection.GetCollectionId(),
			Name:         collection.GetName(),
			Description:  collection.GetDescription(),
			IsPublic:     collection.GetIsPublic(),
		}

		if collection.GetCreatedAt() != nil {
			item.CreatedAt = collection.GetCreatedAt().AsTime()
		}

		if collection.GetLastModified() != nil {
			item.LastModified = collection.GetLastModified().AsTime()
		}

		result = append(result, item)
	}

	return result
}


func (c *RestMapper) ConvertFolderSlice(folders []*mediav1.ProjectFolder) []api.Folder {
	if folders == nil {
		return nil
	}

	result := make([]api.Folder, 0, len(folders))

	for _, folder := range folders {
		if folder == nil {
			continue
		}

		item := api.Folder{
			FolderId:   folder.GetFolderId(),
			FolderName: folder.GetFolderName(),
			FolderSize: folder.GetFolderSize(),
			AssetCount: folder.GetAssetCount(),
		}

		if folder.GetLastUplaod() != nil {
			item.LastUpload = folder.GetLastUplaod().AsTime()
		}

		if folder.GetCreatedAt() != nil {
			item.CreatedAt = folder.GetCreatedAt().AsTime()
		}

		if folder.GetModifiedAt() != nil {
			item.ModifiedAt = folder.GetModifiedAt().AsTime()
		}

		result = append(result, item)
	}

	return result
}