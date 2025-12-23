# Testing Final Invoice Functionality

## Quick Start (3 Steps)

### Step 1: Start the Server

**PowerShell:**
```powershell
cd go
.\start-local.ps1
```

**Git Bash:**
```bash
cd go
./start-local.sh
```

The server will start on `http://localhost:8080`

### Step 2: Test It

**Option A: HTML Test Page (Easiest - Recommended)**
1. Double-click `test-final-invoice.html` (or open it in your browser)
2. Enter your email address
3. Click "🚀 Test Full Flow (Both)"
4. Check the results on the page

**Option B: Batch File**
1. Double-click `test-final-invoice-quick.bat`
2. Enter your email when prompted
3. Check the results in the console

**Option C: PowerShell Script**
```powershell
.\test-final-invoice.ps1
```

### Step 3: Check Results

- **Invoice**: Check your Stripe dashboard at https://dashboard.stripe.com/invoices
- **Email**: Check your email inbox (if email service is configured)

## Manual Testing with curl

### Create Final Invoice
```bash
curl -X POST http://localhost:8080/api/stripe/final-invoice \
  -H "X-Api-Key: test-api-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"email":"your-email@example.com","name":"Test Customer","totalAmount":1000.0,"depositPaid":400.0}'
```

### Send Email
```bash
curl -X POST http://localhost:8080/api/email/final-invoice \
  -H "X-Api-Key: test-api-key-12345" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Customer","email":"your-email@example.com","totalAmount":1000.0,"depositPaid":400.0,"remainingBalance":600.0,"invoiceUrl":"https://invoice.stripe.com/i/..."}'
```

## Troubleshooting

**Server not running?**
- Make sure you ran `start-local.ps1` or `start-local.sh` first
- Check that port 8080 is not in use

**CORS errors in browser?**
- The server should have CORS enabled, but if you see errors, try using the batch file or PowerShell script instead

**Email not sending?**
- Make sure `GMAIL_CREDENTIALS_JSON` or `EMAIL_SERVICE_URL` is configured
- The invoice will still be created in Stripe even without email

**Stripe errors?**
- Verify your Stripe key is set correctly
- Check Stripe dashboard for any errors

