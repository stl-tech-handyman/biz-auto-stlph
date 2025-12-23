# Email API Architecture

## Overview

The Email API is a separate microservice that handles all email-related functionality for BizOps360. This document explains the architecture, deployment, and code organization decisions.

## Architecture Decision

### Current Structure: Monorepo with Separate Deployment

The Email API follows a **monorepo pattern** where:
- **Code**: Lives in the same Go module (`github.com/bizops360/go-api`)
- **Deployment**: Deployed as a separate Cloud Run service in a separate GCP project
- **Isolation**: Runs in its own project (`bizops360-email-dev`) with independent scaling and security

### Why This Approach?

#### ✅ Advantages

1. **Code Reusability**
   - Shared types and interfaces (`internal/ports`, `internal/util`)
   - Common middleware and utilities
   - Single `go.mod` for dependency management

2. **Simplified Development**
   - One repository to clone
   - Shared CI/CD patterns
   - Consistent code style and tooling

3. **Independent Deployment**
   - Deploy email API without affecting main API
   - Separate scaling and resource allocation
   - Isolated security boundaries

4. **Cost Optimization**
   - Separate billing and quotas
   - Independent resource limits
   - Better cost tracking

5. **Security Isolation**
   - Email credentials in separate project
   - Independent IAM policies
   - Reduced blast radius if compromised

#### ❌ Alternative Approaches Considered

**Option 1: Separate Repository**
- ❌ Code duplication
- ❌ Harder to maintain shared code
- ❌ More complex CI/CD setup

**Option 2: Same Project, Same Service**
- ❌ No isolation
- ❌ Shared scaling limits
- ❌ Security concerns

**Option 3: Separate Repository + Shared Library**
- ❌ More complex dependency management
- ❌ Versioning challenges
- ❌ Overhead for small team

## Current Structure

```
go/
├── cmd/
│   ├── api/              # Main API service
│   │   └── main.go
│   └── email-api/        # Email API service
│       └── main.go
├── internal/
│   ├── http/
│   │   ├── handlers/     # Main API handlers
│   │   ├── emailapi/     # Email API handlers
│   │   │   ├── router.go
│   │   │   └── handlers/
│   │   └── middleware/   # Shared middleware
│   └── infra/
│       └── email/        # Email infrastructure
│           ├── gmail.go  # Gmail API implementation
│           └── client.go # HTTP client for email API
├── scripts/
│   ├── deploy-dev.sh           # Deploy main API
│   └── deploy-email-api-dev.sh # Deploy email API
└── Dockerfile.email-api.dev     # Email API Dockerfile
```

## Deployment Architecture

```
┌─────────────────────────────────────┐
│   bizops360-dev (GCP Project)      │
│                                     │
│   ┌─────────────────────────────┐  │
│   │  bizops360-api-go-dev       │  │
│   │  (Cloud Run Service)        │  │
│   │                             │  │
│   │  - Form Events              │  │
│   │  - Triggers                 │  │
│   │  - Stripe                   │  │
│   │  - Estimates                │  │
│   │  - Email Handler            │  │
│   │    └─> HTTP Client          │  │
│   └─────────────────────────────┘  │
└─────────────────────────────────────┘
              │
              │ HTTP + API Key
              │
              ▼
┌─────────────────────────────────────┐
│   bizops360-email-dev (GCP Project)│
│                                     │
│   ┌─────────────────────────────┐  │
│   │  bizops360-email-api-dev    │  │
│   │  (Cloud Run Service)        │  │
│   │                             │  │
│   │  - Email Endpoints          │  │
│   │  - Gmail Integration        │  │
│   │  - Rate Limiting            │  │
│   │  - API Key Auth            │  │
│   └─────────────────────────────┘  │
│                                     │
│   ┌─────────────────────────────┐  │
│   │  Secret Manager             │  │
│   │  - gmail-credentials-json   │  │
│   │  - email-api-key-dev        │  │
│   └─────────────────────────────┘  │
└─────────────────────────────────────┘
              │
              │ Domain-Wide Delegation
              │
              ▼
┌─────────────────────────────────────┐
│   Google Workspace                  │
│   - Gmail API                       │
│   - team@stlpartyhelpers.com        │
└─────────────────────────────────────┘
```

## Code Organization

### Shared Code

The following packages are shared between main API and email API:

- `internal/ports` - Interface definitions (`Mailer`, `SendEmailRequest`)
- `internal/util` - Utility functions (JSON, errors, request ID)
- `internal/http/middleware` - HTTP middleware (auth, logging, rate limiting)
- `internal/infra/log` - Logging infrastructure

### Email-Specific Code

- `cmd/email-api/main.go` - Email API entry point
- `internal/http/emailapi/` - Email API handlers and router
- `internal/infra/email/gmail.go` - Gmail API implementation
- `internal/infra/email/client.go` - HTTP client for calling email API

### Main API Email Integration

- `internal/http/handlers/email_handler.go` - Main API email handler
  - Uses `EmailServiceClient` if `EMAIL_SERVICE_URL` is set
  - Falls back to `GmailSender` if direct Gmail credentials available

## Deployment

### Email API Deployment

```bash
cd go
bash scripts/deploy-email-api-dev.sh
```

**What it does:**
1. Builds Docker image for email API
2. Pushes to Artifact Registry (`bizops360-email-dev`)
3. Deploys to Cloud Run in `bizops360-email-dev` project
4. Configures secrets and environment variables

### Main API Deployment

```bash
cd go
bash scripts/deploy-dev.sh
```

**What it does:**
1. Builds Docker image for main API
2. Pushes to Artifact Registry (`bizops360-dev`)
3. Discovers email API URL automatically
4. Retrieves email API key from secret
5. Deploys with `EMAIL_SERVICE_URL` and `EMAIL_SERVICE_API_KEY` env vars

## Environment Variables

### Email API Service

```bash
ENV=dev
LOG_LEVEL=debug
GMAIL_FROM=team@stlpartyhelpers.com
GMAIL_CREDENTIALS_JSON=<from Secret Manager>
SERVICE_API_KEY=<from Secret Manager>
```

### Main API Service

```bash
ENV=dev
LOG_LEVEL=debug
EMAIL_SERVICE_URL=https://bizops360-email-api-dev-xxx.run.app
EMAIL_SERVICE_API_KEY=<from Secret Manager>
```

## API Communication

### Authentication

Main API authenticates to Email API using:
- Header: `X-Api-Key: {EMAIL_SERVICE_API_KEY}`
- Key stored in Secret Manager (`email-api-key-dev`)

### Request Flow

1. Client → Main API: `POST /api/email/test`
2. Main API → Email API: `POST /api/email/send` (with API key)
3. Email API → Gmail API: Sends email via Gmail API
4. Email API → Main API: Returns success/error
5. Main API → Client: Returns response

## Benefits of Current Architecture

### 1. Separation of Concerns
- Email functionality isolated from main business logic
- Independent scaling and resource allocation
- Clear boundaries between services

### 2. Security
- Email credentials in separate project
- API key authentication between services
- Independent IAM policies

### 3. Maintainability
- Shared code reduces duplication
- Single repository simplifies development
- Consistent tooling and patterns

### 4. Scalability
- Email API can scale independently
- Different resource limits per service
- Isolated performance impact

### 5. Cost Management
- Separate billing for email service
- Better cost tracking
- Independent quota management

## Future Considerations

### When to Split Repository?

Consider splitting into separate repositories if:
- Team grows significantly (>10 developers)
- Services have very different release cycles
- Need independent versioning
- Different technology stacks

### When to Merge Services?

Consider merging if:
- Services always deploy together
- No independent scaling needs
- Security isolation not required
- Simpler deployment preferred

## Best Practices

1. **Keep Shared Code Minimal**
   - Only share truly common utilities
   - Avoid tight coupling between services

2. **Document Dependencies**
   - Clear interface contracts
   - Version shared types carefully

3. **Independent Testing**
   - Test email API independently
   - Mock email API in main API tests

4. **Monitor Separately**
   - Separate logging and metrics
   - Independent alerting

5. **Deploy Independently**
   - Email API changes don't require main API redeploy
   - Main API can work without email API (graceful degradation)

## References

- [Monorepo vs Multi-repo](https://www.atlassian.com/git/tutorials/monorepos)
- [Microservices Patterns](https://microservices.io/patterns/index.html)
- [GCP Project Organization](https://cloud.google.com/resource-manager/docs/cloud-foundation-framework/organize-resources)

