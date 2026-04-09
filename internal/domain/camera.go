package domain

import "time"

type Camera struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	Region        string    `json:"region"`
	StreamSlug    string    `json:"stream_slug,omitempty"`
	IsLive        bool      `json:"is_live"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CameraStream struct {
	CameraID  string    `json:"camera_id"`
	StreamURL string    `json:"stream_url"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CameraConditions struct {
	CameraID      string    `json:"camera_id"`
	WaveHeight    string    `json:"wave_height"`
	WavePeriod    string    `json:"wave_period"`
	WaveDirection string    `json:"wave_direction"`
	WindSpeed     string    `json:"wind_speed"`
	WindDirection string    `json:"wind_direction"`
	WaterTemp     string    `json:"water_temp"`
	AirTemp       string    `json:"air_temp"`
	UVIndex       string    `json:"uv_index"`
	TideTime      string    `json:"tide_time"`
	TideHeight    string    `json:"tide_height"`
	Weather       string    `json:"weather"`
	Humidity      string    `json:"humidity,omitempty"`
	FetchedAt     time.Time `json:"fetched_at"`
}
