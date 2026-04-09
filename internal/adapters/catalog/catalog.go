package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"beach/internal/domain"
)

//go:embed cameras.json
var camerasJSON []byte

type cameraEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Region     string `json:"region"`
	StreamSlug string `json:"stream_slug"`
}

type StaticCatalog struct {
	cameras []domain.Camera
	index   map[string]*domain.Camera
}

func New() (*StaticCatalog, error) {
	var entries []cameraEntry
	if err := json.Unmarshal(camerasJSON, &entries); err != nil {
		return nil, fmt.Errorf("parsing embedded catalog: %w", err)
	}

	now := time.Now().UTC()
	c := &StaticCatalog{
		cameras: make([]domain.Camera, 0, len(entries)),
		index:   make(map[string]*domain.Camera, len(entries)),
	}

	for _, e := range entries {
		cam := domain.Camera{
			ID:         e.ID,
			Name:       e.Name,
			Location:   e.Location,
			Region:     e.Region,
			StreamSlug: e.StreamSlug,
			IsLive:     true,
			UpdatedAt:  now,
		}
		c.cameras = append(c.cameras, cam)
		c.index[cam.ID] = &c.cameras[len(c.cameras)-1]
	}

	return c, nil
}

func (c *StaticCatalog) ListCameras() ([]domain.Camera, error) {
	return c.cameras, nil
}

func (c *StaticCatalog) GetCameraByID(id string) (*domain.Camera, error) {
	cam, ok := c.index[id]
	if !ok {
		return nil, nil
	}
	return cam, nil
}
