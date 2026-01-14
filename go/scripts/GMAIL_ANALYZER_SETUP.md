# Gmail Email Analyzer - Setup Guide

This script analyzes your Gmail emails to extract form submissions, track confirmations, and generate revenue analytics in Google Sheets.

## Features

- ✅ Analyzes 86,000+ emails from Gmail
- ✅ Filters out test submissions
- ✅ Deduplicates by client email (same email = one client)
- ✅ Tracks confirmation status (confirmed vs unconfirmed)
- ✅ Extracts pricing/rate information
- ✅ Breaks down by year, month, and client
- ✅ Calculates revenue with payout percentages (45% payout, 10% self-payout)
- ✅ Creates live Google Sheets with charts
- ✅ Automates processing of new emails

## Setup Instructions

### Step 1: Create Google Apps Script Project

1. Go to [script.google.com](https://script.google.com)
2. Click "New Project"
3. Name it "Gmail Email Analyzer"
4. Delete the default `myFunction` code
5. Copy the entire contents of `gmail-email-analyzer.gs` into the editor

### Step 2: Initialize BOS Folder and Spreadsheet

**No manual setup needed!** The script will automatically:

1. Create a **BOS** folder in your Google Drive (if it doesn't exist)
2. Create a **single Google Sheet** named "Email Revenue Analytics" inside the BOS folder
3. Create all necessary tabs (pages) within that sheet:
   - Raw Data
   - Clients
   - Monthly Revenue
   - Yearly Summary
   - Pattern Discovery
   - Client Matching
   - Processing Log
   - Sample Analysis

**To initialize, simply run:**
```javascript
initializeBOS()
```

This will:
- Create the BOS folder in Google Drive
- Create the spreadsheet with all tabs
- Set up all sheet structures
- Store the IDs for future use

**Note:** You only need to run `initializeBOS()` once. After that, the script remembers the folder and spreadsheet.

### Step 4: Customize Gmail Query (Optional)

The script searches for emails matching:
```javascript
const GMAIL_QUERY = 'from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission"';
```

Adjust this based on your form submission patterns. Common patterns:
- `from:zapier.com` - Zapier form submissions
- `from:forms.google.com` - Google Forms
- `subject:"New Lead"` - Specific subject patterns
- `has:attachment` - Emails with attachments

### Step 3: Authorize and Initialize

1. Click the "Run" button (▶️) in the Apps Script editor
2. Select `initializeBOS` function
3. Click "Run"
4. You'll be prompted to authorize:
   - Click "Review Permissions"
   - Choose your Google account
   - Click "Advanced" → "Go to [Project Name] (unsafe)"
   - Click "Allow"
5. The script will:
   - Create the BOS folder in Google Drive
   - Create the spreadsheet with all tabs
   - Set up all sheet structures
6. Check the execution log for the spreadsheet URL

### Step 4: Run Initial Analysis

1. Select `analyzeAllEmails` function
2. Click "Run"
3. **Note**: Processing 86,000 emails will take time. The script processes in batches of 500 threads.
4. You can monitor progress in the "Execution" tab
5. Check your Google Sheet - data will appear as it's processed

### Step 7: Set Up Automation

1. Select `setupTrigger` function
2. Click "Run"
3. This creates a daily trigger at 2 AM to process new emails
4. You can change the schedule by modifying `setupTrigger()` function

## Understanding the Output

### Raw Data Sheet
- All processed emails with extracted information
- Columns: Email ID, Date, From, Subject, Body Preview, Is Test, Is Confirmation, Client Email, Event Date, Total Cost, Rate, Hours, Helpers, Occasion, Status

### Clients Sheet
- One row per unique client email
- Shows: First/Last submission dates, Total submissions, Total revenue, Confirmed/Unconfirmed events

### Monthly Revenue Sheet
- Revenue breakdown by year and month
- Includes: Gross revenue, Payout (45%), Self payout (10% of 45%), Net revenue, Event counts, Average rate

### Yearly Summary Sheet
- Annual totals and summaries
- Same structure as monthly but aggregated by year

## Customization

### Adjust Payout Percentages

In the `processRawData()` function, find:
```javascript
const payout = grossRevenue * 0.45;  // 45% payout
const selfPayout = payout * 0.10;    // 10% of payout
```

Change these percentages as needed.

### Improve Email Pattern Matching

The `extractEventData()` and `extractPricingData()` functions use regex patterns. You can enhance these based on your specific email formats.

### Add More Analysis

You can add custom sheets and analysis by:
1. Adding a new sheet name to `SHEETS` object
2. Creating the sheet structure in `setupSheet()`
3. Adding processing logic in `processRawData()`

## Troubleshooting

### Script Times Out
- Gmail API has rate limits
- The script processes in batches
- If it times out, run `analyzeAllEmails()` multiple times - it will skip already processed emails (based on Email ID)

### Missing Data
- Check the Gmail query matches your email patterns
- Verify email format matches the extraction patterns
- Review the Raw Data sheet to see what was extracted

### Incorrect Pricing
- Review the `extractPricingData()` function
- Add more regex patterns based on your email format
- Check the Raw Data sheet to see what was extracted

## Advanced: Add Charts

After data is processed, you can add charts in Google Sheets:

1. Select the Monthly Revenue sheet
2. Insert → Chart
3. Choose "Line chart" or "Column chart"
4. Data range: Select the revenue columns
5. X-axis: Year/Month column
6. Y-axis: Revenue columns

## Security Notes

- The script only accesses emails matching your query
- All processing happens in Google's secure environment
- No data leaves Google's servers
- You can review all permissions in Apps Script settings

## Support

If you need help customizing the script:
1. Check the execution logs in Apps Script (View → Execution log)
2. Review the Raw Data sheet to see what was extracted
3. Adjust the extraction patterns based on your email format
