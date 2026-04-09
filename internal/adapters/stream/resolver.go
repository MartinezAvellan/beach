package stream

import (
	"fmt"
	"time"

	"beach/internal/domain"
)

const baseStreamURL = "https://video-auth1.iol.pt/beachcam"

type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) ResolveStream(camera *domain.Camera) (*domain.CameraStream, error) {
	if camera.StreamSlug == "" {
		return nil, fmt.Errorf("no stream slug for camera %s", camera.ID)
	}

	return &domain.CameraStream{
		CameraID:  camera.ID,
		StreamURL: fmt.Sprintf("%s/%s/playlist.m3u8", baseStreamURL, camera.StreamSlug),
		Status:    "online",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}, nil
}
