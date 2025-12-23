# STL Party Helpers - Go Backend API

This is the Go-based backend implementation for the STL Party Helpers business operations system.

## Project Structure

```
stlph/
├── go/                    # Go backend codebase (main)
│   ├── cmd/              # Application entry points
│   ├── internal/         # Internal packages
│   │   ├── domain/       # Business logic and entities
│   │   ├── ports/        # Interface definitions
│   │   ├── infra/        # Infrastructure implementations
│   │   ├── app/          # Application services
│   │   ├── http/         # HTTP handlers and middleware
│   │   ├── config/       # Configuration loading
│   │   └── util/         # Utilities
│   ├── config/           # Environment configs
│   ├── scripts/          # Deployment and setup scripts
│   └── docs/             # Documentation
├── config/                # Business configurations
│   ├── businesses/       # Business YAML configs
│   └── pipelines/        # Pipeline definitions
├── templates/            # HTML email templates
├── archive/              # Archived JavaScript/Apps Script code
└── oldcode/              # Original codebase reference
```

## Quick Start

See `go/README.md` for detailed setup instructions.

### Local Development

```bash
cd go
go run ./cmd/api
```

The server will start on port 8080 (or PORT environment variable).

## Architecture

The Go backend follows clean architecture principles:

- **Domain**: Core business logic and entities (Business, Pipeline, Job)
- **Ports**: Interfaces for external dependencies (PaymentsProvider, Mailer, CRM, etc.)
- **Infrastructure**: Concrete implementations (Stripe, Gmail, Monday.com, etc.)
- **App**: Application services (FormEventsService, TriggersService)
- **HTTP**: HTTP handlers, middleware, and routing

## API Endpoints

### POST /v1/form-events
Processes form submissions (e.g., from WPForms).

### POST /v1/triggers
Processes trigger-based events (e.g., from Monday.com, Cloud Scheduler).

See `go/docs/api/` for complete API documentation.

## Configuration

Business configurations and pipeline definitions are stored in YAML files:
- `config/businesses/` - Business configurations
- `config/pipelines/` - Pipeline definitions
- `templates/` - HTML email templates

## Deployment

See `go/docs/deployment/` for deployment instructions to Google Cloud Run.

## Archived Code

All previous JavaScript/Google Apps Script code has been moved to the `archive/` directory for reference.

## Documentation

- Main Go API docs: `go/README.md`
- API Reference: `go/docs/api/`
- Deployment Guide: `go/docs/deployment/`
- Development Guide: `go/docs/development/`
