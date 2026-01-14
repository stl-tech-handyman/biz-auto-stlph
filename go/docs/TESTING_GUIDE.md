# Quote Engine Testing Guide

## Overview

The Quote Engine has a comprehensive testing system that allows you to:
- Test individual features in isolation
- Run all tests automatically
- Generate detailed reports for analysis
- Validate correctness across date ranges and scenarios

## Test Dashboard

Access the test dashboard at: **`http://localhost:8080/test-dashboard.html`**

The dashboard provides:
- **Individual Test Cards**: Click any test to run it individually
- **Run All Tests**: Execute all tests at once with a single click
- **Real-time Results**: See pass/fail status immediately
- **Detailed Reports**: Download JSON reports for analysis

## Available Tests

### 1. Date Variations Test
Tests all date scenarios:
- Today, tomorrow, 2-3 days (critical)
- 4-7 days (urgent)
- 8-14 days (high)
- 15-30 days (moderate)
- 31+ days (normal)

**What it validates:**
- Urgency level calculations
- Expiration date calculations
- Date normalization (midnight-based)

### 2. Deposit Booking Deadline Test
Calculates how soon we need the deposit from them to book the staffing reservation. Tests deposit deadline rules for different values of days until event:
- ≤3 days: Deposit deadline is midnight today
- 4-7 days: Deposit deadline is in 48 hours
- 8-14 days: Deposit deadline is in 3 days
- 15+ days: Deposit deadline is in 2 weeks

**What it validates:**
- Deposit deadlines are calculated correctly based on days until event
- Deadlines are in the future (unless event is today)
- Correct time calculations
- Proper formatting

### 3. Urgency Level Test
Tests urgency level logic:
- Critical (≤3 days)
- Urgent (4-7 days)
- High (8-14 days)
- Moderate (15-30 days)
- Normal (>30 days)

**What it validates:**
- Correct urgency assignment
- Boundary conditions

### 4. Deposit Calculation Test
Tests deposit amount calculations:
- Various total costs ($100, $200, $500, $1000, $2000, $5000)
- Deposit range (15-30% of total)
- Professional amounts (multiples of $50)

**What it validates:**
- Deposit is within expected range
- Professional rounding
- Consistency

### 5. Pricing & Rates Test
Tests pricing calculations:
- Regular dates
- Special dates (holidays)
- Base rates
- Hourly rates
- Total cost calculations

**What it validates:**
- Rates are positive
- Special date pricing
- Total cost accuracy

### 6. Email Template Test
Tests email HTML generation:
- Regular events
- Urgent events
- High-demand dates
- Returning clients

**What it validates:**
- HTML is generated
- All required fields present
- Template structure correct

### 7. Weather Forecast Test
Tests weather integration:
- Geocoding addresses
- Fetching forecasts
- Weather recommendations

**What it validates:**
- Service availability
- Forecast data structure
- Recommendations logic

### 8. Form Validation Test
Tests form field validations:
- Email format
- Date format
- Helper count
- Hours validation

**What it validates:**
- Input validation rules
- Error handling

## Running Tests

### Via Dashboard (Recommended)
1. Open `http://localhost:8080/test-dashboard.html`
2. Click any test card to run individually
3. Or click "Run All Tests" to execute everything

### Via API
```bash
# Run specific test
curl -X POST http://localhost:8080/api/test/date-variations \
  -H "Content-Type: application/json" \
  -d '{"dateRange": "all", "detailed": true, "generateReport": true}'

# Run all tests
curl -X POST http://localhost:8080/api/test/run-all \
  -H "Content-Type: application/json" \
  -d '{"detailed": true, "generateReport": true}'
```

### Via Unit Tests
```bash
# Run all unit tests
go test ./internal/util -v

# Run specific test
go test ./internal/util -run TestCalculateUrgencyLevel -v

# Run with coverage
go test ./internal/util -cover
```

## Test Reports

When `generateReport: true` is set, reports are saved to:
- **Location**: `./test-reports/`
- **Format**: JSON
- **Naming**: `{test-name}-{timestamp}.json`

Reports include:
- Test summary (total, passed, failed, warnings)
- Individual test results
- Expected vs actual values
- Full test data

## Analyzing Results

### Report Structure
```json
{
  "testName": "Date Variations",
  "timestamp": "2026-01-13T23:50:00Z",
  "summary": {
    "total": 12,
    "passed": 11,
    "failed": 1,
    "warnings": 0
  },
  "results": [
    {
      "name": "Today (0 days)",
      "status": "pass",
      "expected": "critical",
      "actual": "critical",
      "data": {
        "daysUntilEvent": 0,
        "urgencyLevel": "critical"
      }
    }
  ]
}
```

### Using Cursor Agent for Analysis
1. Run tests with `generateReport: true`
2. Open the report file
3. Ask Cursor Agent: "Analyze this test report and identify any issues"
4. Agent will evaluate correctness and suggest fixes

## Test Patterns

Each test follows this pattern:
1. **Setup**: Define test cases with expected values
2. **Execute**: Run the function/feature being tested
3. **Validate**: Compare actual vs expected
4. **Report**: Return structured results

### Example Test Pattern
```go
testCases := []struct {
    name     string
    input    int
    expected string
}{
    {"Today", 0, "critical"},
    {"Tomorrow", 1, "critical"},
}

for _, tc := range testCases {
    actual := CalculateUrgencyLevel(tc.input)
    if actual != tc.expected {
        // Report failure
    }
}
```

## Continuous Testing

### Automated Test Execution
The system supports:
- **Scheduled runs**: Can be integrated with CI/CD
- **On-demand**: Via dashboard or API
- **Report generation**: Automatic JSON reports
- **Error logging**: All failures logged

### Best Practices
1. **Run tests after changes**: Always test after modifying logic
2. **Check reports**: Review generated reports for issues
3. **Fix failures immediately**: Don't let failures accumulate
4. **Test edge cases**: Especially boundary conditions (0, 3, 7, 14, 30 days)

## Unit Tests

### Location
- `go/internal/util/expiration_calc_test.go` - Expiration & urgency tests
- `go/internal/infra/stripe/deposit_calc_test.go` - Deposit calculation tests
- `go/internal/infra/stripe/deposit_calc_test_extended.go` - Extended deposit tests

### Running Unit Tests
```bash
# All tests in a package
go test ./internal/util -v

# Specific test
go test ./internal/util -run TestCalculateUrgencyLevel -v

# With coverage
go test ./internal/util -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Troubleshooting

**Tests not running?**
- Check server is running: `curl http://localhost:8080/api/health`
- Verify routes are registered in router.go
- Check server logs for errors

**Tests failing?**
- Review test report JSON
- Check expected vs actual values
- Verify logic hasn't changed
- Run unit tests to isolate issues

**Reports not generating?**
- Check `test-reports/` directory exists
- Verify write permissions
- Check server logs for errors

## Next Steps

1. **Add more test cases**: Expand date ranges, edge cases
2. **Integration tests**: Test full quote email flow
3. **Performance tests**: Measure response times
4. **Load tests**: Test under high volume
5. **Automated CI/CD**: Integrate with deployment pipeline
