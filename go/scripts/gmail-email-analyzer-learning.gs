/**
 * Learning Gmail Email Analyzer
 * 
 * This version learns patterns as it processes emails, allowing you to:
 * 1. Process a sample batch first to discover patterns
 * 2. Review and refine extraction rules
 * 3. Process all emails with learned patterns
 * 4. Continuously improve as more emails are processed
 */

// Configuration
const FOLDER_NAME = 'BOS'; // Business Operating System folder
const SPREADSHEET_NAME = 'Email Revenue Analytics'; // Name of the Google Sheet
const GMAIL_QUERY = 'from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission" OR subject:"Quote"';
const YOUR_EMAIL = 'team@stlpartyhelpers.com';

// STLPH-related keywords for classification
const STLPH_KEYWORDS = [
  'party', 'event', 'helpers', 'stl party', 'stlpartyhelpers', 'booking', 
  'quote', 'deposit', 'reservation', 'staffing', 'celebration', 'birthday',
  'wedding', 'corporate', 'venue', 'guests', 'hours', 'helpers requested',
  'event date', 'occasion', 'total cost', 'estimate', 'invoice'
];

// Non-STLPH keywords (marketing, personal, etc.)
const NON_STLPH_KEYWORDS = [
  'school', 'pta', 'parent', 'newsletter', 'promotion', 'sale', 'discount',
  'unsubscribe', 'marketing', 'advertisement', 'spam', 'receipt', 'order confirmation',
  'shipping', 'delivery', 'account', 'password', 'verification'
];

// Sheet tabs (pages) within the single Google Sheet
const SHEETS = {
  RAW_DATA: 'Raw Data',
  CLIENTS: 'Clients',
  MONTHLY_REVENUE: 'Monthly Revenue',
  YEARLY_SUMMARY: 'Yearly Summary',
  PATTERN_DISCOVERY: 'Pattern Discovery',
  CLIENT_MATCHING: 'Client Matching',
  PROCESSING_LOG: 'Processing Log',
  SAMPLE_DATA: 'Sample Analysis',
  CONVERSATIONS: 'Conversations',
  EMAIL_CLASSIFICATION: 'Email Classification'
};

// Global variables
let SPREADSHEET_ID = null;
let BOS_FOLDER_ID = null;

// Learned patterns (will be populated as we process)
const LEARNED_PATTERNS = {
  testKeywords: [],
  confirmationKeywords: [],
  clientEmailPatterns: [],
  eventDatePatterns: [],
  pricingPatterns: [],
  ratePatterns: []
};

/**
 * Initialize BOS folder and spreadsheet
 * Run this first to set up the folder and sheet
 */
function initializeBOS() {
  Logger.log('Initializing BOS folder and spreadsheet...');
  
  // Get or create BOS folder
  BOS_FOLDER_ID = getOrCreateBOSFolder();
  Logger.log(`BOS Folder ID: ${BOS_FOLDER_ID}`);
  
  // Get or create spreadsheet
  const ss = getOrCreateSpreadsheet();
  SPREADSHEET_ID = ss.getId();
  Logger.log(`Spreadsheet ID: ${SPREADSHEET_ID}`);
  
  // Store in script properties for future use
  const properties = PropertiesService.getScriptProperties();
  properties.setProperty('BOS_FOLDER_ID', BOS_FOLDER_ID);
  properties.setProperty('SPREADSHEET_ID', SPREADSHEET_ID);
  
  // Initialize all sheet tabs
  setupAllSheets(ss);
  
  Logger.log('BOS initialization complete!');
  Logger.log(`Spreadsheet URL: ${ss.getUrl()}`);
  
  return ss;
}

/**
 * Get or create BOS folder in Google Drive
 */
function getOrCreateBOSFolder() {
  const properties = PropertiesService.getScriptProperties();
  let folderId = properties.getProperty('BOS_FOLDER_ID');
  
  if (folderId) {
    try {
      const folder = DriveApp.getFolderById(folderId);
      // Verify folder still exists
      folder.getName(); // This will throw if folder doesn't exist
      return folderId;
    } catch (e) {
      // Folder doesn't exist, create new one
    }
  }
  
  // Search for existing BOS folder
  const folders = DriveApp.getFoldersByName(FOLDER_NAME);
  if (folders.hasNext()) {
    const folder = folders.next();
    folderId = folder.getId();
    properties.setProperty('BOS_FOLDER_ID', folderId);
    return folderId;
  }
  
  // Create new BOS folder
  const folder = DriveApp.createFolder(FOLDER_NAME);
  folderId = folder.getId();
  properties.setProperty('BOS_FOLDER_ID', folderId);
  Logger.log(`Created new BOS folder: ${folder.getName()}`);
  
  return folderId;
}

/**
 * Get or create the main spreadsheet
 */
function getOrCreateSpreadsheet() {
  const properties = PropertiesService.getScriptProperties();
  let spreadsheetId = properties.getProperty('SPREADSHEET_ID');
  
  if (spreadsheetId) {
    try {
      const ss = SpreadsheetApp.openById(spreadsheetId);
      // Verify spreadsheet still exists
      ss.getName(); // This will throw if spreadsheet doesn't exist
      return ss;
    } catch (e) {
      // Spreadsheet doesn't exist, create new one
    }
  }
  
  // Get BOS folder
  if (!BOS_FOLDER_ID) {
    BOS_FOLDER_ID = getOrCreateBOSFolder();
  }
  const folder = DriveApp.getFolderById(BOS_FOLDER_ID);
  
  // Search for existing spreadsheet in BOS folder
  const files = folder.getFilesByName(SPREADSHEET_NAME);
  if (files.hasNext()) {
    const file = files.next();
    const ss = SpreadsheetApp.openById(file.getId());
    spreadsheetId = ss.getId();
    properties.setProperty('SPREADSHEET_ID', spreadsheetId);
    return ss;
  }
  
  // Create new spreadsheet in BOS folder
  const ss = SpreadsheetApp.create(SPREADSHEET_NAME);
  const file = DriveApp.getFileById(ss.getId());
  folder.addFile(file);
  DriveApp.getRootFolder().removeFile(file); // Remove from root, keep only in BOS folder
  
  spreadsheetId = ss.getId();
  properties.setProperty('SPREADSHEET_ID', spreadsheetId);
  Logger.log(`Created new spreadsheet: ${ss.getName()}`);
  
  return ss;
}

/**
 * Get the spreadsheet (loads from properties if available)
 */
function getSpreadsheet() {
  // If already loaded, use it
  if (SPREADSHEET_ID) {
    try {
      return SpreadsheetApp.openById(SPREADSHEET_ID);
    } catch (e) {
      // Spreadsheet doesn't exist, reset and initialize
      SPREADSHEET_ID = null;
    }
  }
  
  // Try to load from properties
  const properties = PropertiesService.getScriptProperties();
  const spreadsheetId = properties.getProperty('SPREADSHEET_ID');
  
  if (spreadsheetId) {
    try {
      const ss = SpreadsheetApp.openById(spreadsheetId);
      SPREADSHEET_ID = spreadsheetId;
      return ss;
    } catch (e) {
      // Spreadsheet doesn't exist, clear property and initialize
      properties.deleteProperty('SPREADSHEET_ID');
    }
  }
  
  // If no spreadsheet exists, initialize
  Logger.log('No spreadsheet found. Initializing BOS...');
  return initializeBOS();
}

/**
 * Setup all sheet tabs in the spreadsheet
 */
function setupAllSheets(ss) {
  // Create all sheet tabs
  Object.values(SHEETS).forEach(sheetName => {
    let sheet = ss.getSheetByName(sheetName);
    if (!sheet) {
      sheet = ss.insertSheet(sheetName);
    }
  });
  
  // Remove default "Sheet1" if it exists and is empty
  const defaultSheet = ss.getSheetByName('Sheet1');
  if (defaultSheet && defaultSheet.getLastRow() === 0 && ss.getSheets().length > 1) {
    ss.deleteSheet(defaultSheet);
  }
  
  // Setup each sheet's structure
  setupRawDataSheet(ss.getSheetByName(SHEETS.RAW_DATA));
  setupClientsSheet(ss.getSheetByName(SHEETS.CLIENTS));
  setupMonthlyRevenueSheet(ss.getSheetByName(SHEETS.MONTHLY_REVENUE));
  setupYearlySummarySheet(ss.getSheetByName(SHEETS.YEARLY_SUMMARY));
  setupPatternDiscoverySheet(ss.getSheetByName(SHEETS.PATTERN_DISCOVERY));
  setupClientMatchingSheet(ss.getSheetByName(SHEETS.CLIENT_MATCHING));
  setupLogSheet(ss.getSheetByName(SHEETS.PROCESSING_LOG));
  setupSampleDataSheet(ss.getSheetByName(SHEETS.SAMPLE_DATA));
  setupConversationsSheet(ss.getSheetByName(SHEETS.CONVERSATIONS));
  setupEmailClassificationSheet(ss.getSheetByName(SHEETS.EMAIL_CLASSIFICATION));
  
  Logger.log('All sheet tabs initialized');
}

/**
 * Setup Raw Data sheet tab
 */
function setupRawDataSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 21).setValues([[
      'Email ID', 'Thread ID', 'Date', 'From Email', 'Subject', 'Body Preview', 
      'Is Test', 'Is Confirmation', 'Client Email', 'Event Date', 
      'Total Cost', 'Rate', 'Hours', 'Helpers', 'Occasion', 'Status',
      'Guests', 'Deposit', 'Email Type', 'Conversation ID', 'Message Number'
    ]]);
    sheet.getRange(1, 1, 1, 21).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Clients sheet tab
 */
function setupClientsSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 12).setValues([[
      'Client Email', 'First Submission', 'Last Submission', 
      'Total Submissions', 'STLPH Emails', 'Other Emails',
      'Total Revenue', 'Confirmed Events', 'Unconfirmed Events', 
      'Conversations', 'Avg Emails per Conversation', 'Status'
    ]]);
    sheet.getRange(1, 1, 1, 12).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Monthly Revenue sheet tab
 */
function setupMonthlyRevenueSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 10).setValues([[
      'Year', 'Month', 'Total Revenue', 'Gross Revenue', 
      'Payout (45%)', 'Self Payout (10% of 45%)', 
      'Net Revenue', 'Confirmed Events', 'Unconfirmed Events', 'Avg Rate'
    ]]);
    sheet.getRange(1, 1, 1, 10).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Yearly Summary sheet tab
 */
function setupYearlySummarySheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 9).setValues([[
      'Year', 'Total Revenue', 'Gross Revenue', 
      'Total Payout (45%)', 'Self Payout (10% of 45%)', 
      'Net Revenue', 'Total Events', 'Confirmed', 'Unconfirmed'
    ]]);
    sheet.getRange(1, 1, 1, 9).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Pattern Discovery sheet tab
 */
function setupPatternDiscoverySheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 4).setValues([[
      'Pattern Type', 'Pattern Found', 'Frequency', 'Example'
    ]]);
    sheet.getRange(1, 1, 1, 4).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Client Matching sheet tab
 */
function setupClientMatchingSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 6).setValues([[
      'Client Email', 'Alternative Emails', 'Submission Count', 
      'First Seen', 'Last Seen', 'Confidence'
    ]]);
    sheet.getRange(1, 1, 1, 6).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Processing Log sheet tab
 */
function setupLogSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 4).setValues([[
      'Timestamp', 'Status', 'Emails Processed', 'Notes'
    ]]);
    sheet.getRange(1, 1, 1, 4).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Sample Data sheet tab
 */
function setupSampleDataSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 12).setValues([[
      'Date', 'From', 'Subject', 'Body Preview', 'Client Email?', 
      'Event Date?', 'Total Cost?', 'Rate?', 'Is Test?', 'Is Confirmation?',
      'Email Type?', 'Conversation?'
    ]]);
    sheet.getRange(1, 1, 1, 12).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Conversations sheet tab
 */
function setupConversationsSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 10).setValues([[
      'Conversation ID', 'Client Email', 'Subject', 'First Email', 'Last Email',
      'Email Count', 'STLPH Emails', 'Other Emails', 'Status', 'Outcome'
    ]]);
    sheet.getRange(1, 1, 1, 10).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Setup Email Classification sheet tab
 */
function setupEmailClassificationSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 6).setValues([[
      'Email Type', 'Count', 'Percentage', 'Avg per Client', 'Examples', 'Notes'
    ]]);
    sheet.getRange(1, 1, 1, 6).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

/**
 * Step 1: Analyze a sample batch to discover patterns
 */
function discoverPatterns(sampleSize = 100) {
  Logger.log(`Analyzing ${sampleSize} emails to discover patterns...`);
  
  const ss = getSpreadsheet();
  const sampleSheet = ss.getSheetByName(SHEETS.SAMPLE_DATA);
  const patternSheet = ss.getSheetByName(SHEETS.PATTERN_DISCOVERY);
  
  // Setup sample sheet
  sampleSheet.clear();
  sampleSheet.getRange(1, 1, 1, 10).setValues([[
    'Date', 'From', 'Subject', 'Body Preview', 'Client Email?', 
    'Event Date?', 'Total Cost?', 'Rate?', 'Is Test?', 'Is Confirmation?'
  ]]);
  sampleSheet.getRange(1, 1, 1, 10).setFontWeight('bold');
  
  // Setup pattern discovery sheet
  patternSheet.clear();
  patternSheet.getRange(1, 1, 1, 4).setValues([[
    'Pattern Type', 'Pattern Found', 'Frequency', 'Example'
  ]]);
  patternSheet.getRange(1, 1, 1, 4).setFontWeight('bold');
  
  const threads = GmailApp.search(GMAIL_QUERY, 0, sampleSize);
  const samples = [];
  const patternCounts = {
    testKeywords: {},
    confirmationKeywords: {},
    emailPatterns: {},
    datePatterns: {},
    costPatterns: {},
    ratePatterns: {}
  };
  
  threads.forEach(thread => {
    const messages = thread.getMessages();
    messages.forEach(message => {
      const sample = analyzeSampleEmail(message, patternCounts);
      if (sample) {
        samples.push(sample.rowData);
      }
    });
  });
  
  // Write samples
  if (samples.length > 0) {
    sampleSheet.getRange(2, 1, samples.length, 10).setValues(samples);
  }
  
  // Analyze and document discovered patterns
  documentDiscoveredPatterns(patternSheet, patternCounts);
  
  // Store learned patterns
  storeLearnedPatterns(patternCounts);
  
  Logger.log(`Pattern discovery complete! Review ${SHEETS.SAMPLE_DATA} and ${SHEETS.PATTERN_DISCOVERY} sheets.`);
  Logger.log('After reviewing, run refinePatterns() to update extraction rules, then run analyzeAllEmailsLearning()');
}

/**
 * Analyze a single email sample and collect patterns
 */
function analyzeSampleEmail(message, patternCounts) {
  try {
    const subject = message.getSubject();
    const from = message.getFrom();
    const body = message.getPlainBody();
    const htmlBody = message.getBody();
    const date = message.getDate();
    
    // Skip outgoing
    if (from.includes(YOUR_EMAIL) && !subject.toLowerCase().includes('reply')) {
      return null;
    }
    
    const combined = (subject + ' ' + body).toLowerCase();
    
    // Discover test keywords
    const testWords = ['test', 'testing', 'demo', 'sample', 'fake', 'dummy'];
    testWords.forEach(word => {
      if (combined.includes(word)) {
        patternCounts.testKeywords[word] = (patternCounts.testKeywords[word] || 0) + 1;
      }
    });
    
    // Discover confirmation keywords
    const confWords = ['confirm', 'confirmed', 'reservation', 'booking', 'deposit', 'thank you'];
    confWords.forEach(word => {
      if (combined.includes(word)) {
        patternCounts.confirmationKeywords[word] = (patternCounts.confirmationKeywords[word] || 0) + 1;
      }
    });
    
    // Extract and track email patterns
    const emailPattern = /([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)/gi;
    const emails = body.match(emailPattern) || [];
    emails.forEach(email => {
      if (!email.includes(YOUR_EMAIL) && !email.includes('zapier') && !email.includes('noreply')) {
        const domain = email.split('@')[1];
        patternCounts.emailPatterns[domain] = (patternCounts.emailPatterns[domain] || 0) + 1;
      }
    });
    
    // Extract and track date patterns
    const datePatterns = [
      /(\d{1,2}\/\d{1,2}\/\d{4})/g,
      /(\d{4}-\d{2}-\d{2})/g,
      /([A-Za-z]+\s+\d{1,2},?\s+\d{4})/g
    ];
    datePatterns.forEach(pattern => {
      const matches = body.match(pattern) || subject.match(pattern) || [];
      matches.forEach(match => {
        patternCounts.datePatterns[match] = (patternCounts.datePatterns[match] || 0) + 1;
      });
    });
    
    // Extract and track cost patterns
    const costMatches = body.match(/\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/g) || [];
    costMatches.forEach(match => {
      const amount = match.replace(/[$,]/g, '');
      if (parseFloat(amount) > 10) {
        patternCounts.costPatterns[match] = (patternCounts.costPatterns[match] || 0) + 1;
      }
    });
    
    // Extract and track rate patterns
    const rateMatches = body.match(/(?:rate|hourly|per\s+hour)[:\s]+\$?(\d+(?:\.\d+)?)/gi) || [];
    rateMatches.forEach(match => {
      patternCounts.ratePatterns[match] = (patternCounts.ratePatterns[match] || 0) + 1;
    });
    
    // Try to extract data
    const clientEmail = extractClientEmailLearning(from, body, htmlBody);
    const eventDate = extractEventDateLearning(body, subject);
    const pricingData = extractPricingDataLearning(body, subject);
    const isTest = detectTestLearning(combined);
    const isConfirmation = detectConfirmationLearning(combined);
    const emailType = classifyEmailType(subject, body, from);
    
    // Get conversation info
    const thread = message.getThread();
    const threadMessages = thread.getMessages();
    const messageIndex = threadMessages.findIndex(m => m.getId() === message.getId());
    const conversationId = generateConversationId(thread, clientEmail);
    
    return {
      rowData: [
        date,
        from,
        subject,
        body.substring(0, 200),
        clientEmail || 'NOT FOUND',
        eventDate || 'NOT FOUND',
        pricingData.totalCost || 'NOT FOUND',
        pricingData.rate || 'NOT FOUND',
        isTest ? 'Yes' : 'No',
        isConfirmation ? 'Yes' : 'No',
        emailType || 'NOT FOUND',
        conversationId || 'NOT FOUND'
      ]
    };
  } catch (error) {
    Logger.log(`Error analyzing sample: ${error}`);
    return null;
  }
}

/**
 * Document discovered patterns
 */
function documentDiscoveredPatterns(sheet, patternCounts) {
  let row = 2;
  
  // Test keywords
  Object.entries(patternCounts.testKeywords)
    .sort((a, b) => b[1] - a[1])
    .forEach(([word, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Test Keyword', word, count, '']]);
      row++;
    });
  
  // Confirmation keywords
  Object.entries(patternCounts.confirmationKeywords)
    .sort((a, b) => b[1] - a[1])
    .forEach(([word, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Confirmation Keyword', word, count, '']]);
      row++;
    });
  
  // Email domains (top 20)
  Object.entries(patternCounts.emailPatterns)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 20)
    .forEach(([domain, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Email Domain', domain, count, '']]);
      row++;
    });
  
  // Date patterns (top 10)
  Object.entries(patternCounts.datePatterns)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .forEach(([pattern, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Date Pattern', pattern, count, '']]);
      row++;
    });
  
  // Cost patterns (top 10)
  Object.entries(patternCounts.costPatterns)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .forEach(([pattern, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Cost Pattern', pattern, count, '']]);
      row++;
    });
  
  // Rate patterns (top 10)
  Object.entries(patternCounts.ratePatterns)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .forEach(([pattern, count]) => {
      sheet.getRange(row, 1, 1, 4).setValues([['Rate Pattern', pattern, count, '']]);
      row++;
    });
}

/**
 * Store learned patterns for use in full processing
 */
function storeLearnedPatterns(patternCounts) {
  // Store in script properties (persists across runs)
  const properties = PropertiesService.getScriptProperties();
  
  // Store top patterns
  const topTestKeywords = Object.entries(patternCounts.testKeywords)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .map(([word]) => word);
  
  const topConfKeywords = Object.entries(patternCounts.confirmationKeywords)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 10)
    .map(([word]) => word);
  
  properties.setProperty('learned_test_keywords', JSON.stringify(topTestKeywords));
  properties.setProperty('learned_confirmation_keywords', JSON.stringify(topConfKeywords));
  
  Logger.log('Learned patterns stored. Top test keywords:', topTestKeywords);
  Logger.log('Top confirmation keywords:', topConfKeywords);
}

/**
 * Load learned patterns
 */
function loadLearnedPatterns() {
  const properties = PropertiesService.getScriptProperties();
  
  const testKeywords = JSON.parse(properties.getProperty('learned_test_keywords') || '[]');
  const confKeywords = JSON.parse(properties.getProperty('learned_confirmation_keywords') || '[]');
  
  LEARNED_PATTERNS.testKeywords = testKeywords.length > 0 ? testKeywords : ['test', 'testing', 'demo'];
  LEARNED_PATTERNS.confirmationKeywords = confKeywords.length > 0 ? confKeywords : ['confirm', 'confirmed', 'booking'];
  
  Logger.log('Loaded learned patterns:', LEARNED_PATTERNS);
}

/**
 * Step 2: Refine patterns based on sample analysis
 * Run this after reviewing the Pattern Discovery sheet
 */
function refinePatterns() {
  const ss = getSpreadsheet();
  const patternSheet = ss.getSheetByName(SHEETS.PATTERN_DISCOVERY);
  
  if (!patternSheet || patternSheet.getLastRow() <= 1) {
    Logger.log('No patterns found. Run discoverPatterns() first.');
    return;
  }
  
  // You can manually edit the Pattern Discovery sheet, then this function
  // will reload the patterns. Or you can call updatePatternsManually()
  
  loadLearnedPatterns();
  Logger.log('Patterns refined. Ready to process all emails.');
}

/**
 * Manually update patterns (call this after editing Pattern Discovery sheet)
 */
function updatePatternsManually() {
  const ss = getSpreadsheet();
  const patternSheet = ss.getSheetByName(SHEETS.PATTERN_DISCOVERY);
  const properties = PropertiesService.getScriptProperties();
  
  const data = patternSheet.getDataRange().getValues();
  const testKeywords = [];
  const confKeywords = [];
  
  data.slice(1).forEach(row => {
    const type = row[0];
    const pattern = row[1];
    if (type === 'Test Keyword' && pattern) {
      testKeywords.push(pattern);
    } else if (type === 'Confirmation Keyword' && pattern) {
      confKeywords.push(pattern);
    }
  });
  
  properties.setProperty('learned_test_keywords', JSON.stringify(testKeywords));
  properties.setProperty('learned_confirmation_keywords', JSON.stringify(confKeywords));
  
  Logger.log('Patterns updated manually:', { testKeywords, confKeywords });
}

/**
 * Step 3: Process all emails using learned patterns
 */
function analyzeAllEmailsLearning() {
  loadLearnedPatterns();
  
  Logger.log('Starting full email analysis with learned patterns...');
  
  const ss = getSpreadsheet();
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  const logSheet = ss.getSheetByName(SHEETS.PROCESSING_LOG);
  
  // Ensure sheets are set up
  setupRawDataSheet(rawSheet);
  setupLogSheet(logSheet);
  
  logMessage(logSheet, 'Starting full analysis...', 0);
  
  let processedCount = 0;
  let batch = [];
  let lastProcessedId = getLastProcessedId(rawSheet);
  
  // Process in batches
  for (let start = 0; start < 50000; start += 500) {
    try {
      const threads = GmailApp.search(GMAIL_QUERY, start, 500);
      if (threads.length === 0) break;
      
      threads.forEach(thread => {
        const messages = thread.getMessages();
        messages.forEach(message => {
          const messageId = message.getId();
          
          if (isAlreadyProcessed(rawSheet, messageId)) {
            return;
          }
          
          const emailData = extractEmailDataLearning(message);
          if (emailData && !emailData.isTest) {
            batch.push(emailData.rowData);
            processedCount++;
            
            if (batch.length >= 50) {
              writeBatchToSheet(rawSheet, batch);
              batch = [];
              Utilities.sleep(100);
            }
          }
        });
      });
      
      if (batch.length > 0) {
        writeBatchToSheet(rawSheet, batch);
        batch = [];
      }
      
      logMessage(logSheet, `Processed batch ${start}-${start + threads.length}`, processedCount);
      Utilities.sleep(500);
      
    } catch (error) {
      logMessage(logSheet, `Error in batch ${start}: ${error.toString()}`, processedCount);
      break;
    }
  }
  
  logMessage(logSheet, `Analysis complete. Total: ${processedCount}`, processedCount);
  
  // Identify unique clients
  identifyUniqueClients();
  
  // Process aggregated data
  processRawDataLearning();
  
  Logger.log('Full analysis complete!');
}

/**
 * Extract email data using learned patterns
 */
function extractEmailDataLearning(message) {
  try {
    const subject = message.getSubject();
    const from = message.getFrom();
    const body = message.getPlainBody();
    const htmlBody = message.getBody();
    const date = message.getDate();
    const messageId = message.getId();
    const threadId = message.getThread().getId();
    
    if (from.includes(YOUR_EMAIL) && !subject.toLowerCase().includes('reply')) {
      return null;
    }
    
    const combined = (subject + ' ' + body).toLowerCase();
    
    const isTest = detectTestLearning(combined);
    const isConfirmation = detectConfirmationLearning(combined);
    const clientEmail = extractClientEmailLearning(from, body, htmlBody);
    const eventData = extractEventDataLearning(body, subject);
    const pricingData = extractPricingDataLearning(body, subject);
    
    // Classify email type
    const emailType = classifyEmailType(subject, body, from);
    
    // Get conversation info
    const thread = message.getThread();
    const threadMessages = thread.getMessages();
    const messageIndex = threadMessages.findIndex(m => m.getId() === messageId);
    const conversationId = generateConversationId(thread, clientEmail);
    
    return {
      isTest: isTest,
      emailType: emailType,
      conversationId: conversationId,
      messageNumber: messageIndex + 1,
      rowData: [
        messageId,
        threadId,
        date,
        from,
        subject,
        body.substring(0, 300),
        isTest ? 'Yes' : 'No',
        isConfirmation ? 'Yes' : 'No',
        clientEmail || '',
        eventData.eventDate || '',
        pricingData.totalCost || '',
        pricingData.rate || '',
        eventData.hours || '',
        eventData.helpers || '',
        eventData.occasion || '',
        isConfirmation ? 'Confirmed' : (isTest ? 'Test' : 'Pending'),
        eventData.guests || '',
        pricingData.deposit || '',
        emailType,
        conversationId,
        messageIndex + 1
      ]
    };
  } catch (error) {
    Logger.log(`Error extracting data: ${error}`);
    return null;
  }
}

/**
 * Classify email as STLPH-related or other
 */
function classifyEmailType(subject, body, from) {
  const combined = (subject + ' ' + body).toLowerCase();
  const fromLower = from.toLowerCase();
  
  // Check for STLPH keywords
  const stlphScore = STLPH_KEYWORDS.filter(keyword => 
    combined.includes(keyword.toLowerCase())
  ).length;
  
  // Check for non-STLPH keywords
  const nonStlphScore = NON_STLPH_KEYWORDS.filter(keyword =>
    combined.includes(keyword.toLowerCase())
  ).length;
  
  // Check if from your business domain
  if (fromLower.includes('stlpartyhelpers') || fromLower.includes('bizops360')) {
    return 'STLPH';
  }
  
  // Check if from known form sources
  if (fromLower.includes('zapier') || fromLower.includes('forms')) {
    return 'STLPH';
  }
  
  // Score-based classification
  if (stlphScore > nonStlphScore && stlphScore > 0) {
    return 'STLPH';
  } else if (nonStlphScore > stlphScore && nonStlphScore > 0) {
    return 'Other';
  } else if (stlphScore > 0) {
    return 'STLPH';
  } else {
    return 'Other';
  }
}

/**
 * Generate conversation ID for grouping related emails
 */
function generateConversationId(thread, clientEmail) {
  // Use thread ID + client email to create unique conversation ID
  // Same thread with same client = same conversation
  const threadId = thread.getId();
  const normalizedEmail = (clientEmail || '').toLowerCase().trim();
  
  if (normalizedEmail) {
    return `${threadId}_${normalizedEmail}`;
  }
  
  // Fallback: use thread ID and subject
  const subject = thread.getFirstMessageSubject();
  const subjectHash = subject.toLowerCase().replace(/[^a-z0-9]/g, '').substring(0, 20);
  return `${threadId}_${subjectHash}`;
}

/**
 * Detect test using learned patterns
 */
function detectTestLearning(text) {
  return LEARNED_PATTERNS.testKeywords.some(keyword => text.includes(keyword));
}

/**
 * Detect confirmation using learned patterns
 */
function detectConfirmationLearning(text) {
  return LEARNED_PATTERNS.confirmationKeywords.some(keyword => text.includes(keyword));
}

/**
 * Extract client email with multiple strategies
 */
function extractClientEmailLearning(from, body, htmlBody) {
  const emailPattern = /([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)/gi;
  
  // Try HTML first
  if (htmlBody) {
    const htmlMatches = htmlBody.match(emailPattern);
    if (htmlMatches) {
      const filtered = htmlMatches.filter(email => 
        !email.includes(YOUR_EMAIL) &&
        !email.includes('zapier') &&
        !email.includes('noreply') &&
        !email.includes('no-reply') &&
        !email.includes('google') &&
        !email.includes('gmail.com')
      );
      if (filtered.length > 0) return filtered[0];
    }
  }
  
  // Try body
  const bodyMatches = body.match(emailPattern);
  if (bodyMatches) {
    const filtered = bodyMatches.filter(email => 
      !email.includes(YOUR_EMAIL) &&
      !email.includes('zapier') &&
      !email.includes('noreply')
    );
    if (filtered.length > 0) return filtered[0];
  }
  
  // Fallback to sender
  if (!from.includes(YOUR_EMAIL)) {
    const fromMatch = from.match(emailPattern);
    if (fromMatch) return fromMatch[0];
  }
  
  return '';
}

/**
 * Extract event date with multiple patterns
 */
function extractEventDateLearning(body, subject) {
  const datePatterns = [
    /event\s+date[:\s]+([A-Za-z]+\s+\d{1,2},?\s+\d{4})/i,
    /date[:\s]+([A-Za-z]+\s+\d{1,2},?\s+\d{4})/i,
    /(\d{1,2}\/\d{1,2}\/\d{4})/,
    /(\d{4}-\d{2}-\d{2})/,
    /(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},?\s+\d{4}/i
  ];
  
  const searchText = subject + ' ' + body;
  for (const pattern of datePatterns) {
    const match = searchText.match(pattern);
    if (match) return match[1];
  }
  
  return '';
}

/**
 * Extract pricing data
 */
function extractPricingDataLearning(body, subject) {
  const data = {};
  const searchText = subject + ' ' + body;
  
  // Total cost
  const costMatches = searchText.match(/\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/g);
  if (costMatches) {
    const amounts = costMatches.map(m => parseFloat(m.replace(/[$,]/g, ''))).filter(n => n > 10);
    if (amounts.length > 0) {
      data.totalCost = Math.max(...amounts);
    }
  }
  
  // Rate
  const rateMatch = searchText.match(/(?:rate|hourly|per\s+hour)[:\s]+\$?(\d+(?:\.\d+)?)/i);
  if (rateMatch) {
    data.rate = parseFloat(rateMatch[1]);
  }
  
  // Deposit
  const depositMatch = searchText.match(/deposit[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)/i);
  if (depositMatch) {
    data.deposit = parseFloat(depositMatch[1].replace(/[$,]/g, ''));
  }
  
  return data;
}

/**
 * Step 4: Identify unique clients and conversations
 */
function identifyUniqueClients() {
  Logger.log('Identifying unique clients and conversations...');
  
  const ss = getSpreadsheet();
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  const clientSheet = ss.getSheetByName(SHEETS.CLIENT_MATCHING);
  const conversationsSheet = ss.getSheetByName(SHEETS.CONVERSATIONS);
  
  if (rawSheet.getLastRow() <= 1) return;
  
  const data = rawSheet.getDataRange().getValues();
  const rows = data.slice(1).filter(row => row[6] === 'No'); // Non-test (column 7, 0-indexed 6)
  
  // Setup client matching sheet
  clientSheet.clear();
  clientSheet.getRange(1, 1, 1, 6).setValues([[
    'Client Email', 'Alternative Emails', 'Submission Count', 
    'First Seen', 'Last Seen', 'Confidence'
  ]]);
  clientSheet.getRange(1, 1, 1, 6).setFontWeight('bold');
  
  // Group by email (exact match)
  const emailGroups = new Map();
  const conversationMap = new Map();
  
  rows.forEach(row => {
    const email = row[8] || ''; // Client email column (column 9, 0-indexed 8)
    const emailType = row[18] || ''; // Email Type column (column 19, 0-indexed 18)
    const conversationId = row[19] || ''; // Conversation ID column (column 20, 0-indexed 19)
    
    if (!email) return;
    
    const normalizedEmail = email.toLowerCase().trim();
    
    if (!emailGroups.has(normalizedEmail)) {
      emailGroups.set(normalizedEmail, {
        email: normalizedEmail,
        alternatives: new Set(),
        submissions: [],
        stlphEmails: 0,
        otherEmails: 0,
        conversations: new Set(),
        firstDate: new Date(row[2]), // Date column (column 3, 0-indexed 2)
        lastDate: new Date(row[2])
      });
    }
    
    const group = emailGroups.get(normalizedEmail);
    const date = new Date(row[2]);
    if (date < group.firstDate) group.firstDate = date;
    if (date > group.lastDate) group.lastDate = date;
    group.submissions.push(row);
    
    // Count email types
    if (emailType === 'STLPH') {
      group.stlphEmails++;
    } else {
      group.otherEmails++;
    }
    
    // Track conversations
    if (conversationId) {
      group.conversations.add(conversationId);
      
      // Build conversation map
      if (!conversationMap.has(conversationId)) {
        conversationMap.set(conversationId, {
          id: conversationId,
          clientEmail: normalizedEmail,
          subject: row[4] || '', // Subject column (column 5, 0-indexed 4)
          emails: [],
          firstDate: date,
          lastDate: date,
          stlphCount: 0,
          otherCount: 0
        });
      }
      
      const conv = conversationMap.get(conversationId);
      conv.emails.push(row);
      if (date < conv.firstDate) conv.firstDate = date;
      if (date > conv.lastDate) conv.lastDate = date;
      if (emailType === 'STLPH') {
        conv.stlphCount++;
      } else {
        conv.otherCount++;
      }
    }
    
    // Look for similar emails (same name, different domain, etc.)
    const fromEmail = row[3] || ''; // From column (column 4, 0-indexed 3)
    if (fromEmail && fromEmail !== email) {
      group.alternatives.add(fromEmail);
    }
  });
  
  // Write client groups with enhanced data
  const clientData = Array.from(emailGroups.values()).map(group => {
    const totalEmails = group.submissions.length;
    const conversationCount = group.conversations.size;
    const avgEmailsPerConv = conversationCount > 0 ? (totalEmails / conversationCount).toFixed(2) : totalEmails;
    
    return [
      group.email,
      Array.from(group.alternatives).join(', '),
      totalEmails,
      group.firstDate,
      group.lastDate,
      group.submissions.length > 1 ? 'High' : 'Medium',
      group.stlphEmails,
      group.otherEmails,
      conversationCount,
      avgEmailsPerConv
    ];
  });
  
  if (clientData.length > 0) {
    clientSheet.getRange(2, 1, clientData.length, 10).setValues(clientData);
  }
  
  // Write conversations
  setupConversationsSheet(conversationsSheet);
  conversationsSheet.getRange(2, 1, conversationsSheet.getLastRow() - 1, 10).clear();
  
  const conversationData = Array.from(conversationMap.values())
    .sort((a, b) => b.lastDate - a.lastDate)
    .map(conv => {
      const status = conv.stlphCount > 0 ? 'STLPH Lead' : 'Other';
      const outcome = conv.emails.some(e => e[7] === 'Yes') ? 'Confirmed' : 'Pending'; // Is Confirmation column
      
      return [
        conv.id,
        conv.clientEmail,
        conv.subject,
        conv.firstDate,
        conv.lastDate,
        conv.emails.length,
        conv.stlphCount,
        conv.otherCount,
        status,
        outcome
      ];
    });
  
  if (conversationData.length > 0) {
    conversationsSheet.getRange(2, 1, conversationData.length, 10).setValues(conversationData);
  }
  
  // Generate email classification summary
  generateEmailClassificationSummary(ss, emailGroups);
  
  Logger.log(`Identified ${emailGroups.size} unique client emails`);
  Logger.log(`Identified ${conversationMap.size} conversations`);
}

/**
 * Generate email classification summary
 */
function generateEmailClassificationSummary(ss, emailGroups) {
  const classificationSheet = ss.getSheetByName(SHEETS.EMAIL_CLASSIFICATION);
  setupEmailClassificationSheet(classificationSheet);
  classificationSheet.getRange(2, 1, classificationSheet.getLastRow() - 1, 6).clear();
  
  // Calculate totals
  let totalSTLPH = 0;
  let totalOther = 0;
  let totalClients = emailGroups.size;
  
  emailGroups.forEach(group => {
    totalSTLPH += group.stlphEmails;
    totalOther += group.otherEmails;
  });
  
  const totalEmails = totalSTLPH + totalOther;
  const stlphPercentage = totalEmails > 0 ? ((totalSTLPH / totalEmails) * 100).toFixed(1) : 0;
  const otherPercentage = totalEmails > 0 ? ((totalOther / totalEmails) * 100).toFixed(1) : 0;
  
  const avgSTLPHPerClient = totalClients > 0 ? (totalSTLPH / totalClients).toFixed(2) : 0;
  const avgOtherPerClient = totalClients > 0 ? (totalOther / totalClients).toFixed(2) : 0;
  
  const classificationData = [
    ['STLPH', totalSTLPH, `${stlphPercentage}%`, avgSTLPHPerClient, 'Business-related emails', ''],
    ['Other', totalOther, `${otherPercentage}%`, avgOtherPerClient, 'Marketing, personal, etc.', ''],
    ['Total', totalEmails, '100%', ((totalEmails / totalClients) || 0).toFixed(2), '', '']
  ];
  
  classificationSheet.getRange(2, 1, classificationData.length, 6).setValues(classificationData);
  
  Logger.log(`Email Classification: ${totalSTLPH} STLPH (${stlphPercentage}%), ${totalOther} Other (${otherPercentage}%)`);
}

/**
 * Process raw data with learning
 */
function processRawDataLearning() {
  // Use similar logic to original processRawData but with learned patterns
  // This would call the same aggregation functions
  processRawData(); // Reuse original for now
}

function setupRawDataSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 17).setValues([[
      'Email ID', 'Date', 'From Email', 'Subject', 'Body Preview', 
      'Is Test', 'Is Confirmation', 'Client Email', 'Event Date', 
      'Total Cost', 'Rate', 'Hours', 'Helpers', 'Occasion', 'Status',
      'Guests', 'Deposit'
    ]]);
    sheet.getRange(1, 1, 1, 17).setFontWeight('bold');
    sheet.setFrozenRows(1);
  }
}

function setupLogSheet(sheet) {
  if (sheet.getLastRow() === 0) {
    sheet.getRange(1, 1, 1, 4).setValues([['Timestamp', 'Status', 'Emails Processed', 'Notes']]);
    sheet.getRange(1, 1, 1, 4).setFontWeight('bold');
  }
}

function logMessage(logSheet, message, count) {
  const row = logSheet.getLastRow() + 1;
  logSheet.getRange(row, 1, 1, 4).setValues([[new Date(), message, count, '']]);
  Logger.log(`${new Date()}: ${message} (${count} emails)`);
}

function isAlreadyProcessed(sheet, messageId) {
  const data = sheet.getDataRange().getValues();
  return data.some(row => row[0] === messageId);
}

function getLastProcessedId(sheet) {
  const lastRow = sheet.getLastRow();
  if (lastRow > 1) {
    return sheet.getRange(lastRow, 1).getValue();
  }
  return null;
}

function writeBatchToSheet(sheet, batch) {
  if (batch.length === 0) return;
  const lastRow = sheet.getLastRow();
  sheet.getRange(lastRow + 1, 1, batch.length, batch[0].length).setValues(batch);
}

/**
 * Process raw data and create summary sheets (from original script)
 */
function processRawData() {
  const ss = getSpreadsheet();
  const rawSheet = ss.getSheetByName(SHEETS.RAW_DATA);
  const clientsSheet = ss.getSheetByName(SHEETS.CLIENTS);
  const monthlySheet = ss.getSheetByName(SHEETS.MONTHLY_REVENUE);
  const yearlySheet = ss.getSheetByName(SHEETS.YEARLY_SUMMARY);
  
  if (rawSheet.getLastRow() <= 1) return;
  
  const data = rawSheet.getDataRange().getValues();
  const rows = data.slice(1).filter(row => row[6] === 'No'); // Non-test submissions (column 7, 0-indexed 6)
  
  // Group by client email
  const clientsMap = new Map();
  
  rows.forEach(row => {
    const clientEmail = row[8] || ''; // Column 9 (Client Email)
    const emailType = row[18] || ''; // Column 19 (Email Type)
    const conversationId = row[19] || ''; // Column 20 (Conversation ID)
    
    if (!clientEmail) return;
    
    if (!clientsMap.has(clientEmail)) {
      clientsMap.set(clientEmail, {
        email: clientEmail,
        firstDate: new Date(row[2]), // Column 3 (Date)
        lastDate: new Date(row[2]),
        submissions: [],
        confirmed: 0,
        unconfirmed: 0,
        totalRevenue: 0,
        stlphEmails: 0,
        otherEmails: 0,
        conversations: new Set()
      });
    }
    
    const client = clientsMap.get(clientEmail);
    const date = new Date(row[2]); // Column 3 (Date)
    if (date < client.firstDate) client.firstDate = date;
    if (date > client.lastDate) client.lastDate = date;
    
    client.submissions.push(row);
    
    // Count email types
    if (emailType === 'STLPH') {
      client.stlphEmails++;
    } else {
      client.otherEmails++;
    }
    
    // Track conversations
    if (conversationId) {
      client.conversations.add(conversationId);
    }
    
    const cost = parseFloat(row[10]) || 0; // Column 11 (Total Cost)
    client.totalRevenue += cost;
    
    if (row[7] === 'Yes') { // Column 8 (Is Confirmation)
      client.confirmed++;
    } else {
      client.unconfirmed++;
    }
  });
  
  // Write clients data with enhanced metrics
  clientsSheet.getRange(2, 1, clientsSheet.getLastRow(), 12).clear();
  const clientsData = Array.from(clientsMap.values()).map(client => {
    const conversationCount = client.conversations.size;
    const avgEmailsPerConv = conversationCount > 0 
      ? (client.submissions.length / conversationCount).toFixed(2) 
      : client.submissions.length;
    
    return [
      client.email,
      client.firstDate,
      client.lastDate,
      client.submissions.length,
      client.stlphEmails,
      client.otherEmails,
      client.totalRevenue,
      client.confirmed,
      client.unconfirmed,
      conversationCount,
      avgEmailsPerConv,
      client.confirmed > 0 ? 'Active' : 'Lead'
    ];
  });
  
  if (clientsData.length > 0) {
    clientsSheet.getRange(2, 1, clientsData.length, 12).setValues(clientsData);
  }
  
  // Group by year and month
  const monthlyMap = new Map();
  
  rows.forEach(row => {
    const date = new Date(row[2]); // Column 3 (Date)
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
    const cost = parseFloat(row[10]) || 0; // Column 11 (Total Cost)
    monthData.revenue += cost;
    
    if (row[7] === 'Yes') { // Column 8 (Is Confirmation)
      monthData.confirmed++;
    } else {
      monthData.unconfirmed++;
    }
    
    const rate = parseFloat(row[11]) || 0; // Column 12 (Rate)
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
  
  rows.forEach(row => {
    const date = new Date(row[2]); // Column 3 (Date)
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
    const cost = parseFloat(row[10]) || 0; // Column 11 (Total Cost)
    yearData.revenue += cost;
    
    if (row[7] === 'Yes') { // Column 8 (Is Confirmation)
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
}

/**
 * Format the sheets for better readability
 */
function formatSheets(ss) {
  const monthlySheet = ss.getSheetByName(SHEETS.MONTHLY_REVENUE);
  if (monthlySheet.getLastRow() > 1) {
    const monthlyRange = monthlySheet.getRange(2, 3, monthlySheet.getLastRow() - 1, 7);
    monthlyRange.setNumberFormat('$#,##0.00');
  }
  
  const yearlySheet = ss.getSheetByName(SHEETS.YEARLY_SUMMARY);
  if (yearlySheet.getLastRow() > 1) {
    const yearlyRange = yearlySheet.getRange(2, 2, yearlySheet.getLastRow() - 1, 5);
    yearlyRange.setNumberFormat('$#,##0.00');
  }
  
  [monthlySheet, yearlySheet, ss.getSheetByName(SHEETS.CLIENTS)].forEach(sheet => {
    if (sheet && sheet.getLastRow() > 1) {
      sheet.autoResizeColumns(1, sheet.getLastColumn());
    }
  });
}
