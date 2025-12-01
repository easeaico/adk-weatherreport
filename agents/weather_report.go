package agents

import (
	"adk-weatherreport/main/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type getWeatherReportArgs struct {
	City string `json:"city" jsonschema:"The city for which to get the weather report."`
}

type getWeatherReportResult struct {
	Status string `json:"status"`
	Report string `json:"report,omitempty"`
}

// WeatherResponse struct maps the JSON response from OpenWeather API
// Contains city name and main temperature details
type WeatherResponse struct {
	CityName string `json:"name"`
	Weather  []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type ErrorResponse struct {
	Code    string `json:"cod"`
	Message string `json:"message"`
}

func getWeatherReport(ctx tool.Context, args getWeatherReportArgs) (getWeatherReportResult, error) {
	resp, err := getTemperature(args.City)
	if err != nil {
		return getWeatherReportResult{Status: "error", Report: err.Error()}, nil
	}

	report := fmt.Sprintf("The weather in %s is %s with a temperature of %f degrees Celsius.", resp.CityName, resp.Weather[0].Description, resp.Main.Temp)
	return getWeatherReportResult{Status: "success", Report: report}, nil
}

func NewWeatherReportAgent(ctx context.Context) (agent.Agent, error) {
	model, err := models.NewGrokModel(ctx, "grok-4-1-fast", &genai.ClientConfig{
		APIKey: os.Getenv("XAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
		return nil, err
	}

	weatherTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_weather_report",
			Description: "Retrieves the current weather report for a specified city.",
		},
		getWeatherReport,
	)
	if err != nil {
		log.Fatal(err)
	}

	weatherReportAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_report_agent",
		Model:       model,
		Instruction: "You are a helpful assistant that provides weather information and analyzes the sentiment of user feedback. **If the user asks about the weather in a specific city, use the 'get_weather_report' tool to retrieve the weather details.** **If the 'get_weather_report' tool returns a 'success' status, provide the weather report to the user.** **If the 'get_weather_report' tool returns an 'error' status, inform the user that the weather information for the specified city is not available and ask if they have another city in mind.** ",
		Tools:       []tool.Tool{weatherTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	return weatherReportAgent, nil
}

func getTemperature(city string) (*WeatherResponse, error) {
	apiKey := os.Getenv("OWM_API_KEY")
	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&units=metric", apiKey, city)
	res, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if bytes.Index(body, []byte("{\"cod\"")) == 0 {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}

		return nil, fmt.Errorf("failed to get weather report: %s", errorResponse.Message)
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &weatherData, nil
}
