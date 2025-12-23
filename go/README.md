# Go Backend API

This is the Go-based backend implementation for the BizOps360 Platform, designed to run alongside the existing JavaScript/Node.js backend.

## Architecture

The Go backend follows clean architecture principles with clear separation of concerns:

- **Domain**: Core business logic and entities (Business, Pipeline, Job)
- **Ports**: Interfaces for external dependencies (PaymentsProvider, Mailer, CRM, etc.)
- **Infrastructure**: Concrete implementations (Stripe, Gmail, Monday.com, etc.)
- **App**: Application services (FormEventsService, TriggersService)
- **HTTP**: HTTP handlers, middleware, and routing

## Project Structure

```
go/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain models and business logic
│   │   ├── business.go
│   │   ├── pipeline.go
│   │   ├── job.go
│   │   ├── pipeline_runner.go
│   │   └── errors.go
│   ├── ports/                   # Interface definitions
│   │   ├── payments.go
│   │   ├── mailer.go
│   │   ├── crm.go
│   │   └── ...
│   ├── infra/                   # Infrastructure implementations
│   │   ├── db/
│   │   ├── stripe/
│   │   ├── gmail/
│   │   └── ...
│   ├── app/                     # Application services
│   │   ├── form_events_service.go
│   │   ├── triggers_service.go
│   │   └── actions.go
│   ├── http/                    # HTTP layer
│   │   ├── router.go
│   │   ├── handlers/
│   │   └── middleware/
│   ├── config/                  # Configuration loading
│   │   ├── config.go
│   │   └── business_loader.go
│   └── util/                    # Utilities
│       ├── httpjson.go
│       ├── hmac.go
│       └── tracing.go
└── Dockerfile                   # Container build file
```

## Configuration

Business configurations and pipeline definitions are stored in YAML files at the project root:

- `config/businesses/` - Business configurations
- `config/pipelines/` - Pipeline definitions
- `templates/` - HTML email templates

## API Endpoints

### POST /v1/form-events

Processes form submissions (e.g., from WPForms).

**Headers:**
- `X-Business-Id` (optional, can be in body)
- `X-Pipeline-Key` (optional, can be in body)
- `X-Source` (optional)
- `X-Dry-Run` (optional)
- `X-Signature` (optional, for HMAC auth)

**Body:**
```json
{
  "businessId": "stlpartyhelpers",
  "pipelineKey": "quote_and_deposit",
  "dryRun": false,
  "options": {
    "sendQuoteEmail": true
  },
  "fields": {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "event_date": "2025-05-10"
  }
}
```

### POST /v1/triggers

Processes trigger-based events (e.g., from Monday.com, Cloud Scheduler).

**Headers:**
- `X-Business-Id` (optional, can be in body)
- `X-Trigger-Key` (optional, can be in body)
- `X-Pipeline-Key` (optional, can be in body)
- `X-Source` (optional)
- `X-Dry-Run` (optional)
- `X-Signature` (optional, for HMAC auth)

**Body:**
```json
{
  "source": "monday",
  "businessId": "stlpartyhelpers",
  "triggerKey": "send_renewal_offer",
  "resource": {
    "type": "monday_item",
    "boardId": 123456789,
    "itemId": 987654321
  },
  "payload": {
    "event_date": "2024-05-10",
    "client_email": "client@example.com"
  }
}
```

## Building and Running

### Local Development

```bash
cd go
go run ./cmd/api
```

The server will start on port 8080 (or PORT environment variable).

### Docker Build

```bash
docker build -t bizops-api -f go/Dockerfile .
```

### Environment Variables

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (dev/prod, default: dev)
- `CONFIG_DIR` - Path to config directory (default: /app/config)
- `TEMPLATES_DIR` - Path to templates directory (default: /app/templates)
- `LOG_LEVEL` - Log level (debug/info/warn/error, default: info)
- `HMAC_SECRET` - Secret for HMAC signature verification (optional)

## Deployment to Google Cloud Run

See the main project documentation for Cloud Run deployment instructions. The Dockerfile is configured to:

1. Build the Go binary
2. Copy config and templates into the image
3. Run the server on the PORT environment variable (Cloud Run sets this)

## Current Status

This is a skeleton implementation with:

✅ Domain models and pipeline runner
✅ HTTP server with routing and middleware
✅ Config/template loading
✅ Stub action implementations
✅ In-memory job storage

🚧 TODO (future work):
- Implement real Stripe integration
- Implement real Gmail integration
- Implement real Monday.com integration
- Implement real Slack integration
- Add Firestore/Cloud SQL for job persistence
- Implement more pipeline actions
- Add comprehensive error handling
- Add request validation
- Add tests

## Notes

- The Go backend is the **primary** backend implementation
- Config/template files are shared with the Go backend
- All API endpoints are implemented in Go

