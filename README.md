## STL Party Helpers – Google Apps Script Backend

Modular backend for internal operations automation using Google Apps Script.  
Integrates with Zapier, Gmail, Monday CRM, Google Sheets, and Google Maps.

---

## 📁 Project Structure (by Business System)

Each business function lives in its own folder:

```
src/
├── auth/                 # Token validation, secured access
├── email/                # Email sending, templates, forwarding
├── estimates/            # Quote estimation logic
├── geo/                  # Address → lat/lng conversion via Google Maps
├── stripe/               # Stripe integration (deposit handling)
├── health/               # Pingable healthchecks & diagnostics
├── router/               # Versioned GET/POST routing (doGet / doPost)
├── shared/               # Constants, enums, and shared utilities
├── utils/                # Generic utilities (e.g., request parsing, logging)
└── tests/                # All unit and integration tests
```

---

## 🌐 HTTP Routing

The entry point is `doGet(e)` inside `router/httpActions.router.v1.gs`.

It performs:

- API token validation
- Action parameter routing
- Version control (e.g. `v=1`, `v=2`)

Example:

```js
function doGet(e) {
  const action = getAction(e);                  // utils.http
  const token = extractTokenFromRequest(e);     // utils.auth
  const version = getApiVersion(e);             // utils.http

  if (!isValidZapierToken(token)) return text("Unauthorized");

  switch (version) {
    case 2: return handleGetV2(action, e);
    case 1:
    default: return handleGetV1(action, e);
  }
}
```

---

## 🔐 Security: API Token

All public requests must include a valid `api_token`.

### Script Property (required):
- `ZAPIER_TOKEN` – used to authenticate Zapier or any external system

### Validation Logic:
```js
function isValidZapierToken(token) {
  const expected = PropertiesService.getScriptProperties().getProperty("ZAPIER_TOKEN");
  return token && token === expected;
}
```

---

## 🌍 External API Keys

Script expects additional keys in the **Script Properties**:

| Key                  | Purpose                                |
|----------------------|----------------------------------------|
| `ZAPIER_TOKEN`       | Validates incoming Zapier requests     |
| `GOOGLE_MAPS_API_KEY`| Used in address-to-geo API calls       |

To set: `Apps Script → Project Settings → Script Properties`

---

## 🧪 Testing Philosophy

We use a 3-level testing strategy:

1. **Unit Tests**: Test core business logic in isolation
2. **Integration Tests**: Simulate routed `doGet()` calls
3. **Router Tests**: Validate handlers are correctly triggered per action

All test runners live in:

- `tests/`
- Or domain-specific folders (e.g. `email/tests.email.gs`)

To run all tests:

```js
test_all();
```

---

## 📦 Business Domain Code Layout

Each functional domain (like email) is structured as:

```
email/
├── core.email.gs           # Core logic: sendEmail()
├── handler.email.v1.gs     # API-facing logic
├── utils.email.gs          # Optional helpers (template builders)
├── tests.email.gs          # All tests for this module
```

Same pattern applies for `geo`, `stripe`, etc.

---

## ⚙️ CLASP: Apps Script Local Development

### Setup:
```bash
npm install -g @google/clasp
```

### Login:
```bash
clasp login
```

### Clone Existing Project:
```bash
clasp clone <SCRIPT_ID>
```

### Edit `.clasp.json`:
```json
{
  "scriptId": "SCRIPT_ID",
  "rootDir": "src"
}
```

### Push Changes:
```bash
clasp push
```

---

## ✅ Naming Conventions

| Thing            | Convention Example               |
|------------------|----------------------------------|
| Files            | `core.email.gs`, `handler.geo.v1.gs` |
| GET Actions      | `PublicGetActions.TEST_HEALTHCHECK` |
| Params           | `?action=SEND_QUOTE&v=2&api_token=abc123` |
| Versioned routing| `handleGetV1()`, `handleGetV2()` |

---

## 🧠 Versioning Strategy

API version is controlled via URL query:

- `?v=1` (default) → `handleGetV1`
- `?v=2` → `handleGetV2`

Use helpers:
```js
const version = getApiVersion(e); // from utils.http
```

You may optionally version:
- Handler files: `handler.email.v1.gs`
- Tests: `tests.email.v1.gs`
- Router: `httpActions.router.v1.gs`

---

## 🚀 Future Enhancements

- Add `doPost()` versioned router
- Add stackdriver-compatible logging
- Add webhook signature verification (Stripe)
- Add Gmail / Maps mocking for isolated testing