package util

import (
	"testing"
	"time"
)

func TestConsolidateWithFallback(t *testing.T) {
	tests := []struct {
		name         string
		primary      string
		fallback     string
		trigger      string
		defaultValue string
		want         string
	}{
		{
			name:         "normal value",
			primary:      "Birthday Party",
			fallback:     "",
			trigger:      "other",
			defaultValue: "Unspecified",
			want:         "Birthday Party",
		},
		{
			name:         "trigger with fallback",
			primary:      "Other",
			fallback:     "Custom Event",
			trigger:      "other",
			defaultValue: "Unspecified",
			want:         "Custom Event",
		},
		{
			name:         "trigger without fallback",
			primary:      "Other",
			fallback:     "",
			trigger:      "other",
			defaultValue: "Unspecified",
			want:         "Unspecified",
		},
		{
			name:         "contains trigger",
			primary:      "Other - Custom",
			fallback:     "My Custom Event",
			trigger:      "other",
			defaultValue: "Unspecified",
			want:         "My Custom Event",
		},
		{
			name:         "empty primary",
			primary:      "",
			fallback:     "",
			trigger:      "other",
			defaultValue: "Unspecified",
			want:         "Unspecified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConsolidateWithFallback(tt.primary, tt.fallback, tt.trigger, tt.defaultValue)
			if got != tt.want {
				t.Errorf("ConsolidateWithFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFirstInteger(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "standard format",
			input: "I Need 2 Helpers",
			want:  2,
		},
		{
			name:  "just number",
			input: "5",
			want:  5,
		},
		{
			name:  "no number",
			input: "I Need Helpers",
			want:  0,
		},
		{
			name:  "multiple numbers",
			input: "I Need 3 Helpers for 2 events",
			want:  3,
		},
		{
			name:  "range format",
			input: "10 - 25 Guests",
			want:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFirstInteger(tt.input)
			if got != tt.want {
				t.Errorf("ExtractFirstInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFirstFloat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "standard format",
			input: "for 4 Hours (minimum)",
			want:  4.0,
		},
		{
			name:  "decimal hours",
			input: "for 5.5 Hours",
			want:  5.5,
		},
		{
			name:  "just number",
			input: "4",
			want:  4.0,
		},
		{
			name:  "no number",
			input: "for Hours",
			want:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFirstFloat(tt.input)
			if got != tt.want {
				t.Errorf("ExtractFirstFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBooleanFromText(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		positiveIndicators  []string
		want                bool
	}{
		{
			name:               "yes default",
			input:              "Yes, I need a call",
			positiveIndicators: nil, // Uses default
			want:               true,
		},
		{
			name:               "no default",
			input:              "No, I don't need a call",
			positiveIndicators: nil,
			want:               false,
		},
		{
			name:               "contains yes",
			input:              "Yes, please call me",
			positiveIndicators: nil,
			want:               true,
		},
		{
			name:               "empty",
			input:              "",
			positiveIndicators: nil,
			want:               false,
		},
		{
			name:               "custom indicators",
			input:              "Confirmed",
			positiveIndicators: []string{"confirmed", "yes", "true"},
			want:               true,
		},
		{
			name:               "custom indicators false",
			input:              "Cancelled",
			positiveIndicators: []string{"confirmed", "yes"},
			want:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseBooleanFromText(tt.input, tt.positiveIndicators...)
			if got != tt.want {
				t.Errorf("ParseBooleanFromText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTimeComponent(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "standard format",
			input:   "July 10, 2025 4:00 PM",
			want:    "4:00 PM",
			wantErr: false,
		},
		{
			name:    "24 hour format",
			input:   "2025-07-10 16:00",
			want:    "4:00 PM",
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTimeComponent(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTimeComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractTimeComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "YYYY-MM-DD format",
			input:   "2025-07-10",
			want:    time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "full date format",
			input:   "July 10, 2025",
			want:    time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformZapierPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload RawZapierPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: RawZapierPayload{
				FirstName:        "John",
				LastName:         "Doe",
				EmailAddress:     "john@example.com",
				EventDate:        "2025-07-10",
				EventTime:        "4:00 PM",
				HelpersRequested: "I Need 2 Helpers",
				ForHowManyHours:  "for 4 Hours",
				Occasion:         "Birthday Party",
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			payload: RawZapierPayload{
				FirstName: "",
				LastName:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			payload: RawZapierPayload{
				FirstName:        "John",
				LastName:         "Doe",
				EmailAddress:     "john@example.com",
				EventDate:        "2025-07-10",
				HelpersRequested: "I Need 2 Helpers",
				ForHowManyHours:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformZapierPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformZapierPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("TransformZapierPayload() returned nil result without error")
			}
		})
	}
}

