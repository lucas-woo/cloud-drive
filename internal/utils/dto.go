package utils

import (
	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)


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