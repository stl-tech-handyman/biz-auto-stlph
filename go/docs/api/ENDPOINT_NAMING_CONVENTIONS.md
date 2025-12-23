# API Endpoint Naming Conventions

## Overview

This document defines the naming conventions for API endpoints, with special focus on distinguishing **orchestrated endpoints** from **simple/atomic endpoints**.

## Endpoint Types

### 1. Simple/Atomic Endpoints
These endpoints perform a **single, focused operation**.

**Pattern:** `/api/{service}/{resource}`

**Examples:**
- `POST /api/stripe/deposit` - Creates a deposit invoice only
- `POST /api/stripe/final-invoice` - Creates a final invoice only
- `POST /api/email/final-invoice` - Sends an email only

**Characteristics:**
- Single responsibility
- No orchestration of multiple operations
- Can be composed together by the caller
- Lower-level, more granular

### 2. Orchestrated Endpoints
These endpoints **coordinate multiple operations** in a single call, providing a complete workflow.

**Pattern:** `/api/{service}/{resource}/with-{action}` or `/api/{service}/{resource}/orchestrate`

**Examples:**
- `POST /api/stripe/deposit/with-email` - Creates deposit invoice AND sends email
- `POST /api/stripe/final-invoice/with-email` - Creates final invoice AND sends email

**Characteristics:**
- Multiple operations in sequence
- Transaction-like behavior (all or nothing, or graceful degradation)
- Higher-level, business-focused
- Reduces client complexity

## Naming Patterns

### Pattern 1: `/with-{action}` (Recommended for Email Orchestration)
Use when orchestrating with a specific action like email sending.

**Examples:**
- `/api/stripe/deposit/with-email`
- `/api/stripe/final-invoice/with-email`
- `/api/stripe/invoice/with-notification`

**When to use:**
- The orchestration adds a clear, specific action (email, notification, etc.)
- The action is a common pattern across multiple endpoints

### Pattern 2: `/orchestrate` (For Complex Orchestration)
Use when orchestrating multiple complex operations.

**Examples:**
- `/api/lead/orchestrate` - Processes lead: estimate + calendar + email + geocode
- `/api/booking/orchestrate` - Creates booking: invoice + calendar + email + CRM update

**When to use:**
- Multiple operations (3+)
- Complex business workflows
- Operations that don't fit a simple "with-X" pattern

### Pattern 3: `/complete` (Alternative for Full Workflows)
Use as an alternative to `/orchestrate` for complete end-to-end workflows.

**Examples:**
- `/api/booking/complete` - Complete booking workflow
- `/api/payment/complete` - Complete payment workflow

**When to use:**
- Full end-to-end business process
- All steps must succeed (transactional)

## Current Orchestrated Endpoints

| Endpoint | Type | Orchestrates |
|----------|------|--------------|
| `POST /api/stripe/deposit/with-email` | Email Orchestration | Creates deposit invoice + sends email |
| `POST /api/stripe/final-invoice/with-email` | Email Orchestration | Creates final invoice + sends email |

## Guidelines

### When to Create an Orchestrated Endpoint

1. **Common Workflow**: Multiple clients need the same sequence of operations
2. **Transaction Safety**: Operations should succeed or fail together
3. **Client Simplification**: Reduces client-side complexity and API calls
4. **Business Logic**: Represents a complete business process

### When to Keep Separate Endpoints

1. **Flexibility Needed**: Clients need to customize the workflow
2. **Different Use Cases**: Operations are used independently
3. **Granular Control**: Clients need to handle errors at each step
4. **Testing**: Easier to test individual operations

## Response Format

Orchestrated endpoints should return a response that includes:
- Status of each operation
- Any errors that occurred (even if overall succeeded)
- Results from each step

**Example:**
```json
{
  "ok": true,
  "message": "Final invoice created and email sent",
  "invoice": {
    "id": "in_123",
    "url": "https://invoice.stripe.com/...",
    "amount": 600.0
  },
  "email": {
    "sent": true,
    "error": null
  }
}
```

## Documentation

All orchestrated endpoints should be clearly marked in:
1. **OpenAPI Spec**: Tagged with `orchestrated: true` or similar
2. **Endpoint List**: Grouped separately or clearly labeled
3. **Swagger UI**: Include note about orchestration in description

## Future Considerations

As we add more orchestrated endpoints, consider:
- Creating an `/api/orchestrate/` namespace for complex workflows
- Standardizing error handling across orchestrated endpoints
- Adding retry logic for failed steps
- Transaction rollback capabilities

