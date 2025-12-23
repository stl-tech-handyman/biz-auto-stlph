# Go API Documentation

This directory contains all documentation for the BizOps360 Go API.

## Structure

```
docs/
├── README.md                    # This file
├── api/                         # API documentation
│   ├── openapi-ru.yaml         # OpenAPI specification (Russian)
│   ├── OPENAPI_README_RU.md   # OpenAPI usage guide
│   ├── POSTMAN_COLLECTION.md   # Postman collection guide
│   ├── API_COMPARISON.md       # Go vs JS API comparison
│   └── API_KEY_AUTHENTICATION_GUIDE.md  # API key authentication guide
├── architecture/                # Architecture standards
│   ├── README.md               # Architecture documentation index
│   └── API_GATEWAY_STANDARD.md # Tyk Gateway & API architecture standard
├── deployment/                  # Deployment documentation
│   ├── DEPLOY_INSTRUCTIONS.md  # How to deploy
│   ├── ENVIRONMENTS.md         # Environment configuration
│   └── README_ENVIRONMENTS.md  # Environment setup guide
├── development/                 # Development documentation
│   ├── CODING_STANDARDS.md     # Coding standards & best practices
│   ├── GO_API_IMPLEMENTATION.md  # Go API implementation summary
│   ├── GO_BACKEND_SETUP.md      # Go backend setup summary
│   └── GO_DEV_PROD_SETUP.md     # Dev/Prod setup summary
├── TEST_COVERAGE_REPORT.md     # Test coverage report
└── TEST_SUMMARY.md             # Test summary

postman/
└── BizOps360-Go-API.postman_collection.json  # Postman collection
```

## Quick Links

### API Documentation
- [OpenAPI Specification](api/openapi-ru.yaml) - Complete API specification in Russian (OpenAPI 3.0.3)
- [OpenAPI README](api/OPENAPI_README_RU.md) - How to use and validate the OpenAPI spec
- [Postman Collection Guide](api/POSTMAN_COLLECTION.md) - How to use and maintain the Postman collection
- [API Comparison](api/API_COMPARISON.md) - Comparison between Go and JavaScript APIs
- [API Key Authentication](api/API_KEY_AUTHENTICATION_GUIDE.md) - API key authentication guide

### Architecture Standards
- [API Gateway Standard](architecture/API_GATEWAY_STANDARD.md) - Tyk Gateway & API architecture standard
- [Architecture README](architecture/README.md) - Architecture documentation index

### Development
- [Coding Standards](development/CODING_STANDARDS.md) - **MANDATORY** coding standards and best practices
- [Go API Implementation](development/GO_API_IMPLEMENTATION.md) - Implementation summary
- [Go Backend Setup](development/GO_BACKEND_SETUP.md) - Backend setup summary
- [Dev/Prod Setup](development/GO_DEV_PROD_SETUP.md) - Environment setup summary

### Deployment
- [Deployment Instructions](deployment/DEPLOY_INSTRUCTIONS.md) - Step-by-step deployment guide
- [Environments](deployment/ENVIRONMENTS.md) - Environment configuration details
- [Environment Setup](deployment/README_ENVIRONMENTS.md) - Setting up dev/prod environments

### Testing
- [Test Coverage Report](TEST_COVERAGE_REPORT.md) - Detailed test coverage information
- [Test Summary](TEST_SUMMARY.md) - Overview of test suite

## Postman Collection

The Postman collection is located at `go/postman/BizOps360-Go-API.postman_collection.json`.

**Quick Start:**
1. Import the collection into Postman
2. Create an environment with variables:
   - `base_url_dev`: `https://bizops360-api-go-dev-nhrhozfuaq-uc.a.run.app`
   - `api_key`: Your API key from Secret Manager
3. Select your environment and start testing

See [Postman Collection Guide](api/POSTMAN_COLLECTION.md) for detailed instructions.

## Adding New Endpoints

When adding a new endpoint:

1. **Update OpenAPI Specification**: **MANDATORY** - Update `go/docs/api/openapi-ru.yaml`
   - Add endpoint to `paths`
   - Define request/response schemas
   - Update endpoint count in description
   - Validate YAML: `python -c "import yaml; yaml.safe_load(open('go/docs/api/openapi-ru.yaml', encoding='utf-8'))"`
2. **Add to Postman Collection**: Update `go/postman/BizOps360-Go-API.postman_collection.json`
3. **Include curl Command**: Add curl command to the request description
4. **Update Documentation**: Update `api/POSTMAN_COLLECTION.md` with the new endpoint
5. **Update Tyk Config** (if using Tyk Gateway): Update `/deploy/tyk/service.json` if needed
6. **Test**: Verify the endpoint works in Postman

**See [Coding Standards](development/CODING_STANDARDS.md) and [API Gateway Standard](architecture/API_GATEWAY_STANDARD.md) for detailed requirements.**

## Getting API Keys

### Development
```bash
gcloud secrets versions access latest --secret="svc-api-key-dev" --project="bizops360-dev"
```

### Production
```bash
gcloud secrets versions access latest --secret="svc-api-key-prod" --project="bizops360-prod"
```

## Contributing

When updating documentation:
- Keep curl commands in sync with Postman requests
- Update both Postman collection and markdown docs
- Test all examples before committing
- Follow the existing format and structure

