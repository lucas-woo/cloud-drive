package utils

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	mediav1 "github.com/lucas-woo/cloud-drive/api/media/v1"
)

//todo need to do all transformations
func SelectImageProcessor(ctx context.Context, transformation *mediav1.ImageTransformations) (*exec.Cmd, error) {
	
	if transformation.GetScale() != nil {
		// need to add some Validation function for params
		return exec.CommandContext(ctx, "bin/image-processor", "scale", fmt.Sprintf("%d", transformation.Scale.GetWidth()), fmt.Sprintf("%d", transformation.Scale.GetHeight())), nil
	}

	if transformation.GetConversion() != nil {
		// validate transformation.GetConversion().GetFormat()
		return exec.CommandContext(ctx, "bin/image-processor", "convert", transformation.GetConversion().GetFormat()), nil
	}

	if transformation.GetCompression() != nil {
		return nil, errors.New("no compression yet")
	}

	if transformation.GetCrop() != nil {
		return nil, errors.New("no crop yet")
	}

	return nil, nil
}