# Postman Collection Naming Conventions

## Overview

This document defines naming conventions for Postman collections, environments, and requests to ensure consistency and maintainability across the BizOps360 project.

## Collection Naming

### Format
```
{Project}-{Service}-{API-Type}.postman_collection.json
```

### Components
- **Project**: `BizOps360` (always uppercase, no spaces)
- **Service**: Service name (e.g., `Go-API`, `Email-API`, `Maps-API`)
- **API-Type**: Optional, describes API type (e.g., `REST`, `GraphQL`)

### Examples
```
BizOps360-Go-API.postman_collection.json
BizOps360-Email-API.postman_collection.json
BizOps360-Maps-API.postman_collection.json
```

### Rules
- Use PascalCase for service names
- Use hyphens (`-`) as separators
- Keep names concise but descriptive
- Include version number if multiple versions exist: `BizOps360-Go-API-v2.postman_collection.json`

## Environment Naming

### Format
```
{Project}-{Service}-{Environment}.postman_environment.json
```

### Examples
```
BizOps360-Go-API-Dev.postman_environment.json
BizOps360-Go-API-Prod.postman_environment.json
BizOps360-Email-API-Dev.postman_environment.json
```

### Environment Variables Naming

#### Base URLs
```
{{base_url_dev}}
{{base_url_prod}}
{{base_url_staging}}
```

#### API Keys
```
{{api_key_dev}}
{{api_key_prod}}
{{email_api_key_dev}}
{{email_api_key_prod}}
```

#### Service URLs
```
{{email_service_url_dev}}
{{email_service_url_prod}}
{{maps_service_url_dev}}
```

### Rules
- Use lowercase with underscores for variable names
- Prefix with service name if multiple services (e.g., `email_api_key`, `maps_api_key`)
- Use descriptive suffixes (`_dev`, `_prod`, `_staging`)

## Request Naming

### Format
```
{HTTP_METHOD} {Resource} [{Action}]
```

### Examples
```
GET /api/health
POST /api/email/test
POST /api/email/send
GET /api/estimate
POST /api/stripe/deposit/calculate
```

### Folder Structure
```
Collection Name
├── Health
│   ├── GET /api/health
│   ├── GET /api/health/ready
│   └── GET /api/health/live
├── Email
│   ├── POST /api/email/test
│   ├── POST /api/email/send
│   └── POST /api/email/draft
├── Stripe
│   ├── POST /api/stripe/deposit/calculate
│   └── POST /api/stripe/test
└── Estimate
    ├── POST /api/estimate
    └── GET /api/estimate/special-dates
```

### Rules
- Start with HTTP method (uppercase): `GET`, `POST`, `PUT`, `DELETE`, `PATCH`
- Include full path: `/api/resource/action`
- Use descriptive action names when path is ambiguous
- Group related requests in folders
- Use PascalCase for folder names

## Collection Metadata

### Collection Info Structure
```json
{
  "info": {
    "name": "BizOps360 Go API",
    "description": "API collection for BizOps360 Go backend service",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "_postman_id": "unique-uuid-here"
  }
}
```

### Description Template
```markdown
# {Service Name} API Collection

## Overview
Brief description of the API service.

## Environments
- **Dev**: `{{base_url_dev}}`
- **Prod**: `{{base_url_prod}}`

## Authentication
Describe authentication method (API Key, OAuth, etc.)

## Endpoints
List main endpoint categories.

## Usage
1. Import collection
2. Import environment
3. Set environment variables
4. Start testing
```

## Folder Organization

### By Feature/Module
```
Collection
├── Authentication
├── Health & Status
├── Email
├── Payments (Stripe)
├── Estimates
└── Admin
```

### By Resource
```
Collection
├── Users
├── Orders
├── Products
└── Reports
```

### Best Practices
- Group related endpoints together
- Use clear, descriptive folder names
- Limit nesting depth (max 2-3 levels)
- Include a "Common" or "Utilities" folder for shared requests

## Request Examples

### Request Name
```
POST /api/email/test
```

### Request Description
```markdown
## Send Test Email

Sends a test email through the email API.

### Request Body
- `to` (required): Recipient email address
- `subject` (required): Email subject
- `html` (optional): HTML email body
- `text` (optional): Plain text email body

### Response
- `200 OK`: Email sent successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Invalid API key
```

### Pre-request Script Example
```javascript
// Set request ID for tracing
pm.request.headers.add({
    key: 'X-Request-Id',
    value: pm.variables.replaceIn('{{$randomUUID}}')
});
```

### Test Script Example
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has success field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('ok');
    pm.expect(jsonData.ok).to.be.true;
});
```

## Version Control

### File Naming
- Keep collection files in version control
- Use descriptive commit messages when updating collections
- Tag releases if collections are versioned

### Git Structure
```
go/postman/
├── BizOps360-Go-API.postman_collection.json
├── BizOps360-Email-API.postman_collection.json
├── BizOps360-Go-API-Dev.postman_environment.json
├── BizOps360-Go-API-Prod.postman_environment.json
├── NAMING_CONVENTIONS.md
└── README.md
```

## Checklist for New Collections

- [ ] Collection name follows convention: `{Project}-{Service}-{API-Type}`
- [ ] Environment files follow convention: `{Project}-{Service}-{Environment}`
- [ ] All requests have descriptive names with HTTP method
- [ ] Requests are organized in logical folders
- [ ] Environment variables use consistent naming
- [ ] Collection includes description and usage instructions
- [ ] Requests include descriptions and examples
- [ ] Test scripts validate responses
- [ ] Pre-request scripts set required headers
- [ ] Collection is added to version control

## Examples

### Complete Collection Structure
```
BizOps360-Go-API.postman_collection.json
├── Health
│   ├── GET /api/health
│   ├── GET /api/health/ready
│   └── GET /api/health/live
├── Email
│   ├── POST /api/email/test
│   └── POST /api/email/booking-deposit
├── Stripe
│   ├── POST /api/stripe/deposit/calculate
│   └── POST /api/stripe/test
└── Estimate
    ├── POST /api/estimate
    └── GET /api/estimate/special-dates
```

### Environment Variables
```json
{
  "base_url_dev": "https://bizops360-api-go-dev-xxx.run.app",
  "base_url_prod": "https://bizops360-api-go-prod-xxx.run.app",
  "api_key_dev": "dev-api-key-12345",
  "api_key_prod": "prod-api-key-xxxxx",
  "email_service_url_dev": "https://bizops360-email-api-dev-xxx.run.app",
  "email_api_key_dev": "email-api-key-dev-xxxxx"
}
```

## References

- [Postman Collection Format](https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html)
- [Postman Best Practices](https://learning.postman.com/docs/writing-scripts/script-references/test-examples/)
- [API Design Guidelines](https://restfulapi.net/)

