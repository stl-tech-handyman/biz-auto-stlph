/**
 * Gmail Email Analyzer for Business Revenue Analytics
 * 
 * This script analyzes Gmail emails to extract form submissions, track confirmations,
 * and generate revenue analytics in Google Sheets.
 * 
 * Setup:
 * 1. Create a new Google Apps Script project
 * 2. Copy this entire script
 * 3. Create a Google Sheet for output
 * 4. Update SPREADSHEET_ID below
 * 5. Run setupSheet() once to initialize
 * 6. Run analyzeAllEmails() to process historical emails
 * 7. Set up a trigger to run processNewEmails() daily/hourly
 */

// Configuration
const SPREADSHEET_ID = 'YOUR_SPREADSHEET_ID_HERE'; // Replace with your Google Sheet ID
const GMAIL_QUERY = 'from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission"'; // Adjust based on your form submission patterns
const TEST_KEYWORDS = ['test', 'testing', 'demo', 'sample', 'fake']; // Keywords to identify test submissions
const CONFIRMATION_KEYWORDS = ['confirm', 'confirmed', 'reservation', 'booking', 'deposit received', 'thank you for your booking'];
const YOUR_EMAIL = 'team@stlpartyhelpers.com'; // Your business email

// Sheet structure
const SHEETS = {
  RAW_DATA: 'Raw Data',
  CLIENTS: 'Clients',
  MONTHLY_REVENUE: 'Monthly Revenue',
  YEARLY_SUMMARY: 'Yearly Summary',
  ANALYTICS: 'Analytics Dashboard'
};

/**
 * Initialize the Google Sheet with proper structure
 */
function setupSheet() {
  const ss = SpreadsheetApp.openById(SPREADSHEET_ID);
  
  // Create sheets if they don't exist
  Object.values(SHEETS).forEach(sheetName => {
    let sheet = ss.getSheetByName(sheetName);
    if (!sheet) {
      sheet = ss.insertSheet(sheetName);
    }
  });
  
  // Setup Raw Data sheet
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  rawSheet.clear();
  rawSheet.getRange(1, 1, 1, 15).setValues([[
    'Email ID', 'Date', 'From Email', 'Subject', 'Body Preview', 
    'Is Test', 'Is Confirmation', 'Client Email', 'Event Date', 
    'Total Cost', 'Rate', 'Hours', 'Helpers', 'Occasion', 'Status'
  ]]);
  rawSheet.getRange(1, 1, 1, 15).setFontWeight('bold');
  rawSheet.setFrozenRows(1);
  
  // Setup Clients sheet
  const clientsSheet = ss.getSheetByName(SHEETS.CLIENTS);
  clientsSheet.clear();
  clientsSheet.getRange(1, 1, 1, 8).setValues([[
    'Client Email', 'First Submission', 'Last Submission', 
    'Total Submissions', 'Total Revenue', 'Confirmed Events', 
    'Unconfirmed Events', 'Status'
  ]]);
  clientsSheet.getRange(1, 1, 1, 8).setFontWeight('bold');
  clientsSheet.setFrozenRows(1);
  
  // Setup Monthly Revenue sheet
  const monthlySheet = ss.getSheetByName(SHEETS.MONTHLY_REVENUE);
  monthlySheet.clear();
  monthlySheet.getRange(1, 1, 1, 10).setValues([[
    'Year', 'Month', 'Total Revenue', 'Gross Revenue', 
    'Payout (45%)', 'Self Payout (10% of 45%)', 
    'Net Revenue', 'Confirmed Events', 'Unconfirmed Events', 'Avg Rate'
  ]]);
  monthlySheet.getRange(1, 1, 1, 10).setFontWeight('bold');
  monthlySheet.setFrozenRows(1);
  
  // Setup Yearly Summary sheet
  const yearlySheet = ss.getSheetByName(SHEETS.YEARLY_SUMMARY);
  yearlySheet.clear();
  yearlySheet.getRange(1, 1, 1, 9).setValues([[
    'Year', 'Total Revenue', 'Gross Revenue', 
    'Total Payout (45%)', 'Self Payout (10% of 45%)', 
    'Net Revenue', 'Total Events', 'Confirmed', 'Unconfirmed'
  ]]);
  yearlySheet.getRange(1, 1, 1, 9).setFontWeight('bold');
  yearlySheet.setFrozenRows(1);
  
  Logger.log('Sheet setup complete!');
}

/**
 * Analyze all emails in Gmail
 */
function analyzeAllEmails() {
  Logger.log('Starting email analysis...');
  
  const ss = SpreadsheetApp.openById(SPREADSHEET_ID);
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  
  // Search for all relevant emails
  let threads = GmailApp.search(GMAIL_QUERY, 0, 500); // Process in batches
  let processedCount = 0;
  let batch = [];
  
  threads.forEach(thread => {
    const messages = thread.getMessages();
    messages.forEach(message => {
      const emailData = extractEmailData(message);
      if (emailData) {
        batch.push(emailData);
        processedCount++;
        
        // Write in batches of 100
        if (batch.length >= 100) {
          writeBatchToSheet(rawSheet, batch);
          batch = [];
        }
      }
    });
  });
  
  // Write remaining batch
  if (batch.length > 0) {
    writeBatchToSheet(rawSheet, batch);
  }
  
  Logger.log(`Processed ${processedCount} emails`);
  
  // Process the data
  processRawData();
  
  Logger.log('Analysis complete!');
}

/**
 * Process new emails (to be run on a schedule)
 */
function processNewEmails() {
  const ss = SpreadsheetApp.openById(SPREADSHEET_ID);
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  
  // Get last processed date
  const lastRow = rawSheet.getLastRow();
  let lastDate = new Date('2020-01-01'); // Default start date
  
  if (lastRow > 1) {
    const lastDateStr = rawSheet.getRange(lastRow, 2).getValue();
    if (lastDateStr) {
      lastDate = new Date(lastDateStr);
    }
  }
  
  // Search for emails after last date
  const query = `${GMAIL_QUERY} after:${formatDateForQuery(lastDate)}`;
  const threads = GmailApp.search(query, 0, 100);
  
  let batch = [];
  threads.forEach(thread => {
    const messages = thread.getMessages();
    messages.forEach(message => {
      const messageDate = message.getDate();
      if (messageDate > lastDate) {
        const emailData = extractEmailData(message);
        if (emailData) {
          batch.push(emailData);
        }
      }
    });
  });
  
  if (batch.length > 0) {
    writeBatchToSheet(rawSheet, batch);
    processRawData();
    Logger.log(`Processed ${batch.length} new emails`);
  }
}

/**
 * Extract data from an email message
 */
function extractEmailData(message) {
  try {
    const subject = message.getSubject();
    const from = message.getFrom();
    const body = message.getPlainBody();
    const date = message.getDate();
    const messageId = message.getId();
    
    // Skip if it's from your own email (outgoing)
    if (from.includes(YOUR_EMAIL)) {
      return null;
    }
    
    // Check if it's a test submission
    const isTest = TEST_KEYWORDS.some(keyword => 
      subject.toLowerCase().includes(keyword) || 
      body.toLowerCase().includes(keyword)
    );
    
    // Check if it's a confirmation email
    const isConfirmation = CONFIRMATION_KEYWORDS.some(keyword =>
      subject.toLowerCase().includes(keyword) ||
      body.toLowerCase().includes(keyword)
    );
    
    // Extract client email (usually in form submissions)
    let clientEmail = extractClientEmail(from, body);
    
    // Extract event details
    const eventData = extractEventData(body, subject);
    
    // Extract pricing information
    const pricingData = extractPricingData(body, subject);
    
    return [
      messageId,
      date,
      from,
      subject,
      body.substring(0, 200), // Preview
      isTest ? 'Yes' : 'No',
      isConfirmation ? 'Yes' : 'No',
      clientEmail || '',
      eventData.eventDate || '',
      pricingData.totalCost || '',
      pricingData.rate || '',
      eventData.hours || '',
      eventData.helpers || '',
      eventData.occasion || '',
      isConfirmation ? 'Confirmed' : (isTest ? 'Test' : 'Pending')
    ];
  } catch (error) {
    Logger.log(`Error processing email ${message.getId()}: ${error}`);
    return null;
  }
}

/**
 * Extract client email from email content
 */
function extractClientEmail(from, body) {
  // Try to extract from common patterns
  const emailPattern = /([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)/gi;
  const matches = body.match(emailPattern);
  
  if (matches) {
    // Filter out your own email and common system emails
    const filtered = matches.filter(email => 
      !email.includes(YOUR_EMAIL) &&
      !email.includes('zapier') &&
      !email.includes('noreply') &&
      !email.includes('no-reply')
    );
    if (filtered.length > 0) {
      return filtered[0];
    }
  }
  
  // Fallback to sender email if it's not your email
  if (!from.includes(YOUR_EMAIL)) {
    const fromMatch = from.match(emailPattern);
    if (fromMatch) {
      return fromMatch[0];
    }
  }
  
  return '';
}

/**
 * Extract event data from email body
 */
function extractEventData(body, subject) {
  const data = {};
  
  // Extract event date
  const datePatterns = [
    /event date[:\s]+([A-Za-z]+\s+\d{1,2},?\s+\d{4})/i,
    /date[:\s]+([A-Za-z]+\s+\d{1,2},?\s+\d{4})/i,
    /(\d{1,2}\/\d{1,2}\/\d{4})/,
    /(\d{4}-\d{2}-\d{2})/
  ];
  
  for (const pattern of datePatterns) {
    const match = body.match(pattern) || subject.match(pattern);
    if (match) {
      data.eventDate = match[1];
      break;
    }
  }
  
  // Extract hours
  const hoursMatch = body.match(/(\d+(?:\.\d+)?)\s*hours?/i) || 
                     body.match(/for\s+(\d+(?:\.\d+)?)\s+hours?/i);
  if (hoursMatch) {
    data.hours = parseFloat(hoursMatch[1]);
  }
  
  // Extract helpers
  const helpersMatch = body.match(/(\d+)\s+helpers?/i) ||
                       body.match(/helpers?[:\s]+(\d+)/i);
  if (helpersMatch) {
    data.helpers = parseInt(helpersMatch[1]);
  }
  
  // Extract occasion
  const occasionMatch = body.match(/occasion[:\s]+([^\n]+)/i) ||
                        body.match(/for\s+(?:a\s+)?([^\n,]+?)\s+(?:party|event|celebration)/i);
  if (occasionMatch) {
    data.occasion = occasionMatch[1].trim();
  }
  
  return data;
}

/**
 * Extract pricing data from email body
 */
function extractPricingData(body, subject) {
  const data = {};
  
  // Extract total cost
  const costPatterns = [
    /\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/g,
    /total[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/i,
    /estimate[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/i
  ];
  
  for (const pattern of costPatterns) {
    const matches = body.match(pattern) || subject.match(pattern);
    if (matches) {
      // Get the largest number (usually the total)
      const amounts = matches.map(m => {
        const num = parseFloat(m.replace(/[$,]/g, ''));
        return isNaN(num) ? 0 : num;
      });
      if (amounts.length > 0) {
        data.totalCost = Math.max(...amounts);
      }
      break;
    }
  }
  
  // Extract rate (hourly rate)
  const rateMatch = body.match(/(?:rate|hourly)[:\s]+\$?(\d+(?:\.\d+)?)/i);
  if (rateMatch) {
    data.rate = parseFloat(rateMatch[1]);
  }
  
  return data;
}

/**
 * Write a batch of data to the sheet
 */
function writeBatchToSheet(sheet, batch) {
  if (batch.length === 0) return;
  
  const lastRow = sheet.getLastRow();
  sheet.getRange(lastRow + 1, 1, batch.length, batch[0].length).setValues(batch);
}

/**
 * Process raw data and create summary sheets
 */
function processRawData() {
  const ss = SpreadsheetApp.openById(SPREADSHEET_ID);
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  const clientsSheet = ss.getSheetByName(SHEETS.CLIENTS);
  const monthlySheet = ss.getSheetByName(SHEETS.MONTHLY_REVENUE);
  const yearlySheet = ss.getSheetByName(SHEETS.YEARLY_SUMMARY);
  
  const data = rawSheet.getDataRange().getValues();
  if (data.length <= 1) return; // No data
  
  // Skip header row
  const rows = data.slice(1);
  
  // Filter out test submissions
  const realSubmissions = rows.filter(row => row[5] === 'No'); // Column 6 (0-indexed 5) is "Is Test"
  
  // Group by client email
  const clientsMap = new Map();
  
  realSubmissions.forEach(row => {
    const clientEmail = row[7] || ''; // Column 8
    if (!clientEmail) return;
    
    if (!clientsMap.has(clientEmail)) {
      clientsMap.set(clientEmail, {
        email: clientEmail,
        firstDate: new Date(row[1]),
        lastDate: new Date(row[1]),
        submissions: [],
        confirmed: 0,
        unconfirmed: 0,
        totalRevenue: 0
      });
    }
    
    const client = clientsMap.get(clientEmail);
    const date = new Date(row[1]);
    if (date < client.firstDate) client.firstDate = date;
    if (date > client.lastDate) client.lastDate = date;
    
    client.submissions.push(row);
    
    const cost = parseFloat(row[9]) || 0; // Column 10
    client.totalRevenue += cost;
    
    if (row[6] === 'Yes') { // Column 7 is "Is Confirmation"
      client.confirmed++;
    } else {
      client.unconfirmed++;
    }
  });
  
  // Write clients data
  clientsSheet.getRange(2, 1, clientsSheet.getLastRow(), 8).clear();
  const clientsData = Array.from(clientsMap.values()).map(client => [
    client.email,
    client.firstDate,
    client.lastDate,
    client.submissions.length,
    client.totalRevenue,
    client.confirmed,
    client.unconfirmed,
    client.confirmed > 0 ? 'Active' : 'Lead'
  ]);
  
  if (clientsData.length > 0) {
    clientsSheet.getRange(2, 1, clientsData.length, 8).setValues(clientsData);
  }
  
  // Group by year and month
  const monthlyMap = new Map();
  
  realSubmissions.forEach(row => {
    const date = new Date(row[1]);
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const key = `${year}-${month}`;
    
    if (!monthlyMap.has(key)) {
      monthlyMap.set(key, {
        year: year,
        month: month,
        revenue: 0,
        confirmed: 0,
        unconfirmed: 0,
        rates: []
      });
    }
    
    const monthData = monthlyMap.get(key);
    const cost = parseFloat(row[9]) || 0;
    monthData.revenue += cost;
    
    if (row[6] === 'Yes') {
      monthData.confirmed++;
    } else {
      monthData.unconfirmed++;
    }
    
    const rate = parseFloat(row[10]) || 0;
    if (rate > 0) {
      monthData.rates.push(rate);
    }
  });
  
  // Write monthly data
  monthlySheet.getRange(2, 1, monthlySheet.getLastRow(), 10).clear();
  const monthlyData = Array.from(monthlyMap.values())
    .sort((a, b) => {
      if (a.year !== b.year) return a.year - b.year;
      return a.month - b.month;
    })
    .map(month => {
      const grossRevenue = month.revenue;
      const payout = grossRevenue * 0.45;
      const selfPayout = payout * 0.10;
      const netRevenue = grossRevenue - payout;
      const avgRate = month.rates.length > 0 
        ? month.rates.reduce((a, b) => a + b, 0) / month.rates.length 
        : 0;
      
      return [
        month.year,
        month.month,
        grossRevenue,
        grossRevenue,
        payout,
        selfPayout,
        netRevenue,
        month.confirmed,
        month.unconfirmed,
        avgRate
      ];
    });
  
  if (monthlyData.length > 0) {
    monthlySheet.getRange(2, 1, monthlyData.length, 10).setValues(monthlyData);
  }
  
  // Group by year
  const yearlyMap = new Map();
  
  realSubmissions.forEach(row => {
    const date = new Date(row[1]);
    const year = date.getFullYear();
    
    if (!yearlyMap.has(year)) {
      yearlyMap.set(year, {
        year: year,
        revenue: 0,
        confirmed: 0,
        unconfirmed: 0
      });
    }
    
    const yearData = yearlyMap.get(year);
    const cost = parseFloat(row[9]) || 0;
    yearData.revenue += cost;
    
    if (row[6] === 'Yes') {
      yearData.confirmed++;
    } else {
      yearData.unconfirmed++;
    }
  });
  
  // Write yearly data
  yearlySheet.getRange(2, 1, yearlySheet.getLastRow(), 9).clear();
  const yearlyData = Array.from(yearlyMap.values())
    .sort((a, b) => a.year - b.year)
    .map(year => {
      const grossRevenue = year.revenue;
      const payout = grossRevenue * 0.45;
      const selfPayout = payout * 0.10;
      const netRevenue = grossRevenue - payout;
      
      return [
        year.year,
        grossRevenue,
        grossRevenue,
        payout,
        selfPayout,
        netRevenue,
        year.confirmed + year.unconfirmed,
        year.confirmed,
        year.unconfirmed
      ];
    });
  
  if (yearlyData.length > 0) {
    yearlySheet.getRange(2, 1, yearlyData.length, 9).setValues(yearlyData);
  }
  
  // Format sheets
  formatSheets(ss);
  
  Logger.log('Data processing complete!');
}

/**
 * Format the sheets for better readability
 */
function formatSheets(ss) {
  // Format monthly revenue sheet
  const monthlySheet = ss.getSheetByName(SHEETS.MONTHLY_REVENUE);
  const monthlyRange = monthlySheet.getRange(2, 3, monthlySheet.getLastRow() - 1, 7);
  monthlyRange.setNumberFormat('$#,##0.00');
  
  // Format yearly summary sheet
  const yearlySheet = ss.getSheetByName(SHEETS.YEARLY_SUMMARY);
  const yearlyRange = yearlySheet.getRange(2, 2, yearlySheet.getLastRow() - 1, 5);
  yearlyRange.setNumberFormat('$#,##0.00');
  
  // Auto-resize columns
  [monthlySheet, yearlySheet, ss.getSheetByName(SHEETS.CLIENTS)].forEach(sheet => {
    sheet.autoResizeColumns(1, sheet.getLastColumn());
  });
}

/**
 * Helper function to format date for Gmail query
 */
function formatDateForQuery(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}/${month}/${day}`;
}

/**
 * Set up time-driven trigger for processing new emails
 * Run this once to set up automation
 */
function setupTrigger() {
  // Delete existing triggers
  const triggers = ScriptApp.getProjectTriggers();
  triggers.forEach(trigger => {
    if (trigger.getHandlerFunction() === 'processNewEmails') {
      ScriptApp.deleteTrigger(trigger);
    }
  });
  
  // Create new trigger (runs daily at 2 AM)
  ScriptApp.newTrigger('processNewEmails')
    .timeBased()
    .everyDays(1)
    .atHour(2)
    .create();
  
  Logger.log('Trigger set up successfully!');
}
