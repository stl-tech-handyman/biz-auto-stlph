package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/util"
	"gopkg.in/yaml.v3"
)

// SettingsHandler handles settings-related endpoints
type SettingsHandler struct {
	businessLoader *config.BusinessLoader
	logger         *slog.Logger
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(businessLoader *config.BusinessLoader, logger *slog.Logger) *SettingsHandler {
	return &SettingsHandler{
		businessLoader: businessLoader,
		logger:         logger,
	}
}

// EmailTemplateSettingsResponse represents the response for email template settings
type EmailTemplateSettingsResponse struct {
	DefaultTemplate    string   `json:"defaultTemplate"`
	AvailableTemplates []string `json:"availableTemplates"`
	BusinessID         string   `json:"businessId"`
}

// HandleGetEmailTemplateSettings handles GET /api/settings/email-template
func (h *SettingsHandler) HandleGetEmailTemplateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	businessID := r.URL.Query().Get("businessId")
	if businessID == "" {
		businessID = "stlpartyhelpers" // Default
	}

	var settings EmailTemplateSettingsResponse
	settings.BusinessID = businessID
	settings.AvailableTemplates = []string{"original", "apple_style"}
	settings.DefaultTemplate = "original" // Fallback

	if h.businessLoader != nil {
		if businessConfig, err := h.businessLoader.LoadBusiness(r.Context(), businessID); err == nil && businessConfig != nil {
			if businessConfig.Templates.EmailTemplateSettings.DefaultTemplate != "" {
				settings.DefaultTemplate = businessConfig.Templates.EmailTemplateSettings.DefaultTemplate
			}
			if len(businessConfig.Templates.EmailTemplateSettings.AvailableTemplates) > 0 {
				settings.AvailableTemplates = businessConfig.Templates.EmailTemplateSettings.AvailableTemplates
			}
		}
	}

	util.WriteJSON(w, http.StatusOK, settings)
}

// HandleUpdateEmailTemplateSettings handles POST /api/settings/email-template
// This updates the YAML config file (dev only)
func (h *SettingsHandler) HandleUpdateEmailTemplateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		BusinessID      string `json:"businessId"`
		DefaultTemplate string `json:"defaultTemplate"`
	}

	if err := util.ReadJSON(r, &body); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	businessID := body.BusinessID
	if businessID == "" {
		businessID = "stlpartyhelpers"
	}

	// Validate template
	validTemplates := []string{"original", "apple_style"}
	isValid := false
	for _, t := range validTemplates {
		if body.DefaultTemplate == t {
			isValid = true
			break
		}
	}
	if !isValid {
		util.WriteError(w, http.StatusBadRequest, "invalid template. Must be one of: original, apple_style")
		return
	}

	// Load current config
	if h.businessLoader == nil {
		util.WriteError(w, http.StatusInternalServerError, "business loader not available")
		return
	}

	businessConfig, err := h.businessLoader.LoadBusiness(r.Context(), businessID)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to load business config: "+err.Error())
		return
	}

	// Update template settings
	businessConfig.Templates.EmailTemplateSettings.DefaultTemplate = body.DefaultTemplate
	if len(businessConfig.Templates.EmailTemplateSettings.AvailableTemplates) == 0 {
		businessConfig.Templates.EmailTemplateSettings.AvailableTemplates = validTemplates
	}

	// Save to YAML file (dev only - in production this should use a database or config service)
	// Get config path from business loader
	configPath := h.businessLoader.GetBusinessConfigPath(businessID)
	if configPath == "" {
		util.WriteError(w, http.StatusInternalServerError, "config path not available")
		return
	}

	// Read existing YAML to preserve other settings
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to read config file: "+err.Error())
		return
	}

	// Parse YAML
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(configBytes, &configMap); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to parse config file: "+err.Error())
		return
	}

	// Update template settings in map
	if templates, ok := configMap["templates"].(map[string]interface{}); ok {
		if emailTemplateSettings, ok := templates["emailTemplateSettings"].(map[string]interface{}); ok {
			emailTemplateSettings["defaultTemplate"] = body.DefaultTemplate
		} else {
			if templates["emailTemplateSettings"] == nil {
				templates["emailTemplateSettings"] = make(map[string]interface{})
			}
			templates["emailTemplateSettings"].(map[string]interface{})["defaultTemplate"] = body.DefaultTemplate
			templates["emailTemplateSettings"].(map[string]interface{})["availableTemplates"] = validTemplates
		}
	} else {
		configMap["templates"] = map[string]interface{}{
			"emailTemplateSettings": map[string]interface{}{
				"defaultTemplate":    body.DefaultTemplate,
				"availableTemplates": validTemplates,
			},
		}
	}

	// Write back to file
	updatedBytes, err := yaml.Marshal(configMap)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to marshal config: "+err.Error())
		return
	}

	if err := os.WriteFile(configPath, updatedBytes, 0644); err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to write config file: "+err.Error())
		return
	}

	// Invalidate cache so changes take effect immediately
	h.businessLoader.InvalidateCache(businessID)

	// #region agent log
	if f, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{
			"sessionId": "debug-session",
			"runId":     "run1",
			"location":  "settings_handler.go:HandleUpdateEmailTemplateSettings",
			"message":   "Email template settings updated",
			"data": map[string]interface{}{
				"businessId":      businessID,
				"defaultTemplate": body.DefaultTemplate,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.Close()
	}
	// #endregion

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success":         true,
		"defaultTemplate": body.DefaultTemplate,
		"message":         "Settings updated successfully",
	})
}

// HandleGetAllSettings handles GET /api/settings
// Returns all available settings grouped by category
func (h *SettingsHandler) HandleGetAllSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	businessID := r.URL.Query().Get("businessId")
	if businessID == "" {
		businessID = "stlpartyhelpers"
	}

	settings := map[string]interface{}{
		"businessId": businessID,
		"categories": map[string]interface{}{
			"emailTemplates": map[string]interface{}{
				"defaultTemplate":    "original",
				"availableTemplates": []string{"original", "apple_style"},
				"description":        "Email template settings for quote emails",
			},
		},
	}

	// Load from config if available
	if h.businessLoader != nil {
		if businessConfig, err := h.businessLoader.LoadBusiness(r.Context(), businessID); err == nil && businessConfig != nil {
			if businessConfig.Templates.EmailTemplateSettings.DefaultTemplate != "" {
				settings["categories"].(map[string]interface{})["emailTemplates"].(map[string]interface{})["defaultTemplate"] = businessConfig.Templates.EmailTemplateSettings.DefaultTemplate
			}
			if len(businessConfig.Templates.EmailTemplateSettings.AvailableTemplates) > 0 {
				settings["categories"].(map[string]interface{})["emailTemplates"].(map[string]interface{})["availableTemplates"] = businessConfig.Templates.EmailTemplateSettings.AvailableTemplates
			}
		}
	}

	util.WriteJSON(w, http.StatusOK, settings)
}
