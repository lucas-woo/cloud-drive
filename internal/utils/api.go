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

// func (c *RestMapper) ConvertProject