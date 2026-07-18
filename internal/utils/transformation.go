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
	
	if transformation.GetScale() == nil {
		return nil, errors.New("no Scale")
	}

	cmd := exec.CommandContext(ctx, "bin/image-processor", "scale", fmt.Sprintf("%d", transformation.Scale.GetWidth()), fmt.Sprintf("%d", transformation.Scale.GetHeight()))

	return cmd, nil
}