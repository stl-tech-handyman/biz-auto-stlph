package domain

// BusinessConfig represents the configuration for a business
type BusinessConfig struct {
	ID          string                `yaml:"id" json:"id"`
	DisplayName string                `yaml:"displayName" json:"displayName"`
	Timezone    string                `yaml:"timezone" json:"timezone"`
	Currency    string                `yaml:"currency" json:"currency"`
	Location    LocationConfig         `yaml:"location" json:"location"`
	Monday      MondayConfig          `yaml:"monday" json:"monday"`
	Stripe      StripeConfig          `yaml:"stripe" json:"stripe"`
	Gmail       GmailConfig           `yaml:"gmail" json:"gmail"`
	Slack       SlackConfig           `yaml:"slack" json:"slack"`
	Templates   TemplateConfig        `yaml:"templates" json:"templates"`
	Pipelines   BusinessPipelineConfig `yaml:"pipelines" json:"pipelines"`
}

// MondayConfig holds Monday.com integration settings
type MondayConfig struct {
	APITokenEnv string            `yaml:"apiTokenEnv" json:"apiTokenEnv"`
	Boards      map[string]int64   `yaml:"boards" json:"boards"`
}

// StripeConfig holds Stripe payment settings
type StripeConfig struct {
	APIKeyEnv      string `yaml:"apiKeyEnv" json:"apiKeyEnv"`
	DefaultCurrency string `yaml:"defaultCurrency" json:"defaultCurrency"`
}

// GmailConfig holds Gmail/email settings
type GmailConfig struct {
	Sender     string `yaml:"sender" json:"sender"`
	SenderName string `yaml:"senderName" json:"senderName"`
}

// SlackConfig holds Slack notification settings
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	WebhookEnv string `yaml:"webhookEnv" json:"webhookEnv"`
}

// TemplateConfig holds template paths
type TemplateConfig struct {
	Email map[string]string `yaml:"email" json:"email"`
}

// BusinessPipelineConfig holds pipeline configuration
type BusinessPipelineConfig struct {
	DefaultForm string            `yaml:"defaultForm" json:"defaultForm"`
	Triggers    map[string]string `yaml:"triggers" json:"triggers"`
}

// LocationConfig holds business location settings for distance calculations
// This is separate from the business address to allow flexibility in choosing
// the origin point for distance calculations (e.g., warehouse vs office)
type LocationConfig struct {
	// OfficeAddress is the business office address (for display purposes)
	// Example: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
	OfficeAddress string `yaml:"officeAddress" json:"officeAddress"`
	
	// DistanceOrigin is the address or coordinates used as the starting point
	// for all distance calculations (travel fees, service area checks, etc.)
	// This can be different from OfficeAddress if you want to calculate from
	// a warehouse, depot, or other location.
	// 
	// Can be:
	// - An address string: "4220 Duncan Ave Suite 201, St. Louis, MO 63110"
	// - Coordinates: "38.6255,-90.2456" (lat,lng format)
	//
	// If not specified, defaults to OfficeAddress
	// NOTE: If Lat/Lng fields are provided, they take precedence over distanceOrigin
	DistanceOrigin string `yaml:"distanceOrigin" json:"distanceOrigin"`
	
	// Lat and Lng are the latitude and longitude of the distance origin
	// If provided, these take precedence over distanceOrigin string
	// These coordinates are used for all distance calculations (travel fees, service area checks)
	// Recommended for accuracy - avoids geocoding API calls
	Lat float64 `yaml:"lat" json:"lat"`
	Lng float64 `yaml:"lng" json:"lng"`
	
	// ServiceRadiusMiles is the service area radius in miles
	// Locations within this radius have no travel fee
	// Default: 15.0 miles
	ServiceRadiusMiles float64 `yaml:"serviceRadiusMiles" json:"serviceRadiusMiles"`
}
