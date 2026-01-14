# Email Classification & Conversation Tracking Features

## New Features Added

### 1. Email Classification (STLPH vs Other)

The script now automatically classifies each email as:
- **STLPH**: Business-related emails (form submissions, quotes, bookings, etc.)
- **Other**: Marketing emails, personal emails, school emails, etc.

**How it works:**
- Analyzes subject and body for STLPH keywords (party, event, helpers, booking, quote, etc.)
- Checks for non-STLPH keywords (school, marketing, newsletter, etc.)
- Checks sender domain (zapier, forms, stlpartyhelpers, etc.)
- Uses scoring system to classify

**Customization:**
Edit these arrays in the script:
```javascript
const STLPH_KEYWORDS = ['party', 'event', 'helpers', ...];
const NON_STLPH_KEYWORDS = ['school', 'pta', 'marketing', ...];
```

### 2. Conversation/Thread Identification

The script groups related emails into conversations:
- **Same thread + same client email** = one conversation
- Tracks back-and-forth exchanges
- Identifies multi-email leads

**Conversation ID Format:**
- `{threadId}_{clientEmail}` - Groups emails in same thread with same client
- Falls back to `{threadId}_{subjectHash}` if no client email

### 3. Average Emails per Client

Calculates:
- **Total emails per client**
- **Number of conversations per client**
- **Average emails per conversation** = Total emails ÷ Conversations

This helps identify:
- Clients with many back-and-forth exchanges
- Simple one-email leads vs complex multi-email conversations
- Communication patterns

## New Sheet Tabs

### Conversations Tab
Shows all identified conversations with:
- Conversation ID
- Client Email
- Subject
- First/Last email dates
- Email count (total, STLPH, Other)
- Status (STLPH Lead vs Other)
- Outcome (Confirmed vs Pending)

### Email Classification Tab
Summary statistics:
- Email Type (STLPH vs Other)
- Count and percentage
- Average per client
- Examples and notes

### Updated Clients Tab
Now includes:
- **STLPH Emails**: Count of business-related emails
- **Other Emails**: Count of non-business emails
- **Conversations**: Number of unique conversations
- **Avg Emails per Conversation**: Average emails per conversation

### Updated Raw Data Tab
New columns:
- **Thread ID**: Gmail thread identifier
- **Email Type**: STLPH or Other
- **Conversation ID**: Groups related emails
- **Message Number**: Position in conversation thread

## Understanding the Data

### Email Classification Logic

**STLPH emails are identified by:**
1. Keywords in subject/body (party, event, booking, quote, etc.)
2. Sender domain (zapier, forms, stlpartyhelpers)
3. Higher STLPH keyword score than non-STLPH keywords

**Other emails are:**
- Marketing/newsletters
- Personal emails
- School/PTA emails
- Receipts/confirmations from other services
- Spam

### Conversation Grouping

**Same conversation = same thread + same client:**
- Initial form submission
- Your reply
- Client's follow-up
- Your confirmation
- All grouped as one conversation

**Different conversations:**
- Same client, different thread = different conversation
- Different client, same thread = different conversation

### Average Emails per Conversation

**Example:**
- Client A: 5 emails, 2 conversations → Avg: 2.5 emails/conversation
- Client B: 3 emails, 1 conversation → Avg: 3.0 emails/conversation
- Client C: 12 emails, 3 conversations → Avg: 4.0 emails/conversation

**Interpretation:**
- Lower average = simpler leads (fewer back-and-forth)
- Higher average = more complex leads (more discussion needed)

## Use Cases

### Filter by Email Type
- View only STLPH emails for business analysis
- Filter out marketing/personal emails
- Focus on revenue-generating communications

### Analyze Conversation Patterns
- Identify leads that required more communication
- Find clients with complex booking processes
- Track follow-up effectiveness

### Client Communication Analysis
- See which clients send many emails
- Identify clients with multiple conversations
- Understand communication volume per client

## Customization

### Adjust Classification Keywords

**Add more STLPH keywords:**
```javascript
const STLPH_KEYWORDS = [
  'party', 'event', 'helpers', 
  'your-custom-keyword'  // Add here
];
```

**Add more non-STLPH keywords:**
```javascript
const NON_STLPH_KEYWORDS = [
  'school', 'marketing',
  'your-custom-keyword'  // Add here
];
```

### Improve Conversation Grouping

The `generateConversationId()` function can be customized to:
- Group by subject similarity
- Group by date proximity
- Group by client name (if email varies)

## Example Output

### Clients Sheet Example:
| Client Email | Total | STLPH | Other | Conversations | Avg/Conv | Status |
|--------------|-------|-------|-------|---------------|----------|--------|
| client@example.com | 8 | 6 | 2 | 2 | 4.0 | Active |
| lead@example.com | 3 | 3 | 0 | 1 | 3.0 | Lead |

### Conversations Sheet Example:
| Conversation ID | Client | Subject | Emails | STLPH | Other | Status | Outcome |
|-----------------|--------|---------|--------|-------|-------|--------|---------|
| thread123_client@... | client@... | Quote Request | 5 | 4 | 1 | STLPH Lead | Confirmed |
| thread456_lead@... | lead@... | Event Inquiry | 3 | 3 | 0 | STLPH Lead | Pending |

### Email Classification Sheet Example:
| Email Type | Count | Percentage | Avg per Client | Examples |
|------------|-------|------------|----------------|----------|
| STLPH | 1,234 | 78.5% | 12.3 | Form submissions, quotes |
| Other | 338 | 21.5% | 3.4 | Marketing, school, receipts |

## Benefits

✅ **Better Focus**: Filter to STLPH emails only for business analysis  
✅ **Conversation Tracking**: See complete lead journeys  
✅ **Communication Metrics**: Understand email volume patterns  
✅ **Client Insights**: Identify high-engagement vs low-engagement clients  
✅ **Data Quality**: Separate business data from noise  
