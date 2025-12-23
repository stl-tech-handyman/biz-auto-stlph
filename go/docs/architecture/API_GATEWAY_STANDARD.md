# API & Gateway Architecture Standard (Tyk Edition)

**Version 1.0 — BizOps360 Architecture Standard**

## 📋 Purpose

This document establishes the standard for building, documenting, and exposing microservice APIs using Tyk API Gateway, OpenAPI specifications, and consistent service patterns.

---

## 🎯 Part 1: Tyk API Gateway — Engineer-Oriented Summary

**TYK API GATEWAY — ENGINEER-ORIENTED SUMMARY**

Tyk is a lightweight, high-performance, open-source API Gateway written in Go.

It provides routing, authentication, rate limiting, analytics (optional), and API lifecycle management.

It is designed to be simple, declarative, and suitable for microservice architectures.

### Key Features

- **Declarative API definitions** (JSON/YAML)
- **Built-in authentication**: API Keys, JWT, OAuth2, HMAC, Basic Auth
- **Rate limiting and quotas**
- **Middleware & plugins** (Go, JS, Python, gRPC)
- **Request/response transformation** (body mapping, header injection, rewrites)
- **Service discovery** (K8s, Consul, static)
- **Gateway + Dashboard** (optional) — dashboard is paid, gateway is free
- **Very fast** (Go-based), minimal latency
- **Cloud, hybrid, and open-source on-prem versions**

### Architecture Components

- **Tyk Gateway (core)** — open-source reverse proxy + API manager
- **Tyk Dashboard (optional)** — UI, analytics, developer portal
- **Tyk Pump** — moves analytics logs to DB (Mongo, Redis, etc.)

### Why Engineers Like Tyk

- Clean JSON configs
- Predictable routing
- Easy to automate
- Powerful middleware system
- Runs anywhere: Docker, Kubernetes, VM, local

### Good Use Cases

- Microservices with multiple backends
- API consolidation
- Lightweight centralized auth
- Building internal API catalog
- Developer-friendly self-hosted gateway

### Notable Advantages vs Kong

- Easier declarative configuration
- More friendly Go ecosystem
- No Lua (simpler debugging)
- Good for companies that want open-source but not NGINX complexity

---

## 🏗️ Part 2: API Architecture Standard

### 2.1 Microservice Structure Standard

Каждый сервис должен иметь следующую структуру:

```
/service-name
  /cmd
  /internal
  /pkg
  /api
    openapi.yaml
    openapi-ru.yaml          # Russian version (if needed)
    examples/
    mocks/
  /deploy
    docker/
    k8s/
    tyk/
      service-name.json      # Tyk API Gateway config
  /configs
  main.go
  README.md
```

**Обязательное:**

- `api/openapi.yaml` — всегда существует
- Ни один endpoint не существует без описания в OpenAPI
- Tyk config в `/deploy/tyk/service-name.json`

### 2.2 OpenAPI Standard

Минимальный набор должен содержать:

```yaml
openapi: 3.0.3
info:
  title: <Service Name> API
  version: 1.0.0

servers:
  - url: /<service-name>

paths:
  /resource:
    get:
      summary: Get resource list
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ResourceList'

components:
  schemas:
    Resource:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
```

**Правила:**

- Каждый endpoint должен быть описан
- Все тела запросов/ответов — через JSON Schema
- Versioning через `/v1`, `/v2` — обязательно при изменении контракта
- Описания должны быть минимально достаточными, но не пустыми
- **При любом изменении API — немедленно обновлять OpenAPI спецификацию**

### 2.3 Tyk Gateway Standard

Каждый сервис экспонируется через файл:

`/deploy/tyk/<service-name>.json`

**Шаблон:**

```json
{
  "name": "<SERVICE-NAME>",
  "api_id": "<UUID>",
  "org_id": "default",
  "definition": {
    "location": "api/openapi.yaml",
    "key": "openapi"
  },
  "listen_port": 8080,
  "protocol": "http",
  "enable_coprocess_auth": false,

  "proxy": {
    "target_url": "http://service-name:8080",
    "strip_listen_path": true,
    "listen_path": "/<service-name>/"
  },

  "auth": {
    "auth_header_name": "Authorization",
    "use_keyless": true
  },

  "rate_limit": {
    "rate": 100,
    "per": 1
  }
}
```

**Правила:**

- Каждый сервис → отдельный Tyk API config
- Конфиги должны быть храниться в Git
- Production deployments управляются CI/CD пайплайнами
- OpenAPI spec должен быть доступен для Tyk

### 2.4 API Naming Standard

**Формат:**

```
/<service>/<resource>
/<service>/<resource>/{id}
/<service>/<resource>/{id}/sub-resource
```

**Согласованная нотация:**

- Ресурсы — множественное число
- Идентификаторы — `{id}`
- Никакой бизнес-логики в URL (`/processDataNow` = плохо)
- Версионирование через `/v1`, `/v2` в пути

**Примеры:**

```
✅ /api/business/{businessId}/process-lead
✅ /api/stripe/deposit
✅ /v1/form-events
❌ /api/processDataNow
❌ /api/doSomething
```

### 2.5 Authentication Standard

**По умолчанию во всех сервисах:**

```
Authorization: Bearer <jwt>
```

**Token validation выполняется на уровне:**

- **Tyk Gateway** → проверяет JWT (public key)
- **Сервис** → доверяет заголовкам от Gateway (например `X-User-ID`)

**Для внутренних сервисов (без Tyk):**

```
X-Api-Key: <api-key>
```

**Правила:**

- API ключи хранятся в Google Secret Manager
- JWT токены валидируются на уровне Gateway
- Сервисы не должны валидировать токены самостоятельно (доверяют Gateway)

### 2.6 Logging & Monitoring Standard

**Каждый сервис должен:**

- Логировать `trace_id` / `request_id`
- Прокидывать `X-Request-ID`
- Отдавать health-check endpoint `/healthz` или `/api/health`

**Tyk подключен к:**

- **Prometheus** (gateway metrics)
- **Loki / ELK** (access logs)
- **Grafana dashboards** (latency, errors, RPS)

**Структурированное логирование:**

```go
logger.Info("request processed",
    "requestId", requestID,
    "businessId", businessID,
    "endpoint", "/api/business/process-lead",
    "status", "success",
    "duration", duration,
)
```

### 2.7 Developer Workflow Standard

**New endpoint workflow:**

1. Добавить описание в `openapi.yaml` (и `openapi-ru.yaml` если используется)
2. Сгенерировать Go types (опционально: `oapi-codegen`)
3. Реализовать handler
4. Обновить Tyk config (если новый path)
5. Запустить локально через Docker compose
6. Пройти тесты
7. Закоммитить и пройти CI

**Обязательные проверки:**

- [ ] OpenAPI spec обновлен
- [ ] OpenAPI YAML валиден
- [ ] Endpoint count обновлен
- [ ] Tyk config обновлен (если нужно)
- [ ] Unit tests написаны
- [ ] Логирование добавлено
- [ ] Error handling реализован

---

## 🤖 Part 3: AI Prompt Standard for Cursor/AI

**AI Prompt Template for API Work (Tyk + Go + OpenAPI)**

### Context

We use Tyk API Gateway and OpenAPI-first development.

Every microservice must include:

- `api/openapi.yaml`
- A Go implementation following the spec
- Tyk config in `/deploy/tyk/service.json`

### Rules to Follow

- Always start from OpenAPI
- Generate Go structs from schemas
- Return clean, idiomatic Go (Chi or Echo preferred)
- Generate corresponding Tyk config if needed
- Ensure versioning and naming standards are respected
- **Update OpenAPI spec for ANY API change**

### Task Template

```
Context:
We use Tyk API Gateway and OpenAPI-first development.
Every microservice must include:
- api/openapi.yaml
- a Go implementation following the spec
- Tyk config in /deploy/tyk/service.json

Rules to follow:
- Always start from OpenAPI
- Generate Go structs from schemas
- Return clean, idiomatic Go
- Generate corresponding Tyk config if needed
- Ensure versioning and naming standards are respected
- Update OpenAPI spec for ANY API change

Task:
<DESCRIBE WHAT YOU NEED BUILT>
```

### Example Usage

```
Context:
We use Tyk API Gateway and OpenAPI-first development.
Every microservice must include:
- api/openapi.yaml
- a Go implementation following the spec
- Tyk config in /deploy/tyk/service.json

Rules to follow:
- Always start from OpenAPI
- Generate Go structs from schemas
- Return clean, idiomatic Go
- Generate corresponding Tyk config if needed
- Ensure versioning and naming standards are respected
- Update OpenAPI spec for ANY API change

Task:
Create a new endpoint POST /api/notifications/send that:
- Accepts { to, subject, body, type }
- Sends notification via configured provider
- Returns { success, messageId, error }
- Requires API key authentication
- Rate limited to 50 requests per minute
```

---

## 📋 Compliance Checklist

Before deploying any API changes:

- [ ] OpenAPI spec updated (`openapi.yaml` and `openapi-ru.yaml` if used)
- [ ] OpenAPI YAML validated
- [ ] Endpoint count updated in OpenAPI description
- [ ] Tyk config updated (if new endpoint or path change)
- [ ] Go handler implemented following OpenAPI spec
- [ ] Unit tests written
- [ ] Error handling implemented
- [ ] Structured logging added
- [ ] Authentication configured correctly
- [ ] Rate limiting configured
- [ ] Health check endpoint exists
- [ ] Request ID middleware in place

---

## 🔄 Integration with Existing Standards

This standard integrates with:

- **[CODING_STANDARDS.md](../development/CODING_STANDARDS.md)** — Code quality and generalization rules
- **OpenAPI Specification** — `go/docs/api/openapi-ru.yaml`
- **Postman Collection** — `go/postman/BizOps360-Go-API.postman_collection.json`

**Key Integration Points:**

1. **OpenAPI Maintenance** (from CODING_STANDARDS.md):
   - Any API change MUST update OpenAPI spec
   - Validate YAML after changes
   - Update endpoint count

2. **Naming Standards**:
   - Follow API naming standard from this document
   - Use business ID routing: `/api/business/{businessId}/...`

3. **Documentation**:
   - OpenAPI is the source of truth
   - Postman collection should match OpenAPI
   - README files reference OpenAPI spec

---

## 📚 References

- [Tyk Documentation](https://tyk.io/docs/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [CODING_STANDARDS.md](../development/CODING_STANDARDS.md)
- [OpenAPI README](../api/OPENAPI_README_RU.md)

---

**Last Updated**: 2025-01-XX  
**Version**: 1.0.0  
**Maintainer**: BizOps360 Architecture Team



