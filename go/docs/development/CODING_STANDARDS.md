# Coding Standards & Best Practices

**This document defines the mandatory coding standards for the BizOps360 Go API project.**

## 🚫 Hard Coding Standards - NEVER VIOLATE

### 1. **Generalization Over Specificity**

**❌ NEVER create overly specific utility functions:**
```go
// BAD - Too specific to Zapier
func ConsolidateOccasion(occasion, occasionAsYouSeeIt string) string { ... }
func ConsolidateRole(eventRole, eventRoleAsYouSeeIt string) string { ... }
func ProcessScheduleCall(scheduleCallStr string) bool { ... }
func FormatHelpersRequested(helpersStr string) int { ... }
func FormatDuration(durationStr string) float64 { ... }
```

**✅ ALWAYS create generalized, reusable functions:**
```go
// GOOD - General purpose, reusable
func ConsolidateWithFallback(primary, fallback, trigger, defaultValue string) string { ... }
func ParseBooleanFromText(text string, positiveIndicators ...string) bool { ... }
func ExtractFirstInteger(s string) int { ... }
func ExtractFirstFloat(s string) float64 { ... }
```

**Rule**: If you find yourself creating multiple functions that do the same thing with different names/contexts, create ONE generalized function instead.

### 2. **No Hard-Coded Values**

**❌ NEVER hard-code:**
- API endpoints
- Configuration values
- Magic numbers
- Default values that should be configurable

**✅ ALWAYS use:**
- Environment variables
- Configuration files
- Constants with clear names
- Dependency injection

### 3. **Single Responsibility Principle**

**❌ NEVER create functions that do multiple unrelated things:**
```go
// BAD
func ProcessLeadAndSendEmailAndCreateCalendar(data LeadData) error {
    // calculates estimate
    // sends email
    // creates calendar
    // geocodes address
}
```

**✅ ALWAYS separate concerns:**
```go
// GOOD
func CalculateEstimate(data LeadData) (*Estimate, error) { ... }
func SendEmail(email EmailRequest) error { ... }
func CreateCalendarEvent(event EventRequest) error { ... }
func GeocodeAddress(address string) (*Location, error) { ... }

// Orchestration happens at service layer
func (s *LeadService) ProcessLead(data LeadData) error {
    estimate, _ := CalculateEstimate(data)
    SendEmail(...)
    CreateCalendarEvent(...)
    GeocodeAddress(...)
}
```

### 4. **Proper Error Handling**

**❌ NEVER:**
- Ignore errors silently
- Use generic error messages
- Return `nil` errors when something actually failed
- Log errors without context

**✅ ALWAYS:**
- Handle every error explicitly
- Provide context in error messages
- Log errors with structured logging
- Return meaningful error types

```go
// BAD
result, _ := someFunction() // Error ignored!

// GOOD
result, err := someFunction()
if err != nil {
    logger.Error("failed to process data", 
        "error", err,
        "input", inputData,
        "step", "dataProcessing",
    )
    return fmt.Errorf("failed to process data: %w", err)
}
```

### 5. **Structured Logging with Levels**

**❌ NEVER:**
- Use `fmt.Println` or `log.Print`
- Log everything at the same level
- Log without context

**✅ ALWAYS:**
- Use structured logging (`slog`)
- Use appropriate log levels:
  - `Debug`: Detailed information for debugging (only in dev)
  - `Info`: General informational messages (important events)
  - `Warn`: Warning messages (non-critical issues)
  - `Error`: Error messages (failures that need attention)
- Include context in every log:
  - Request ID
  - User/Business ID
  - Operation being performed
  - Relevant data (sanitized)

```go
// BAD
fmt.Println("Processing lead")
log.Printf("Error: %v", err)

// GOOD
logger.Debug("processing lead",
    "businessId", businessID,
    "clientName", clientName,
    "requestId", requestID,
)
logger.Error("failed to process lead",
    "error", err,
    "businessId", businessID,
    "step", "emailSending",
    "requestId", requestID,
)
```

**Log Level Configuration:**
- Set via `LOG_LEVEL` environment variable
- Default: `info` in production, `debug` in development
- Must be easily switchable without code changes

### 6. **Modular Architecture**

**❌ NEVER:**
- Put business logic in handlers
- Create circular dependencies
- Mix infrastructure concerns with domain logic
- Create god objects/functions

**✅ ALWAYS:**
- Follow clean architecture layers:
  - **Domain**: Business entities and rules (no dependencies)
  - **Ports**: Interfaces (no dependencies)
  - **Infrastructure**: External services (implements ports)
  - **Services**: Business logic orchestration
  - **Handlers**: HTTP request/response handling
- Keep layers independent
- Use dependency injection
- One file, one responsibility

```
internal/
├── domain/          # Business entities (no deps)
├── ports/           # Interfaces (no deps)
├── infra/           # External services (implements ports)
├── services/         # Business logic (uses ports)
├── http/            # HTTP handlers (uses services)
└── util/            # Shared utilities (no deps)
```

### 7. **Comprehensive Testing**

**❌ NEVER:**
- Skip tests
- Write tests that test multiple things
- Use production data in tests
- Write tests that depend on external services

**✅ ALWAYS:**
- Write unit tests for ALL utility functions
- Write unit tests for ALL service methods
- Write integration tests for critical paths
- Use table-driven tests
- Mock external dependencies
- Aim for 90%+ code coverage

```go
// GOOD - Table-driven test
func TestExtractFirstInteger(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  int
    }{
        {"standard format", "I Need 2 Helpers", 2},
        {"just number", "5", 5},
        {"no number", "I Need Helpers", 0},
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
```

### 8. **Code Quality Standards**

**❌ NEVER:**
- Write functions longer than 50 lines
- Create functions with more than 5 parameters
- Use unclear variable names (`x`, `data`, `temp`)
- Leave TODO comments without context
- Comment out code (delete it or use version control)

**✅ ALWAYS:**
- Keep functions small and focused
- Use descriptive names
- Extract complex logic into separate functions
- Document public APIs
- Use meaningful variable names

```go
// BAD
func p(d string) (t time.Time, e error) { ... }

// GOOD
func ParseDate(dateString string) (time.Time, error) { ... }
```

### 9. **Business-Specific Routing**

**❌ NEVER:**
- Create endpoints that are hard-coded to one business
- Mix business logic in route handlers

**✅ ALWAYS:**
- Use business ID in routes: `/api/business/{businessId}/...`
- Load business configuration
- Allow different implementations per business
- Keep business logic in services, not handlers

```go
// GOOD
mux.HandleFunc("/api/business/{businessId}/process-lead", handler.HandleProcessLead)

func (h *Handler) HandleProcessLead(w http.ResponseWriter, r *http.Request) {
    businessID := extractBusinessIDFromPath(r.URL.Path)
    businessConfig, err := h.businessLoader.LoadBusiness(ctx, businessID)
    // ... use businessConfig for business-specific logic
}
```

### 10. **Documentation Standards**

**❌ NEVER:**
- Leave public functions undocumented
- Write comments that just repeat the code
- Use outdated documentation
- Make API changes without updating OpenAPI specification

**✅ ALWAYS:**
- Document all public functions
- Explain WHY, not WHAT
- Keep documentation up to date
- Include examples for complex functions
- **Update OpenAPI specification for ANY API change**

```go
// GOOD
// ConsolidateWithFallback consolidates a primary value with a fallback value.
// If primary contains the trigger word (case-insensitive), returns fallback if available, otherwise default.
// Otherwise returns primary if available, otherwise default.
// This is useful for handling "Other" options in forms where users can specify custom values.
func ConsolidateWithFallback(primary, fallback, trigger, defaultValue string) string { ... }
```

### 11. **OpenAPI Specification Maintenance**

**❌ NEVER:**
- Add/modify/remove API endpoints without updating OpenAPI spec
- Change request/response schemas without updating OpenAPI spec
- Change authentication requirements without updating OpenAPI spec
- Leave OpenAPI spec outdated or incomplete
- Deploy API changes without OpenAPI updates

**✅ ALWAYS:**
- **Update `go/docs/api/openapi-ru.yaml` for ANY API change:**
  - Adding new endpoints
  - Modifying existing endpoints
  - Changing request/response structures
  - Changing authentication requirements
  - Adding/removing parameters
  - Changing error responses
- Update endpoint count in OpenAPI description
- Add/update examples in OpenAPI spec
- Validate OpenAPI YAML after changes
- Update both Russian (`openapi-ru.yaml`) and English (if exists) versions
- **Update Tyk config** if using Tyk Gateway (see [API_GATEWAY_STANDARD.md](../architecture/API_GATEWAY_STANDARD.md))

**Workflow:**
1. Make API code changes
2. **IMMEDIATELY update OpenAPI specification**
3. Validate OpenAPI YAML
4. Update endpoint count if needed
5. Update Tyk config if needed
6. Commit all changes together

**Validation:**
```bash
# Validate OpenAPI YAML
python -c "import yaml; yaml.safe_load(open('go/docs/api/openapi-ru.yaml', encoding='utf-8'))"

# Or use swagger-cli
npx @apidevtools/swagger-cli validate go/docs/api/openapi-ru.yaml
```

**See Also:** [API Gateway Architecture Standard](../architecture/API_GATEWAY_STANDARD.md) for Tyk Gateway integration standards.

## 📋 Code Review Checklist

Before submitting code, ensure:

- [ ] All utility functions are generalized (not overly specific)
- [ ] No hard-coded values
- [ ] Proper error handling with context
- [ ] Structured logging with appropriate levels
- [ ] Code follows clean architecture layers
- [ ] All functions have unit tests
- [ ] Tests are table-driven where appropriate
- [ ] Business-specific logic uses business ID routing
- [ ] Functions are small and focused
- [ ] Variable names are descriptive
- [ ] Public functions are documented
- [ ] No commented-out code
- [ ] Code compiles without warnings
- [ ] **OpenAPI specification updated for any API changes**
- [ ] OpenAPI YAML validated and correct
- [ ] Endpoint count updated if endpoints added/removed

## 🎯 When Adding New Features

1. **Plan**: Identify if you need new functions or can reuse existing ones
2. **Generalize**: If creating utilities, make them reusable
3. **Layer**: Put code in the correct architectural layer
4. **Test**: Write comprehensive tests first (TDD)
5. **Log**: Add structured logging at each step
6. **Document**: Update documentation
7. **Update OpenAPI**: **MANDATORY** - Update OpenAPI spec for any API changes
8. **Review**: Self-review against this checklist

## 🔍 Common Violations to Watch For

1. **Creating Zapier-specific functions** → Use generalized string/number extraction
2. **Hard-coding business logic** → Use business config loading
3. **Mixing concerns** → Separate into appropriate layers
4. **Ignoring errors** → Always handle and log errors
5. **No tests** → Every function must have tests
6. **Generic error messages** → Include context in errors
7. **Single log level** → Use appropriate levels (debug/info/warn/error)
8. **API changes without OpenAPI update** → **ALWAYS update OpenAPI spec immediately**

---

**Remember**: Code quality matters. Write code that is maintainable, testable, and follows these standards. If you find yourself violating any of these rules, refactor immediately.

**Last Updated**: 2025-01-XX  
**Version**: 1.0.0

