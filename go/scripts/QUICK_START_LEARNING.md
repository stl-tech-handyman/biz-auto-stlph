# Quick Start: Learning Email Analyzer

## Fast Track Setup (5 minutes)

### 1. Create Google Apps Script Project
- Go to [script.google.com](https://script.google.com)
- New Project → Name it "Gmail Email Analyzer Learning"
- Copy entire `gmail-email-analyzer-learning.gs` code

### 2. Initialize BOS Folder and Spreadsheet
- Run: `initializeBOS()`
- This automatically creates:
  - **BOS folder** in Google Drive
  - **Single Google Sheet** with all tabs (pages)
  - All sheet structures
- Check execution log for the spreadsheet URL

### 3. Run Pattern Discovery (Start Small!)
```javascript
discoverPatterns(100)  // Analyze first 100 emails
```

### 4. Review Results
- Check **Sample Analysis** sheet - see what was extracted
- Check **Pattern Discovery** sheet - see what patterns were found
- Look for "NOT FOUND" entries - these need pattern refinement

### 5. Refine Patterns (if needed)
- Edit **Pattern Discovery** sheet directly
- Add/remove keywords based on what you see
- Run: `updatePatternsManually()`

### 6. Process More Samples
```javascript
discoverPatterns(500)   // Learn more patterns
discoverPatterns(1000)  // Even more patterns
```

### 7. Process All Emails
```javascript
analyzeAllEmailsLearning()  // Process all 86,000 emails
```

### 8. Identify Clients
```javascript
identifyUniqueClients()  // Group by unique client email
```

### 9. Generate Reports
```javascript
processRawDataLearning()  // Create summary sheets
```

## Key Functions Reference

| Function | Purpose | When to Run |
|----------|---------|-------------|
| `discoverPatterns(n)` | Analyze sample emails, learn patterns | Start here! Run with 100, then 500, then 1000 |
| `updatePatternsManually()` | Save pattern edits from sheet | After editing Pattern Discovery sheet |
| `analyzeAllEmailsLearning()` | Process all emails with learned patterns | After patterns are stable |
| `identifyUniqueClients()` | Group emails by unique client | After full analysis |
| `processRawDataLearning()` | Create summary/revenue sheets | After client identification |

## What Gets Created

1. **Sample Analysis** - Sample emails with extraction results
2. **Pattern Discovery** - All patterns found (keywords, formats, etc.)
3. **Raw Data** - All processed emails with extracted data
4. **Client Matching** - Unique clients identified
5. **Clients** - Client summary (one per unique email)
6. **Monthly Revenue** - Revenue by month with payout calculations
7. **Yearly Summary** - Annual totals
8. **Processing Log** - Progress tracking

## Iterative Improvement Process

```
discoverPatterns(100)
    ↓
Review Sample Analysis & Pattern Discovery
    ↓
Refine patterns (edit sheet, run updatePatternsManually)
    ↓
discoverPatterns(500)  // More patterns
    ↓
Review again, refine if needed
    ↓
discoverPatterns(1000)  // Even more patterns
    ↓
Review, finalize patterns
    ↓
analyzeAllEmailsLearning()  // Process all 86K emails
    ↓
identifyUniqueClients()  // Group by client
    ↓
processRawDataLearning()  // Generate reports
    ↓
Review results, refine patterns if needed
    ↓
Set up automation (setupTrigger)
```

## Common Issues & Solutions

**"NOT FOUND" in Sample Analysis**
→ Check where the data appears in email body
→ Add new patterns to Pattern Discovery sheet
→ Update extraction functions

**Too many test emails not filtered**
→ Review test keywords in Pattern Discovery
→ Add more keywords you notice
→ Run updatePatternsManually()

**Client emails not being extracted**
→ Check Sample Analysis - see where emails appear
→ Look at "From" field patterns
→ Update extractClientEmailLearning() function

**Script times out**
→ Normal! Just run again - it skips processed emails
→ Check Processing Log to see progress
→ Process in smaller batches if needed

## Pro Tips

1. **Start small**: Always begin with 100 emails to learn patterns
2. **Review before scaling**: Don't process all 86K until patterns are good
3. **Iterate**: Each sample run teaches you more about your email format
4. **Manual refinement**: Sometimes manual pattern addition is faster than code changes
5. **Check the logs**: Processing Log shows exactly what's happening

## Next Steps

After initial analysis:
- Set up daily automation: `setupTrigger()`
- Create charts in Google Sheets
- Export data for further analysis
- Integrate with your CRM

The learning approach means you'll get better results as you process more emails and refine patterns!
