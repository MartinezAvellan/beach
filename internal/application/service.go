package application

import (
	"errors"

	"beach/internal/domain"
)

var (
	ErrCameraNotFound        = errors.New("camera_not_found")
	ErrStreamUnavailable     = errors.New("stream_unavailable")
	ErrConditionsUnavailable = errors.New("conditions_unavailable")
)

type CameraService struct {
	repo       domain.CameraRepository
	resolver   domain.StreamResolver
	conditions domain.ConditionsResolver
}

func NewCameraService(repo domain.CameraRepository, resolver domain.StreamResolver, conditions domain.ConditionsResolver) *CameraService {
	return &CameraService{
		repo:       repo,
		resolver:   resolver,
		conditions: conditions,
	}
}

func (s *CameraService) ListCameras() ([]domain.Camera, error) {
	return s.repo.ListCameras()
}

func (s *CameraService) GetCameraByID(id string) (*domain.Camera, error) {
	camera, err := s.repo.GetCameraByID(id)
	if err != nil {
		return nil, err
	}
	if camera == nil {
		return nil, ErrCameraNotFound
	}
	return camera, nil
}

func (s *CameraService) ResolveCameraStream(id string) (*domain.CameraStream, error) {
	camera, err := s.GetCameraByID(id)
	if err != nil {
		return nil, err
	}

	stream, err := s.resolver.ResolveStream(camera)
	if err != nil {
		return nil, ErrStreamUnavailable
	}
	return stream, nil
}

func (s *CameraService) ResolveCameraConditions(id string) (*domain.CameraConditions, error) {
	camera, err := s.GetCameraByID(id)
	if err != nil {
		return nil, err
	}

	cond, err := s.conditions.ResolveConditions(camera)
	if err != nil {
		return nil, ErrConditionsUnavailable
	}
	return cond, nil
}
