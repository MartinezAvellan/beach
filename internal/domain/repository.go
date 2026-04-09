package domain

type CameraRepository interface {
	ListCameras() ([]Camera, error)
	GetCameraByID(id string) (*Camera, error)
}

type StreamResolver interface {
	ResolveStream(camera *Camera) (*CameraStream, error)
}

type ConditionsResolver interface {
	ResolveConditions(camera *Camera) (*CameraConditions, error)
}
