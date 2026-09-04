package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
	"github.com/lucas-woo/cloud-drive/internal/dto"
)

type S3Utils struct {
}
//todo need to do all transformations
func SelectImageProcessor(ctx context.Context, transformation *mediav1.ImageTransformations) (*exec.Cmd, error) {

	processorPath := os.Getenv("IMAGE_PROCESSOR_PATH")

	if processorPath == "" {
			processorPath = "/app/bin/image-processor"
	}
	

	if transformation.GetScale() != nil {
		// need to add some Validation function for params
		return exec.CommandContext(ctx, processorPath, "scale", fmt.Sprintf("%d", transformation.Scale.GetWidth()), fmt.Sprintf("%d", transformation.Scale.GetHeight())), nil
	}

	if transformation.GetConversion() != nil {
		// validate transformation.GetConversion().GetFormat()
		return exec.CommandContext(ctx, processorPath, "convert", transformation.GetConversion().GetFormat()), nil
	}

	if transformation.GetCompression() != nil {
		return nil, errors.New("no compression yet")
	}

	if transformation.GetCrop() != nil {
		return nil, errors.New("no crop yet")
	}

	return nil, nil
}

func TransformationsToMetadata(transformations *mediav1.ImageTransformations) map[string]string {
	if transformations == nil {
		return nil
	}
	metadata := make(map[string]string, 0)

	if transformations.Scale != nil {

		data := dto.CropTransformation{
			Height: transformations.Scale.GetHeight(),
			Width: transformations.Scale.GetWidth(),
		}
		
		jsonBytes, _ := json.Marshal(data)
		metadata["crop"] = string(jsonBytes)
	}
	if transformations.Crop != nil {

	}
	if transformations.Conversion != nil {

	}
	if transformations.Compression != nil {

	}
	return metadata
}

func NewS3Util() *S3Utils {
	return &S3Utils{}
}