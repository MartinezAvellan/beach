package ipma

import (
	"fmt"
	"math"
	"time"

	"beach/internal/domain"
)

// SeaLocation maps to IPMA sea forecast locations
type SeaLocation struct {
	GlobalIDLocal int
	Name          string
	Lat, Lon      float64
}

// WeatherLocation maps to IPMA district forecast locations
type WeatherLocation struct {
	GlobalIDLocal int
	Name          string
	Lat, Lon      float64
}

// Station maps to IPMA real-time observation stations
type Station struct {
	ID       string
	Name     string
	Lat, Lon float64
}

// Known IPMA sea forecast locations (from sea-locations.json)
var seaLocations = []SeaLocation{
	{1080517, "Faro", 37.0017, -8.0},
	{1081525, "Sagres", 37.0, -8.9383},
	{1110600, "Lisboa", 38.65, -9.31},
	{1060300, "Figueira da Foz", 40.1417, -8.8783},
	{1131200, "Porto", 41.175, -8.76},
	{1160900, "Viana do Castelo", 41.67, -8.8333},
	{1151300, "Sines", 37.95, -8.8833},
	{2310300, "Funchal", 32.64, -16.91},
	{2320126, "Porto Santo", 33.04, -16.34},
	{3420300, "Ponta Delgada", 37.68, -25.67},
	{3470100, "Horta", 38.57, -28.47},
	{3480200, "Santa Cruz das Flores", 39.45, -31.13},
}

// Known IPMA weather forecast districts (from distrits-islands.json)
var weatherLocations = []WeatherLocation{
	{1010500, "Aveiro", 40.6413, -8.6535},
	{1020500, "Beja", 38.02, -7.87},
	{1030300, "Braga", 41.5475, -8.4227},
	{1040200, "Bragança", 41.8076, -7.7606},
	{1050200, "Castelo Branco", 39.8217, -7.4957},
	{1060300, "Coimbra", 40.2081, -8.4194},
	{1070500, "Évora", 38.5701, -7.9104},
	{1080500, "Faro", 37.0146, -7.9331},
	{1081505, "Sagres", 37.0168, -8.9403},
	{1081100, "Portimão", 37.15, -8.52},
	{1090700, "Guarda", 40.5379, -7.2647},
	{1100900, "Leiria", 39.7473, -8.8069},
	{1110600, "Lisboa", 38.766, -9.1286},
	{1121400, "Portalegre", 39.29, -7.42},
	{1131200, "Porto", 41.158, -8.6294},
	{1141600, "Santarém", 39.2, -8.74},
	{1151200, "Setúbal", 38.5246, -8.8856},
	{1151300, "Sines", 37.956, -8.8643},
	{1160900, "Viana do Castelo", 41.6952, -8.8365},
	{1171400, "Vila Real", 41.3053, -7.744},
	{1182300, "Viseu", 40.6585, -7.912},
	{2310300, "Funchal", 32.6485, -16.9084},
	{2320100, "Porto Santo", 33.07, -16.34},
	{3420300, "Ponta Delgada", 37.7415, -25.6677},
	{3470100, "Horta", 38.5363, -28.6315},
}

// Known coastal observation stations (subset of stations.json)
var stations = []Station{
	{"1210881", "Olhão", 37.033, -7.821},
	{"1210883", "Tavira", 37.122, -7.621},
	{"1240610", "Viana do Castelo", 41.695, -8.829},
	{"6212126", "Esposende", 41.526, -8.780},
	{"1200535", "Lisboa/Gago Coutinho", 38.77, -9.13},
	{"1210762", "Peniche/Cabo Carvoeiro", 39.36, -9.41},
	{"1200579", "Cascais", 38.69, -9.42},
	{"1210702", "Leiria", 39.83, -8.88},
	{"1210761", "Nazaré", 39.60, -9.07},
	{"1210755", "Figueira da Foz", 40.15, -8.86},
	{"1131200", "Porto/Pedras Rubras", 41.24, -8.68},
	{"1080500", "Faro", 37.02, -7.97},
	{"1081505", "Sagres", 37.01, -8.95},
	{"1151300", "Sines", 37.95, -8.88},
	{"1151200", "Setúbal", 38.52, -8.90},
}

// CameraLocation stores approximate coordinates for camera regions
type CameraLocation struct {
	Lat, Lon float64
}

// Known camera locations (approximate coordinates for Portuguese beaches)
var cameraLocations = map[string]CameraLocation{
	// Peniche area
	"peniche": {39.36, -9.38},
	// Nazaré area
	"nazare": {39.60, -9.07},
	// Costa da Caparica
	"costadacaparica": {38.63, -9.24},
	"costa":           {38.63, -9.24},
	// Cascais / Estoril / Carcavelos
	"cascais": {38.69, -9.42}, "estoril": {38.70, -9.40},
	"carcavelos": {38.68, -9.34}, "parede": {38.69, -9.36},
	"santoamaro": {38.69, -9.35}, "torre": {38.68, -9.32},
	// Ericeira
	"ericeira": {38.96, -9.42},
	// Algarve
	"faro": {37.02, -7.93}, "lagos": {37.10, -8.67},
	"portimao": {37.12, -8.54}, "albufeira": {37.09, -8.25},
	"sagres": {37.01, -8.95}, "tavira": {37.13, -7.65},
	// Porto area
	"porto": {41.15, -8.63}, "matosinhos": {41.18, -8.69},
	"espinho": {41.01, -8.64}, "gaia": {41.12, -8.65},
	// North
	"viana": {41.69, -8.85}, "esposende": {41.53, -8.79},
	// Central
	"figueiradafoz": {40.15, -8.86}, "aveiro": {40.64, -8.65},
	// Setúbal
	"setubal": {38.52, -8.90}, "sines": {37.95, -8.88},
	// Islands
	"funchal": {32.65, -16.91}, "portosanto": {33.07, -16.34},
	"pontadelgada": {37.74, -25.67},
	// Lisbon area
	"lisboa": {38.72, -9.14}, "almada": {38.68, -9.16},
}

type ConditionsResolver struct {
	client *Client
}

func NewConditionsResolver(client *Client) *ConditionsResolver {
	return &ConditionsResolver{client: client}
}

func (r *ConditionsResolver) ResolveConditions(camera *domain.Camera) (*domain.CameraConditions, error) {
	loc := guessCameraLocation(camera)

	cond := &domain.CameraConditions{
		CameraID:  camera.ID,
		FetchedAt: time.Now().UTC(),
	}

	// Sea forecast (wave height, period, direction, SST)
	seaData, err := r.client.GetSeaForecast()
	if err == nil {
		nearest := findNearestSea(loc, seaData)
		if nearest != nil {
			cond.WaveHeight = nearest.WaveHighMax + " m"
			cond.WavePeriod = nearest.WavePeriodMax + " s"
			cond.WaveDirection = nearest.PredWaveDir
			cond.WaterTemp = nearest.SSTMax + " ºC"
		}
	}

	// Weather forecast (temp, wind, weather type)
	weatherData, err := r.client.GetWeatherForecast()
	if err == nil {
		nearest := findNearestWeather(loc, weatherData)
		if nearest != nil {
			cond.AirTemp = fmt.Sprintf("%.0f ºC", nearest.TMax)
			cond.WindDirection = nearest.PredWindDir
			cond.WindSpeed = windClassToSpeed(nearest.ClassWindSpeed)

			types, _ := r.client.GetWeatherTypes()
			if wt, ok := types[nearest.IDWeatherType]; ok {
				cond.Weather = wt.DescPT
			}
		}
	}

	// UV
	uvData, err := r.client.GetUV()
	if err == nil {
		nearest := findNearestUV(loc, uvData)
		if nearest != nil {
			cond.UVIndex = nearest.IUV
		}
	}

	// Real-time observation (actual wind speed, temperature)
	nearestStation := findNearestStation(loc)
	if nearestStation != nil {
		obs := r.client.GetLatestObservation(nearestStation.ID)
		if obs != nil {
			if obs.WindSpeedKM != nil && *obs.WindSpeedKM >= 0 {
				cond.WindSpeed = fmt.Sprintf("%.1f km/h", *obs.WindSpeedKM)
			}
			if obs.Temperature != nil && *obs.Temperature > -90 {
				cond.AirTemp = fmt.Sprintf("%.1f ºC", *obs.Temperature)
			}
			if obs.WindDirection != nil {
				cond.WindDirection = windDirCodeToText(*obs.WindDirection)
			}
			if obs.Humidity != nil && *obs.Humidity >= 0 {
				cond.Humidity = fmt.Sprintf("%.0f%%", *obs.Humidity)
			}
		}
	}

	return cond, nil
}

func guessCameraLocation(camera *domain.Camera) CameraLocation {
	id := camera.ID

	// Try direct match
	if loc, ok := cameraLocations[id]; ok {
		return loc
	}

	// Try prefix match
	for prefix, loc := range cameraLocations {
		if len(prefix) >= 4 && len(id) >= len(prefix) && id[:len(prefix)] == prefix {
			return loc
		}
	}

	// Try substring match
	for keyword, loc := range cameraLocations {
		if len(keyword) >= 4 && contains(id, keyword) {
			return loc
		}
	}

	// Default: Lisboa
	return CameraLocation{38.72, -9.14}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && searchSubstring(s, substr))
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func findNearestSea(loc CameraLocation, data []SeaForecast) *SeaForecast {
	var nearest *SeaForecast
	minDist := math.MaxFloat64
	for i := range data {
		lat := parseFloat(data[i].Latitude)
		lon := parseFloat(data[i].Longitude)
		d := haversine(loc.Lat, loc.Lon, lat, lon)
		if d < minDist {
			minDist = d
			nearest = &data[i]
		}
	}
	return nearest
}

func findNearestWeather(loc CameraLocation, data []WeatherForecast) *WeatherForecast {
	var nearest *WeatherForecast
	minDist := math.MaxFloat64
	for i := range data {
		lat := parseFloat(data[i].Latitude)
		lon := parseFloat(data[i].Longitude)
		d := haversine(loc.Lat, loc.Lon, lat, lon)
		if d < minDist {
			minDist = d
			nearest = &data[i]
		}
	}
	return nearest
}

func findNearestUV(loc CameraLocation, data []UVEntry) *UVEntry {
	// UV entries use globalIdLocal, match to weather locations
	today := time.Now().Format("2006-01-02")
	var best *UVEntry
	minDist := math.MaxFloat64

	for i := range data {
		if data[i].Data != today {
			continue
		}
		// Find coordinates for this globalIdLocal
		for _, wl := range weatherLocations {
			if wl.GlobalIDLocal == data[i].GlobalIDLocal {
				d := haversine(loc.Lat, loc.Lon, wl.Lat, wl.Lon)
				if d < minDist {
					minDist = d
					best = &data[i]
				}
				break
			}
		}
	}
	return best
}

func findNearestStation(loc CameraLocation) *Station {
	var nearest *Station
	minDist := math.MaxFloat64
	for i := range stations {
		d := haversine(loc.Lat, loc.Lon, stations[i].Lat, stations[i].Lon)
		if d < minDist {
			minDist = d
			nearest = &stations[i]
		}
	}
	return nearest
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371 // km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return r * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func windClassToSpeed(class int) string {
	switch class {
	case 1:
		return "Fraco"
	case 2:
		return "Moderado"
	case 3:
		return "Forte"
	case 4:
		return "Muito forte"
	default:
		return ""
	}
}

func windDirCodeToText(code int) string {
	dirs := map[int]string{
		0: "", 1: "Norte", 2: "Nordeste", 3: "Este", 4: "Sudeste",
		5: "Sul", 6: "Sudoeste", 7: "Oeste", 8: "Noroeste", 9: "Norte",
	}
	if d, ok := dirs[code]; ok {
		return d
	}
	return ""
}
