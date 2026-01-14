# Iterative Email Analysis Workflow

This guide explains how to use the learning analyzer to iteratively improve email pattern recognition and client identification.

## The Problem

When analyzing 86,000 emails, you can't know all patterns upfront. The learning approach:
1. **Discovers patterns** by analyzing samples
2. **Lets you review and refine** extraction rules
3. **Processes all emails** with learned patterns
4. **Continuously improves** as more data is processed

## Step-by-Step Workflow

### Phase 1: Pattern Discovery (Start Here!)

1. **Run `discoverPatterns(100)`**
   - Analyzes first 100 emails
   - Creates two sheets:
     - **Sample Analysis**: Shows what was extracted from each email
     - **Pattern Discovery**: Lists all patterns found (keywords, formats, etc.)

2. **Review the Sample Analysis sheet**
   - Check if client emails are being found correctly
   - Verify event dates are extracted properly
   - Confirm pricing/rates are captured
   - Note any patterns you see

3. **Review the Pattern Discovery sheet**
   - See which test keywords appear most frequently
   - Identify confirmation patterns
   - Check common email domains
   - Review date/cost/rate formats

4. **Manually refine patterns** (if needed)
   - Edit the Pattern Discovery sheet directly
   - Add/remove keywords based on your review
   - Run `updatePatternsManually()` to save changes

### Phase 2: Process Larger Sample

5. **Run `discoverPatterns(500)` or `discoverPatterns(1000)`**
   - Processes more emails to find more patterns
   - Updates the Pattern Discovery sheet
   - Review again and refine

6. **Repeat until patterns are stable**
   - You'll see the same patterns appearing consistently
   - This means your extraction rules are good

### Phase 3: Full Analysis

7. **Run `analyzeAllEmailsLearning()`**
   - Processes ALL emails using learned patterns
   - Creates Raw Data sheet with all extracted information
   - May take several hours for 86,000 emails (processes in batches)

8. **Monitor progress**
   - Check Processing Log sheet
   - Script processes in batches of 500 threads
   - If it times out, run again - it skips already processed emails

### Phase 4: Client Identification

9. **Run `identifyUniqueClients()`**
   - Analyzes the Raw Data
   - Creates Client Matching sheet showing:
     - Unique client emails
     - Alternative emails (same person, different address)
     - Submission counts per client
     - Confidence scores

10. **Review Client Matching sheet**
    - Check if same clients are being grouped correctly
    - Look for patterns in alternative emails
    - Note any clients that should be merged

### Phase 5: Refinement Loop

11. **Review Raw Data sheet**
    - Look for emails where extraction failed
    - Identify new patterns you didn't catch
    - Note any systematic issues

12. **Update patterns based on findings**
    - Add new keywords to Pattern Discovery
    - Update extraction functions if needed
    - Re-run `updatePatternsManually()`

13. **Re-process if needed**
    - If you find major issues, you can:
      - Clear the Raw Data sheet
      - Update patterns
      - Re-run `analyzeAllEmailsLearning()`

## Understanding the Output Sheets

### Sample Analysis Sheet
- Shows extraction results for sample emails
- Columns show what was found (or "NOT FOUND")
- Use this to see extraction quality

### Pattern Discovery Sheet
- Lists all discovered patterns with frequency
- Organized by type (Test Keywords, Confirmation Keywords, etc.)
- Higher frequency = more common pattern

### Raw Data Sheet
- All processed emails with extracted data
- One row per email
- Use filters to find issues

### Client Matching Sheet
- Unique clients identified
- Shows alternative emails found
- Submission counts and date ranges

### Processing Log Sheet
- Tracks progress during full analysis
- Shows batch numbers and email counts
- Helps identify where processing stopped

## Tips for Better Results

### Improving Client Email Extraction

If client emails aren't being found:
1. Check Sample Analysis - see where emails appear in the body
2. Look for patterns in the "From" field
3. Check if emails are in HTML vs plain text
4. Update `extractClientEmailLearning()` function with new patterns

### Improving Date Extraction

If event dates aren't being found:
1. Review date formats in Pattern Discovery
2. Check if dates are in specific fields (like "Event Date:")
3. Add new regex patterns to `extractEventDateLearning()`

### Improving Pricing Extraction

If costs/rates aren't being found:
1. Check cost patterns in Pattern Discovery
2. Look for currency symbols and formats
3. Verify if pricing is in tables or specific sections
4. Update `extractPricingDataLearning()` with new patterns

### Better Test Detection

If test submissions aren't being filtered:
1. Review test keywords in Pattern Discovery
2. Check Sample Analysis for false positives/negatives
3. Add domain-based rules (e.g., "test@example.com")
4. Update `detectTestLearning()` function

## Advanced: Custom Extraction Rules

You can add custom extraction logic based on your specific email formats:

```javascript
function extractCustomField(body, subject) {
  // Add your custom pattern matching here
  // Example: Extract from specific email template
  const match = body.match(/Your Custom Pattern: (.*)/);
  return match ? match[1] : '';
}
```

## Troubleshooting

### Script Times Out
- It processes in batches, so just run again
- Check Processing Log to see where it stopped
- Already processed emails are skipped

### Low Extraction Rate
- Review Sample Analysis to see what's missing
- Check Pattern Discovery for common formats
- Manually add patterns you notice

### Too Many False Positives
- Review Sample Analysis for incorrect extractions
- Refine patterns to be more specific
- Add negative patterns (exclude certain formats)

### Clients Not Being Deduplicated
- Check Client Matching sheet
- Look for alternative email patterns
- You may need to manually merge some clients

## Next Steps After Analysis

Once you have clean data:
1. Run `processRawDataLearning()` to create summary sheets
2. Create charts in Google Sheets
3. Set up automation with `setupTrigger()`
4. Review and refine monthly

## Continuous Improvement

The beauty of this approach:
- **Start small**: Analyze 100 emails, learn patterns
- **Scale up**: Process all 86,000 with learned patterns
- **Refine**: Update patterns as you find issues
- **Automate**: Set up daily processing of new emails
- **Iterate**: Continuously improve extraction quality

Each time you process new emails, you'll discover new patterns and can refine the extraction rules!
