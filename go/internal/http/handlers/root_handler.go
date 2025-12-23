package handlers

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

//go:embed endpoint_manager.html
var endpointManagerHTML embed.FS

// EndpointInfo represents an API endpoint
type EndpointInfo struct {
	Method       string
	Path         string
	Summary      string
	AuthRequired bool
	Tags         []string
}

// EndpointManagerData is the data structure for the endpoint manager template
type EndpointManagerData struct {
	Environment        string
	Version            string
	TotalEndpoints     int
	PublicEndpoints    int
	ProtectedEndpoints int
	Endpoints          []EndpointInfo
}

// RootHandler handles the root endpoint
type RootHandler struct {
	environment string
}

// NewRootHandler creates a new root handler
func NewRootHandler(environment string) *RootHandler {
	return &RootHandler{
		environment: environment,
	}
}

// getEndpoints returns all API endpoints information
func getEndpoints() []EndpointInfo {
	return []EndpointInfo{
		// Health endpoints
		{Method: "GET", Path: "/", Summary: "Корневой эндпоинт - Endpoint Manager", AuthRequired: false, Tags: []string{"Здоровье и информация"}},
		{Method: "GET", Path: "/api/health", Summary: "Проверка здоровья сервиса", AuthRequired: false, Tags: []string{"Здоровье и информация"}},
		{Method: "GET", Path: "/api/health/ready", Summary: "Проверка готовности для load balancers", AuthRequired: false, Tags: []string{"Здоровье и информация"}},
		{Method: "GET", Path: "/api/health/live", Summary: "Проверка жизнеспособности для container orchestration", AuthRequired: false, Tags: []string{"Здоровье и информация"}},

		// Stripe endpoints
		{Method: "POST", Path: "/api/stripe/deposit", Summary: "Создание депозита - генерация инвойса", AuthRequired: true, Tags: []string{"Stripe"}},
		{Method: "GET", Path: "/api/stripe/deposit/calculate", Summary: "Расчет рекомендуемого депозита", AuthRequired: true, Tags: []string{"Stripe"}},
		{Method: "POST", Path: "/api/stripe/deposit/with-email", Summary: "Создание депозита с отправкой email (ORCHESTRATED)", AuthRequired: true, Tags: []string{"Stripe"}},
		{Method: "POST", Path: "/api/stripe/final-invoice", Summary: "Создание финального инвойса", AuthRequired: true, Tags: []string{"Stripe"}},
		{Method: "POST", Path: "/api/stripe/final-invoice/with-email", Summary: "Создание финального инвойса с отправкой email (ORCHESTRATED)", AuthRequired: true, Tags: []string{"Stripe"}},
		{Method: "POST", Path: "/api/stripe/test", Summary: "Тест интеграции со Stripe", AuthRequired: true, Tags: []string{"Stripe"}},

		// Estimate endpoints
		{Method: "POST", Path: "/api/estimate", Summary: "Расчет стоимости мероприятия", AuthRequired: true, Tags: []string{"Расчет стоимости"}},
		{Method: "GET", Path: "/api/estimate/special-dates", Summary: "Получение списка специальных дат", AuthRequired: true, Tags: []string{"Расчет стоимости"}},

		// Email endpoints
		{Method: "POST", Path: "/api/email/test", Summary: "Тест отправки email", AuthRequired: true, Tags: []string{"Email"}},
		{Method: "POST", Path: "/api/email/booking-deposit", Summary: "Отправка email с информацией о депозите", AuthRequired: true, Tags: []string{"Email"}},

		// Calendar endpoints
		{Method: "POST", Path: "/api/calendar/create", Summary: "Создание события в Google Calendar", AuthRequired: false, Tags: []string{"Календарь"}},

		// Lead processing endpoints
		{Method: "POST", Path: "/api/business/{businessId}/process-lead", Summary: "Обработка лида для конкретного бизнеса", AuthRequired: false, Tags: []string{"Обработка лидов"}},
		{Method: "POST", Path: "/api/zapier/process-lead", Summary: "Обработка лида от Zapier (legacy)", AuthRequired: false, Tags: []string{"Обработка лидов"}},

		// V1 Pipeline endpoints
		{Method: "POST", Path: "/v1/form-events", Summary: "Обработка событий форм", AuthRequired: false, Tags: []string{"V1 Pipeline"}},
		{Method: "POST", Path: "/v1/triggers", Summary: "Обработка триггеров", AuthRequired: false, Tags: []string{"V1 Pipeline"}},
	}
}

// HandleRoot handles GET /
func (h *RootHandler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	// #region agent log
	if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "root_handler.go:88", "message": "HandleRoot called (ROOT HANDLER MATCHING?)", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion

	// Don't handle Swagger routes - they should be handled by SwaggerHandler
	if r.URL.Path == "/swagger" || r.URL.Path == "/swagger-ui" || r.URL.Path == "/swagger.html" ||
		r.URL.Path == "/swagger-simple" || r.URL.Path == "/api/openapi.json" || r.URL.Path == "/api/openapi.yaml" {
		// #region agent log
		if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "root_handler.go:96", "message": "Root handler rejecting Swagger path", "data": map[string]interface{}{"path": r.URL.Path, "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		http.NotFound(w, r)
		return
	}

	// Don't handle API routes - they should be handled by their specific handlers
	// If we're here, it means ServeMux matched "/" instead of a specific route
	// This shouldn't happen, but we'll reject API paths explicitly
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/v1/") {
		// #region agent log
		if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "H", "location": "root_handler.go:107", "message": "Root handler rejecting API path (should not match)", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method, "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	env := h.environment
	if env == "" {
		env = os.Getenv("ENV")
		if env == "" {
			env = "dev"
		}
	}

	// Check if client wants JSON (for API clients)
	if r.Header.Get("Accept") == "application/json" || r.URL.Query().Get("format") == "json" {
		response := map[string]interface{}{
			"service":     "bizops360-api-go",
			"version":     "1.0.0",
			"status":      "running",
			"environment": env,
			"endpoints": map[string]interface{}{
				"health": []string{
					"GET /api/health",
					"GET /api/health/ready",
					"GET /api/health/live",
				},
				"stripe": []string{
					"POST /api/stripe/deposit",
					"GET /api/stripe/deposit/calculate",
					"POST /api/stripe/deposit/with-email",
					"POST /api/stripe/test",
				},
				"estimate": []string{
					"POST /api/estimate",
					"GET /api/estimate/special-dates",
				},
				"email": []string{
					"POST /api/email/test",
					"POST /api/email/booking-deposit",
				},
				"v1": []string{
					"POST /v1/form-events",
					"POST /v1/triggers",
				},
			},
			"documentation": "https://github.com/bizops360/bizops360",
		}
		util.WriteJSON(w, http.StatusOK, response)
		return
	}

	// Serve HTML endpoint manager
	endpoints := getEndpoints()
	publicCount := 0
	protectedCount := 0
	for _, ep := range endpoints {
		if ep.AuthRequired {
			protectedCount++
		} else {
			publicCount++
		}
	}

	data := EndpointManagerData{
		Environment:        env,
		Version:            "1.0.0",
		TotalEndpoints:     len(endpoints),
		PublicEndpoints:    publicCount,
		ProtectedEndpoints: protectedCount,
		Endpoints:          endpoints,
	}

	// Load and parse template
	tmplContent, err := endpointManagerHTML.ReadFile("endpoint_manager.html")
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to load template: "+err.Error())
		return
	}

	tmpl, err := template.New("endpoint_manager").Parse(string(tmplContent))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to parse template: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to execute template: "+err.Error())
		return
	}
}
