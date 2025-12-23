# BizOps360 Go API - Complete Endpoints Summary

**Service URL:** https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app  
**Total Endpoints:** 15

---

## 📋 All Endpoints

### **1. Health & Info (4 endpoints - No Auth Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Root endpoint - Service information and available endpoints |
| `GET` | `/api/health` | Comprehensive health check |
| `GET` | `/api/health/ready` | Readiness probe for load balancers |
| `GET` | `/api/health/live` | Liveness probe for container orchestration |

### **2. Stripe (4 endpoints - API Key Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/stripe/deposit` | Generate booking deposit invoice |
| `GET` | `/api/stripe/deposit/calculate` | Calculate recommended deposit from estimate |
| `POST` | `/api/stripe/deposit/with-email` | Generate invoice and send email (end-to-end) |
| `POST` | `/api/stripe/test` | Test Stripe integration |

### **3. Estimate (2 endpoints - API Key Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/estimate` | Calculate event estimate based on date, duration, and helpers |
| `GET` | `/api/estimate/special-dates` | Get all special dates (holidays + surge dates) |

### **4. Email (2 endpoints - API Key Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/email/test` | Send test email |
| `POST` | `/api/email/booking-deposit` | Send booking deposit email |

### **5. Calendar (1 endpoint - No Auth Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/calendar/create` | Create a calendar event - Replicates Apps Script createEvent function |

### **6. Zapier (1 endpoint - No Auth Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/zapier/process-lead` | Process lead from Zapier - Replicates Apps Script flow (calculate estimate, send quote email, create calendar event, geocode address) |

### **7. V1 Pipeline (2 endpoints - No Auth Required)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/v1/form-events` | Process form submissions (e.g., from WPForms) |
| `POST` | `/v1/triggers` | Process trigger-based events (e.g., from Monday.com) |

---

## 🔐 Authentication

**API Key Required for:**
- All `/api/stripe/*` endpoints
- All `/api/estimate/*` endpoints
- All `/api/email/*` endpoints

**No Auth Required for:**
- `/` (root)
- `/api/health*` endpoints
- `/api/calendar/*` endpoints
- `/api/zapier/*` endpoints
- `/v1/*` endpoints

**Header:**
```
X-Api-Key: your-api-key-here
```

---

## 📝 Request Examples

### **POST /api/estimate**
```json
{
  "eventDate": "2025-06-15",
  "durationHours": 4.0,
  "numHelpers": 2
}
```

### **POST /api/stripe/deposit**
```json
{
  "email": "test@example.com",
  "name": "Test User",
  "estimatedTotal": 1000.0,
  "depositValue": null,
  "deposit": null,
  "helpersCount": 2,
  "hours": 4.0,
  "useTest": false,
  "dryRun": false,
  "mockStripe": false
}
```

### **POST /api/stripe/deposit/with-email**
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "eventType": "party",
  "eventDateTimeLocal": "2025-06-15 17:00",
  "eventDate": "2025-06-15",
  "helpersCount": 2,
  "hours": 4.0,
  "duration": 4.0,
  "estimate": 1000.0,
  "estimatedTotal": 1000.0,
  "depositValue": null,
  "useTest": false,
  "dryRun": false,
  "saveAsDraft": false
}
```

### **POST /v1/form-events**
```json
{
  "businessId": "stlpartyhelpers",
  "pipelineKey": "quote_and_deposit",
  "dryRun": false,
  "options": {
    "sendQuoteEmail": true
  },
  "fields": {
    "name": "John Doe",
    "email": "john@example.com",
    "event_date": "2025-06-15",
    "duration_hours": 4,
    "num_helpers": 2
  }
}
```

### **POST /v1/triggers**
```json
{
  "source": "monday",
  "businessId": "stlpartyhelpers",
  "triggerKey": "send_renewal_offer",
  "pipelineKey": "renewal_followup",
  "resource": {
    "type": "monday_item",
    "boardId": 123456789,
    "itemId": 987654321
  },
  "payload": {
    "event_date": "2025-05-10",
    "client_email": "client@example.com"
  }
}
```

### **POST /api/calendar/create**
```json
{
  "calendarId": "c_f8c0098141f20b9bcb25d5e3c05d54c450301eb4f21bff9c75a04b1612138b54@group.calendar.google.com",
  "clientName": "John Doe",
  "occasion": "Birthday Party",
  "guestCount": 50,
  "eventDate": "2025-07-10",
  "eventTime": "4:00 PM",
  "phone": "314-555-5555",
  "location": "2300 Hitzert Ct, Fenton, MO 63026",
  "numHelpers": 2,
  "duration": 5.0,
  "totalCost": 400.0,
  "emailId": "john@example.com",
  "threadId": "",
  "dataSource": "zapier",
  "status": "Pending"
}
```

### **POST /api/zapier/process-lead**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email_address": "john@example.com",
  "phone_number": "3145555555",
  "event_date": "2025-07-10",
  "event_time": "4:00 PM",
  "event_location": "2300 Hitzert Ct, Fenton, MO 63026",
  "helpers_requested": "I Need 2 Helpers",
  "for_how_many_hours": "for 5 Hours",
  "occasion": "Birthday Party",
  "guests_expected": "50",
  "dryRun": false
}
```

---

## ✅ Postman Collection Status

**All 15 endpoints are included in:** `BizOps360-Go-API.postman_collection.json`

**Updated:**
- ✅ Base URL corrected to: `https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app`
- ✅ Request body examples updated to match actual handlers
- ✅ Query parameters documented
- ✅ All endpoints verified against router.go

---

**Last Updated:** $(date)

