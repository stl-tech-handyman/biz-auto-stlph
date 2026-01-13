# Quote Email Availability & Urgency Messaging - Brainstorm

## Goal
Add messaging to quote emails that:
1. Emphasizes pricing is locked but availability is subject to change
2. Creates urgency for holiday/busy seasons
3. Encourages taking action sooner without being pushy

## Current Context
- System already detects holidays: New Year's Day, Thanksgiving, Christmas Eve, Christmas Day, New Year's Eve
- EstimateResult includes `IsSpecialDate` and `SpecialLabel` fields
- QuoteEmailData currently doesn't include holiday/busy season info

## Messaging Options

### Option 1: General + Special Date (Recommended - Updated)
**For Regular Dates:**
> "Please note: While this quote's pricing is locked in, availability is subject to change. We recommend securing your date with a deposit as soon as possible to guarantee your event staffing."

**For Holiday/Surge Dates (Concise):**
> "Please note: While this quote's pricing is locked in, availability is subject to change.  
> **[Dark red, own line:]** Since your event falls during a high-demand period, dates fill up faster than usual.  
> We recommend securing your date with a deposit as soon as possible to guarantee your event staffing."

**Alternative (Even more concise):**
> "Please note: While this quote's pricing is locked in, availability is subject to change. Your event date falls during a high-demand period, so we recommend securing your date with a deposit as soon as possible."

### Option 2: More Direct
**For Regular Dates:**
> "Pricing is locked in, but availability changes daily. Secure your date now to guarantee your event staffing."

**For Holiday Dates:**
> "Pricing is locked in, but [Holiday Name] season dates fill up quickly. Secure your date now to guarantee your event staffing."

### Option 3: Friendly & Informative
**For Regular Dates:**
> "This quote's pricing is guaranteed, but please note that our availability changes as we receive bookings. To ensure we can staff your event, we recommend securing your date with a deposit when you're ready."

**For Holiday Dates:**
> "This quote's pricing is guaranteed, but please note that during the [Holiday Name] season, our availability fills up faster than usual. To ensure we can staff your event, we recommend securing your date with a deposit as soon as possible."

### Option 4: Value-Focused
**For Regular Dates:**
> "Your quote pricing is locked in for 72 hours. Availability is limited and changes daily, so we recommend securing your date with a deposit to guarantee your preferred helpers for your event."

**For Holiday Dates:**
> "Your quote pricing is locked in for 72 hours. Since your event falls during the [Holiday Name] season, availability is especially limited and fills up quickly. We recommend securing your date with a deposit as soon as possible to guarantee your preferred helpers."

## Recommended Approach: Option 1 (General + Special Date - Updated)

### Rationale:
- Professional and informative
- Not pushy or salesy
- Works for both holidays AND dynamic surge dates
- Laconic/concise wording
- Doesn't require specific holiday names (works for any high-demand period)
- Balances urgency with helpfulness
- Clear call-to-action

### Key Change:
Use "high-demand period" instead of specific holiday names. This works for:
- Fixed holidays (Thanksgiving, Christmas, etc.)
- Dynamic surge dates (any date marked as high-demand)
- Any date where pricing is adjusted due to demand

## Implementation Considerations

### 1. Data Flow
- Add `IsSpecialDate bool` to `QuoteEmailData` (we don't need SpecialLabel since we're using generic "high-demand period")
- Pass `estimate.IsSpecialDate` when creating `QuoteEmailData`
- This already exists in EstimateResult, just needs to flow through
- No need to differentiate between holidays vs. surge dates in messaging (both are "high-demand periods")

### 2. Holiday Detection Enhancement
Current holidays detected:
- New Year's Day (Jan 1)
- Thanksgiving (4th Thursday of November)
- Christmas Eve (Dec 24)
- Christmas Day (Dec 25)
- New Year's Eve (Dec 31)

Potential additions (would need new logic):
- **Easter** (varies by year - first Sunday after first full moon after March 21)
- **Graduation Season** (May-June, could detect by month)
- **Mother's Day** (2nd Sunday in May)
- **Father's Day** (3rd Sunday in June)
- **Memorial Day Weekend** (last Monday in May)
- **4th of July Weekend** (July 1-5)
- **Labor Day Weekend** (first Monday in September)

### 3. Template Placement
**Recommended location:** After the "Rates & Pricing" section, before "Secure Your Date" section

This placement:
- Comes after they've seen the pricing (value anchor)
- Sets context before the CTA
- Creates logical flow: pricing → availability concern → action

### 4. Visual Design
- Use a subtle background color (light gray or light yellow)
- Small italic text or regular text with slightly muted color
- Not too prominent (don't want to seem desperate)
- Could use a subtle border or padding to separate it
- **High-demand period line (when IsSpecialDate is true):**
  - Should be on its own line (separate paragraph)
  - Use dark red color (e.g., #cc0000 or #b91c1c)
  - Maintains readability while creating visual emphasis

## Example Template Structure

```
[Rates & Pricing Section]
...
[Pricing table]
...

[Availability Notice Section - NEW]
<p style="font-size: 11px; color: #666666; padding: 8px; background-color: #f9f9f9; border-left: 3px solid #d97706; margin: 8px 0;">
  Please note: While this quote's pricing is locked in, availability is subject to change.
  {{if .IsSpecialDate}}
  <br />
  <span style="color: #b91c1c;">Since your event falls during a high-demand period, dates fill up faster than usual.</span>
  {{end}}
  <br />
  We recommend securing your date with a deposit as soon as possible to guarantee your event staffing.
</p>

[Secure Your Date Section]
...
[Deposit button]
...
```

Note: Since we're using inline HTML generation (not Go templates), we'll use conditional logic in Go code instead.

## Final Recommended Messaging (Laconic)

**For Regular Dates:**
> "Please note: While this quote's pricing is locked in, availability is subject to change. We recommend securing your date with a deposit as soon as possible to guarantee your event staffing."

**For Holiday/Surge Dates (High-Demand Periods):**
> "Please note: While this quote's pricing is locked in, availability is subject to change.  
> **[In dark red, on its own line:]** Since your event falls during a high-demand period, dates fill up faster than usual.  
> We recommend securing your date with a deposit as soon as possible to guarantee your event staffing."

## Next Steps
1. ✅ Approved messaging (Option 1 - laconic version)
2. Add `IsSpecialDate bool` field to `QuoteEmailData`
3. Update code that creates `QuoteEmailData` to pass `estimate.IsSpecialDate`
4. Update template to conditionally show high-demand message when `IsSpecialDate` is true
5. Place message after "Rates & Pricing" section, before "Secure Your Date" section
