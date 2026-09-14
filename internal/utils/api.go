package utils

import (
	"path/filepath"
	"strings"

	iamv1 "github.com/lucas-woo/cloud-drive/api/iam/v1"
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/api"
)

type RestMapper struct {}



func (c *RestMapper) ConvertAssetCursor(cursor *mediav1.AssetCursor) *api.AssetCursor {
	if cursor == nil {
		return nil
	}
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

		if ts := folder.GetLastUplaod(); ts != nil {
			t := ts.AsTime()
			item.LastUpload = &t
		}

		if ts := folder.GetCreatedAt(); ts != nil {
			t := ts.AsTime()
			item.CreatedAt = &t
		}

		if ts := folder.GetModifiedAt(); ts != nil {
			t := ts.AsTime()
			item.ModifiedAt = &t
		}

		result = append(result, item)
	}

	return result
}

func (m *RestMapper) FindFormat(filename string) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "JPEG"
	case ".png":
		return "PNG"
	case ".gif":
		return "GIF"
	case ".webp":
		return "WEBP"
	case ".svg":
		return "SVG"
	case ".bmp":
		return "BMP"
	case ".tiff", ".tif":
		return "TIFF"
	case ".ico":
		return "ICO"
	case ".avif":
		return "AVIF"

	case ".pdf":
		return "PDF"
	case ".doc", ".docx":
		return "WORD"
	case ".xls", ".xlsx":
		return "EXCEL"
	case ".ppt", ".pptx":
		return "POWERPOINT"
	case ".txt":
		return "TEXT"
	case ".csv":
		return "CSV"
	case ".json":
		return "JSON"
	case ".xml":
		return "XML"

	case ".zip":
		return "ZIP"
	case ".rar":
		return "RAR"
	case ".7z":
		return "7Z"
	case ".tar":
		return "TAR"
	case ".gz", ".gzip":
		return "GZIP"

	case ".mp4":
		return "MP4"
	case ".mov":
		return "MOV"
	case ".avi":
		return "AVI"
	case ".mkv":
		return "MKV"
	case ".webm":
		return "WEBM"

	case ".mp3":
		return "MP3"
	case ".wav":
		return "WAV"
	case ".flac":
		return "FLAC"
	case ".aac":
		return "AAC"
	case ".ogg":
		return "OGG"

	case ".html", ".htm":
		return "HTML"
	case ".css":
		return "CSS"
	case ".js":
		return "JAVASCRIPT"
	case ".ts":
		return "TYPESCRIPT"
	case ".go":
		return "GO"
	case ".py":
		return "PYTHON"
	case ".java":
		return "JAVA"
	case ".c":
		return "C"
	case ".cpp", ".cc", ".cxx":
		return "C++"
	case ".rs":
		return "RUST"

	default:
		return strings.TrimPrefix(strings.ToUpper(filepath.Ext(filename)), ".")
	}
}


func (m *RestMapper) ConvertApiKeys(apiKeys []*iamv1.ApiKey) []*api.ApiKey {
	result := make([]*api.ApiKey, 0, len(apiKeys))

	for _, key := range apiKeys {
		if key == nil {
			continue
		}

		result = append(result, &api.ApiKey{
			ApiKey:    key.ApiKey,
			Name:      key.KeyName,
			CreatedAt: key.CreatedAt.AsTime(),
			IsActive: key.IsActive,
		})
	}

	return result
}

func (m *RestMapper) ToProtoTransformations(
	t *api.ImageTransformations,
) *mediav1.ImageTransformations {
	return &mediav1.ImageTransformations{
		Crop: &mediav1.Crop{
			Width:  t.Crop.Width,
			Height: t.Crop.Height,
		},
		Scale: &mediav1.Scale{
			Width:  t.Scale.Width,
			Height: t.Scale.Height,
		},
		Compression: &mediav1.Compression{
			Compress: t.Compression.Compress,
		},
		Conversion: &mediav1.Conversion{
			Format: t.Conversion.Format,
		},
	}
}

func (m *RestMapper) ConvertAllFoldersToApiResponse(folders []*mediav1.ProjectFolder) *api.GetAllFoldersApiResponse {
	response := &api.GetAllFoldersApiResponse{
		Folders: make([]api.FolderApi, 0, len(folders)),
	}

	for _, folder := range folders {
		if folder == nil {
			continue
		}

		response.Folders = append(response.Folders, api.FolderApi{
			FolderId:   folder.FolderId,
			FolderName: folder.FolderName,
		})
	}

	return response
}