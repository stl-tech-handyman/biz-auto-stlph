package domain

// BusinessConfig represents the configuration for a business
type BusinessConfig struct {
	ID          string                `yaml:"id" json:"id"`
	DisplayName string                `yaml:"displayName" json:"displayName"`
	Timezone    string                `yaml:"timezone" json:"timezone"`
	Currency    string                `yaml:"currency" json:"currency"`
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

