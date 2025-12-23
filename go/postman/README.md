# Postman Collection

## Quick Start

1. **Import Collection**: Import `BizOps360-Go-API.postman_collection.json` into Postman
2. **Create Environment**: Create a new environment with these variables:
   - `base_url_dev`: `https://bizops360-api-go-dev-nhrhozfuaq-uc.a.run.app`
   - `base_url_prod`: `https://bizops360-api-go-prod-XXXXX-uc.a.run.app` (update when prod is deployed)
   - `api_key`: Your API key (get it from Secret Manager - see below)
3. **Select Environment**: Select your environment in Postman
4. **Set API Key**: Update the `api_key` variable with your actual API key
5. **Test**: All endpoints are ready to test!

## Getting API Key

### Development
```bash
gcloud secrets versions access latest --secret="svc-api-key-dev" --project="bizops360-dev"
```

### Production
```bash
gcloud secrets versions access latest --secret="svc-api-key-prod" --project="bizops360-prod"
```

## Features

- ✅ All API endpoints included
- ✅ curl commands in every request description
- ✅ Example request bodies
- ✅ Environment variables configured
- ✅ Ready to import and test

## Adding New Endpoints

**IMPORTANT**: When adding a new endpoint, you MUST:

1. Add it to this Postman collection
2. Include a curl command in the request description
3. Update `docs/api/POSTMAN_COLLECTION.md`
4. Test it before committing

## Documentation

See [docs/api/POSTMAN_COLLECTION.md](../docs/api/POSTMAN_COLLECTION.md) for detailed documentation.


