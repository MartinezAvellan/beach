package ipma

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"beach/internal/adapters/cache"
)

const baseURL = "https://api.ipma.pt/open-data"

type Client struct {
	http  *http.Client
	cache *cache.MemoryCache
}

func NewClient(httpClient *http.Client, c *cache.MemoryCache) *Client {
	return &Client{http: httpClient, cache: c}
}

// --- Forecast: sea/ocean ---

type SeaForecastResponse struct {
	Owner        string        `json:"owner"`
	ForecastDate string        `json:"forecastDate"`
	DataUpdate   string        `json:"dataUpdate"`
	Data         []SeaForecast `json:"data"`
}

type SeaForecast struct {
	GlobalIDLocal int     `json:"globalIdLocal"`
	Latitude      string  `json:"latitude"`
	Longitude     string  `json:"longitude"`
	WaveHighMin   string  `json:"waveHighMin"`
	WaveHighMax   string  `json:"waveHighMax"`
	WavePeriodMin string  `json:"wavePeriodMin"`
	WavePeriodMax string  `json:"wavePeriodMax"`
	PredWaveDir   string  `json:"predWaveDir"`
	TotalSeaMin   float64 `json:"totalSeaMin"`
	TotalSeaMax   float64 `json:"totalSeaMax"`
	SSTMin        string  `json:"sstMin"`
	SSTMax        string  `json:"sstMax"`
}

func (c *Client) GetSeaForecast() ([]SeaForecast, error) {
	if cached, ok := c.cache.Get("sea_forecast"); ok {
		return cached.([]SeaForecast), nil
	}

	url := baseURL + "/forecast/oceanography/daily/hp-daily-sea-forecast-day0.json"
	var resp SeaForecastResponse
	if err := c.fetchJSON(url, &resp); err != nil {
		return nil, err
	}

	c.cache.Set("sea_forecast", resp.Data)
	return resp.Data, nil
}

// --- Forecast: weather ---

type WeatherForecastResponse struct {
	Owner        string            `json:"owner"`
	ForecastDate string            `json:"forecastDate"`
	DataUpdate   string            `json:"dataUpdate"`
	Data         []WeatherForecast `json:"data"`
}

type WeatherForecast struct {
	GlobalIDLocal  int     `json:"globalIdLocal"`
	Latitude       string  `json:"latitude"`
	Longitude      string  `json:"longitude"`
	TMin           float64 `json:"tMin"`
	TMax           float64 `json:"tMax"`
	PrecipitaProb  string  `json:"precipitaProb"`
	IDWeatherType  int     `json:"idWeatherType"`
	PredWindDir    string  `json:"predWindDir"`
	ClassWindSpeed int     `json:"classWindSpeed"`
}

func (c *Client) GetWeatherForecast() ([]WeatherForecast, error) {
	if cached, ok := c.cache.Get("weather_forecast"); ok {
		return cached.([]WeatherForecast), nil
	}

	url := baseURL + "/forecast/meteorology/cities/daily/hp-daily-forecast-day0.json"
	var resp WeatherForecastResponse
	if err := c.fetchJSON(url, &resp); err != nil {
		return nil, err
	}

	c.cache.Set("weather_forecast", resp.Data)
	return resp.Data, nil
}

// --- UV ---

type UVEntry struct {
	IDPeriodo     int    `json:"idPeriodo"`
	IntervaloHora string `json:"intervaloHora"`
	Data          string `json:"data"`
	GlobalIDLocal int    `json:"globalIdLocal"`
	IUV           string `json:"iUv"`
}

func (c *Client) GetUV() ([]UVEntry, error) {
	if cached, ok := c.cache.Get("uv"); ok {
		return cached.([]UVEntry), nil
	}

	url := baseURL + "/forecast/meteorology/uv/uv.json"
	var entries []UVEntry
	if err := c.fetchJSON(url, &entries); err != nil {
		return nil, err
	}

	c.cache.Set("uv", entries)
	return entries, nil
}

// --- Observations (real-time) ---

// Observations keyed by timestamp -> stationID -> data
type Observations map[string]map[string]*StationObs

type StationObs struct {
	WindSpeedKM   *float64 `json:"intensidadeVentoKM"`
	Temperature   *float64 `json:"temperatura"`
	Radiation     *float64 `json:"radiacao"`
	WindDirection *int     `json:"idDireccVento"`
	Precipitation *float64 `json:"precAcumulada"`
	WindSpeed     *float64 `json:"intensidadeVento"`
	Humidity      *float64 `json:"humidade"`
	Pressure      *float64 `json:"pressao"`
}

func (c *Client) GetObservations() (Observations, error) {
	if cached, ok := c.cache.Get("observations"); ok {
		return cached.(Observations), nil
	}

	url := baseURL + "/observation/meteorology/stations/observations.json"
	var obs Observations
	if err := c.fetchJSON(url, &obs); err != nil {
		return nil, err
	}

	c.cache.Set("observations", obs)
	return obs, nil
}

// --- Weather types ---

type WeatherTypesResponse struct {
	Data []WeatherType `json:"data"`
}

type WeatherType struct {
	DescPT        string `json:"descWeatherTypePT"`
	DescEN        string `json:"descWeatherTypeEN"`
	IDWeatherType int    `json:"idWeatherType"`
}

var weatherTypes map[int]WeatherType

func (c *Client) GetWeatherTypes() (map[int]WeatherType, error) {
	if weatherTypes != nil {
		return weatherTypes, nil
	}

	url := baseURL + "/weather-type-classe.json"
	var resp WeatherTypesResponse
	if err := c.fetchJSON(url, &resp); err != nil {
		return nil, err
	}

	weatherTypes = make(map[int]WeatherType, len(resp.Data))
	for _, wt := range resp.Data {
		weatherTypes[wt.IDWeatherType] = wt
	}
	return weatherTypes, nil
}

func (c *Client) fetchJSON(url string, target any) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned status %d", url, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// GetLatestObservation returns the most recent observation for a given station ID
func (c *Client) GetLatestObservation(stationID string) *StationObs {
	obs, err := c.GetObservations()
	if err != nil {
		return nil
	}

	var latestTime time.Time
	var latest *StationObs

	for ts, stations := range obs {
		stObs, ok := stations[stationID]
		if !ok || stObs == nil {
			continue
		}
		t, err := time.Parse("2006-01-02T15:04", ts)
		if err != nil {
			continue
		}
		if latest == nil || t.After(latestTime) {
			latestTime = t
			latest = stObs
		}
	}

	return latest
}
