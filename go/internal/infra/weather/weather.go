package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// WeatherService handles weather forecast fetching
type WeatherService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewWeatherService creates a new weather service
func NewWeatherService() (*WeatherService, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		// Try alternative env var name
		apiKey = os.Getenv("WEATHER_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHERMAP_API_KEY or WEATHER_API_KEY environment variable is not set")
	}

	return &WeatherService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.openweathermap.org/data/2.5",
	}, nil
}

// WeatherForecast represents weather forecast data
type WeatherForecast struct {
	Date          time.Time
	Temperature   float64 // in Fahrenheit
	Condition     string  // e.g., "Clear", "Clouds", "Rain"
	Description   string  // e.g., "clear sky", "light rain"
	Humidity      int     // percentage
	WindSpeed     float64 // mph
	Precipitation float64 // mm (0 if no rain)
	Icon          string  // weather icon code
}

// GetForecastForDate gets weather forecast for a specific date and location
// Returns forecast for the event date if within 10 days, nil otherwise
func (w *WeatherService) GetForecastForDate(ctx context.Context, lat, lng float64, eventDate time.Time) (*WeatherForecast, error) {
	// Only fetch weather for events within 10 days
	daysUntilEvent := int(time.Until(eventDate).Hours() / 24)
	if daysUntilEvent < 0 || daysUntilEvent > 10 {
		return nil, nil // Not an error, just outside our forecast window
	}

	// Use forecast API (5-day/3-hour forecast) or daily forecast API
	// For simplicity, we'll use the 5-day/3-hour forecast and find the closest match
	url := fmt.Sprintf("%s/forecast?lat=%.4f&lon=%.4f&appid=%s&units=imperial", w.baseURL, lat, lng, w.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var result struct {
		List []struct {
			DT   int64 `json:"dt"` // Unix timestamp
			Main struct {
				Temp      float64 `json:"temp"`
				Humidity  int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Wind struct {
				Speed float64 `json:"speed"` // m/s, we'll convert to mph
			} `json:"wind"`
			Rain struct {
				ThreeH float64 `json:"3h"` // precipitation in mm for last 3 hours
			} `json:"rain"`
		} `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	if len(result.List) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	// Find the forecast closest to the event date/time
	var closestIndex int = -1
	minDiff := time.Duration(1<<63 - 1) // Max duration

	eventUnix := eventDate.Unix()
	for i := range result.List {
		diff := time.Duration(abs(int64(result.List[i].DT) - eventUnix))
		if diff < minDiff {
			minDiff = diff
			closestIndex = i
		}
	}

	if closestIndex == -1 {
		return nil, fmt.Errorf("could not find forecast for event date")
	}

	closestForecast := result.List[closestIndex]

	// Extract weather data
	weather := closestForecast.Weather[0]
	windSpeedMph := closestForecast.Wind.Speed * 2.237 // Convert m/s to mph

	forecast := &WeatherForecast{
		Date:          eventDate,
		Temperature:   closestForecast.Main.Temp,
		Condition:     weather.Main,
		Description:   weather.Description,
		Humidity:      closestForecast.Main.Humidity,
		WindSpeed:     windSpeedMph,
		Precipitation: closestForecast.Rain.ThreeH,
		Icon:          weather.Icon,
	}

	return forecast, nil
}

// GetWeatherRecommendation returns recommendations based on weather forecast
// Especially useful for outdoor events
func GetWeatherRecommendation(forecast *WeatherForecast, isOutdoor bool) string {
	if forecast == nil {
		return ""
	}

	var recommendations []string

	// Temperature recommendations
	if forecast.Temperature < 50 {
		recommendations = append(recommendations, "Consider providing heaters or warming stations")
	} else if forecast.Temperature > 85 {
		recommendations = append(recommendations, "Consider providing shade, fans, or cooling stations")
	}

	// Precipitation recommendations
	if forecast.Precipitation > 0 {
		if isOutdoor {
			recommendations = append(recommendations, "⚠️ Rain expected — consider backup indoor space or tent coverage")
		} else {
			recommendations = append(recommendations, "Rain expected — ensure easy access to indoor space")
		}
	}

	// Wind recommendations
	if forecast.WindSpeed > 15 {
		if isOutdoor {
			recommendations = append(recommendations, "High winds expected — secure decorations and consider tent stability")
		}
	}

	// Condition-specific recommendations
	switch forecast.Condition {
	case "Rain", "Drizzle":
		if isOutdoor {
			recommendations = append(recommendations, "Wet conditions — plan for covered areas and non-slip surfaces")
		}
	case "Snow":
		if isOutdoor {
			recommendations = append(recommendations, "Snow expected — ensure clear pathways and safe access")
		}
	case "Extreme":
		recommendations = append(recommendations, "⚠️ Extreme weather conditions — consider rescheduling or enhanced safety measures")
	}

	if len(recommendations) == 0 {
		return ""
	}

	result := "Weather considerations: " + recommendations[0]
	for i := 1; i < len(recommendations); i++ {
		result += "; " + recommendations[i]
	}

	return result
}

// abs returns absolute value of int64
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
