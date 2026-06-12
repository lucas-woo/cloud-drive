package media

import "github.com/lucas-woo/cloud-drive/internal/database"

type Service struct {
	mediaResources *database.MediaResources
}

func (s *Service) CreateNewProject() {
	
}

func NewMediaService(mediaResources *database.MediaResources) *Service {
	return &Service{
		mediaResources: mediaResources,
	}
}