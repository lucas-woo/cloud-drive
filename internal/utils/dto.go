package utils

import (
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	assetsmodels "github.com/lucas-woo/cloud-drive/internal/models/assets"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertAssetCursorToResponse(assetCursor *dto.AssetCursor) *mediav1.AssetCursor {
	if assetCursor == nil {
		return nil
	}
	return &mediav1.AssetCursor{
		CreatedAt: timestamppb.New(assetCursor.CreatedAt),
		ObjectId: assetCursor.ObjectId.String(),
	}
}

func ConvertProjectObjectsToResponse(objects []*dto.ProjectObject) ([]*mediav1.ProjectObject) {
	resp := make([]*mediav1.ProjectObject, 0, len(objects))

	for _, obj := range objects {
		resp = append(resp, &mediav1.ProjectObject{
			ProjectId:    obj.ProjectId.String(),
			CollectionId: obj.CollectionId.String(),
			FolderId:     obj.FolderId.String(),
			ObjectId:     obj.ObjectId.String(),
			FileSize:     obj.FileSize,
			Format:       obj.Format,
			IsActive:     obj.IsActive,
			CreatedAt:    timestamppb.New(obj.CreatedAt),
			ModifiedAt:   timestamppb.New(obj.ModifiedAt),
		})
	}

	return resp
}


func ConvertProjectFoldersToResponse(folders []*dto.ProjectFolder) []*mediav1.ProjectFolder {
	if folders == nil {
		return nil
	}

	resp := make([]*mediav1.ProjectFolder, 0, len(folders))

	for _, folder := range folders {
		if folder == nil {
			continue
		}

		var lastUpload, createdAt, modifiedAt *timestamppb.Timestamp

		if folder.LastUpload != nil {
			lastUpload = timestamppb.New(*folder.LastUpload)
		}

		if folder.CreatedAt != nil {
			createdAt = timestamppb.New(*folder.CreatedAt)
		}

		if folder.ModifiedAt != nil {
			modifiedAt = timestamppb.New(*folder.ModifiedAt)
		}

		resp = append(resp, &mediav1.ProjectFolder{
			FolderId:   folder.FolderId.String(),
			FolderName: folder.FolderName,
			FolderSize: folder.FolderSize,
			AssetCount: folder.AssetCount,
			LastUplaod: lastUpload,
			CreatedAt:  createdAt,
			ModifiedAt: modifiedAt,
		})
	}

	return resp
}


func ConvertProjectCollectionsToResponse(collections []*assetsmodels.CollectionModel) []*mediav1.ProjectCollection {
	if collections == nil {
		return nil
	}

	resp := make([]*mediav1.ProjectCollection, 0, len(collections))

	for _, collection := range collections {
		if collection == nil {
			continue
		}

		resp = append(resp, &mediav1.ProjectCollection{
			CollectionId: collection.CollectionId.String(),
			Name:         collection.Name,
			Description:  collection.Description,
			CreatedAt:    timestamppb.New(collection.CreatedAt),
			LastModified: timestamppb.New(collection.LastModified),
			IsPublic:     collection.IsPublic,
		})
	}

	return resp
}