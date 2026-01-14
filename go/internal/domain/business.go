package domain

// BusinessConfig represents the configuration for a business
type BusinessConfig struct {
	ID          string                 `yaml:"id" json:"id"`
	DisplayName string                 `yaml:"displayName" json:"displayName"`
	Timezone    string                 `yaml:"timezone" json:"timezone"`
	Currency    string                 `yaml:"currency" json:"currency"`
	Location    LocationConfig         `yaml:"location" json:"location"`
	Monday      MondayConfig           `yaml:"monday" json:"monday"`
	Stripe      StripeConfig           `yaml:"stripe" json:"stripe"`
	Gmail       GmailConfig            `yaml:"gmail" json:"gmail"`
	Contact     ContactConfig          `yaml:"contact" json:"contact"`
	Slack       SlackConfig            `yaml:"slack" json:"slack"`
	Templates   TemplateConfig         `yaml:"templates" json:"templates"`
	Pipelines   BusinessPipelineConfig `yaml:"pipelines" json:"pipelines"`
}

// MondayConfig holds Monday.com integration settings
type MondayConfig struct {
	APITokenEnv string           `yaml:"apiTokenEnv" json:"apiTokenEnv"`
	Boards      map[string]int64 `yaml:"boards" json:"boards"`
}

// StripeConfig holds Stripe payment settings
type StripeConfig struct {
	APIKeyEnv       string `yaml:"apiKeyEnv" json:"apiKeyEnv"`
	DefaultCurrency string `yaml:"defaultCurrency" json:"defaultCurrency"`
}

// GmailConfig holds Gmail/email settings
type GmailConfig struct {
	Sender     string `yaml:"sender" json:"sender"`
	SenderName string `yaml:"senderName" json:"senderName"`
}

// ContactConfig holds business contact information for email templates
type ContactConfig struct {
	// SupportEmail is the support/contact email address (e.g., "team@example.com")
	// If not set, defaults to extracting domain from Gmail sender or using "support@{businessId}.com"
	SupportEmail string `yaml:"supportEmail" json:"supportEmail"`

	// WebsiteURL is the business website URL (e.g., "https://example.com")
	// If not set, defaults to "https://{businessId}.com" or derived from Gmail sender domain
	WebsiteURL string `yaml:"websiteURL" json:"websiteURL"`

	// LogoURL is the full URL to the business logo image
	// If not set, defaults to "{websiteURL}/wp-content/uploads/logo.jpg" or similar
	LogoURL string `yaml:"logoURL" json:"logoURL"`

	// BookAppointmentURL is the URL for booking appointments
	// If not set, defaults to "{websiteURL}/book-appointment"
	BookAppointmentURL string `yaml:"bookAppointmentURL" json:"bookAppointmentURL"`

	// Phone is the business phone number (optional)
	Phone string `yaml:"phone" json:"phone"`
}

// SlackConfig holds Slack notification settings
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	WebhookEnv string `yaml:"webhookEnv" json:"webhookEnv"`
}

// TemplateConfig holds template paths and settings
type TemplateConfig struct {
	Email map[string]string `yaml:"email" json:"email"`
	// EmailTemplateSettings holds email template configuration
	EmailTemplateSettings EmailTemplateSettings `yaml:"emailTemplateSettings" json:"emailTemplateSettings"`
}

// EmailTemplateSettings holds email template configuration
type EmailTemplateSettings struct {
	// DefaultTemplate is the default template to use for quote emails
	// Options: "original", "apple_style"
	// Default: "original"
	DefaultTemplate string `yaml:"defaultTemplate" json:"defaultTemplate"`
	// AvailableTemplates lists all available templates
	AvailableTemplates []string `yaml:"availableTemplates" json:"availableTemplates"`
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
