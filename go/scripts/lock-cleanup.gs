/**
 * Lock Cleanup Service for Email Analyzer
 * 
 * This Apps Script runs on a schedule to clean up expired locks.
 * Set up a time-driven trigger to run this every minute.
 * 
 * Setup:
 * 1. Open Apps Script editor
 * 2. Create new project or add to existing
 * 3. Copy this code
 * 4. Set up trigger: Edit > Current project's triggers
 * 5. Add trigger: cleanupExpiredLocks, Time-driven, Every minute
 */

// Configuration
const SPREADSHEET_ID_PROPERTY = 'EMAIL_ANALYZER_SPREADSHEET_ID'; // Set this in script properties

/**
 * Main cleanup function - runs on schedule
 */
function cleanupExpiredLocks() {
  const properties = PropertiesService.getScriptProperties();
  const spreadsheetId = properties.getProperty(SPREADSHEET_ID_PROPERTY);
  
  if (!spreadsheetId) {
    Logger.log('No spreadsheet ID configured. Set EMAIL_ANALYZER_SPREADSHEET_ID in script properties.');
    return;
  }
  
  try {
    const ss = SpreadsheetApp.openById(spreadsheetId);
    const locksSheet = ss.getSheetByName('Locks');
    
    if (!locksSheet) {
      Logger.log('Locks sheet not found');
      return;
    }
    
    const data = locksSheet.getDataRange().getValues();
    if (data.length <= 1) {
      return; // No locks
    }
    
    const now = new Date();
    let expiredCount = 0;
    const rowsToUpdate = [];
    
    // Check each lock (skip header row)
    for (let i = 1; i < data.length; i++) {
      const row = data[i];
      if (row.length < 4) continue;
      
      const status = row[3];
      if (status !== 'ACTIVE') continue;
      
      const expiresAtStr = row[2];
      if (!expiresAtStr) continue;
      
      const expiresAt = new Date(expiresAtStr);
      
      // If expired, mark as EXPIRED
      if (expiresAt < now) {
        rowsToUpdate.push({row: i + 1, agentId: row[0]});
        expiredCount++;
      }
    }
    
    // Update expired locks
    for (const update of rowsToUpdate) {
      locksSheet.getRange(update.row, 4).setValue('EXPIRED');
      Logger.log(`Marked lock as EXPIRED: Agent ${update.agentId}`);
    }
    
    if (expiredCount > 0) {
      Logger.log(`Cleaned up ${expiredCount} expired locks`);
    }
    
  } catch (error) {
    Logger.log(`Error cleaning up locks: ${error}`);
  }
}

/**
 * Setup function - run once to configure
 */
function setupLockCleanup() {
  const properties = PropertiesService.getScriptProperties();
  
  // Prompt for spreadsheet ID (you'll need to set this manually)
  // Or get it from the spreadsheet URL
  const spreadsheetId = properties.getProperty(SPREADSHEET_ID_PROPERTY);
  
  if (!spreadsheetId) {
    Logger.log('Please set EMAIL_ANALYZER_SPREADSHEET_ID in script properties');
    Logger.log('Go to: File > Project settings > Script properties');
    Logger.log('Add: EMAIL_ANALYZER_SPREADSHEET_ID = your-spreadsheet-id');
    return;
  }
  
  Logger.log(`Lock cleanup configured for spreadsheet: ${spreadsheetId}`);
  Logger.log('Set up trigger: Edit > Current project\'s triggers');
  Logger.log('Add: cleanupExpiredLocks, Time-driven, Every minute');
}

/**
 * Manual cleanup - can be run manually for testing
 */
function manualCleanup() {
  cleanupExpiredLocks();
}

/**
 * Get current lock status
 */
function getLockStatus() {
  const properties = PropertiesService.getScriptProperties();
  const spreadsheetId = properties.getProperty(SPREADSHEET_ID_PROPERTY);
  
  if (!spreadsheetId) {
    Logger.log('No spreadsheet ID configured');
    return;
  }
  
  try {
    const ss = SpreadsheetApp.openById(spreadsheetId);
    const locksSheet = ss.getSheetByName('Locks');
    
    if (!locksSheet) {
      Logger.log('Locks sheet not found');
      return;
    }
    
    const data = locksSheet.getDataRange().getValues();
    if (data.length <= 1) {
      Logger.log('No locks found');
      return;
    }
    
    const now = new Date();
    Logger.log('Current lock status:');
    
    for (let i = 1; i < data.length; i++) {
      const row = data[i];
      if (row.length < 4) continue;
      
      const agentId = row[0];
      const createdAt = row[1];
      const expiresAt = row[2];
      const status = row[3];
      
      const expiresDate = new Date(expiresAt);
      const isExpired = expiresDate < now;
      const timeLeft = isExpired ? 'EXPIRED' : `${Math.round((expiresDate - now) / 1000)}s`;
      
      Logger.log(`  Agent: ${agentId} | Status: ${status} | Expires: ${expiresAt} | Time left: ${timeLeft}`);
    }
    
  } catch (error) {
    Logger.log(`Error getting lock status: ${error}`);
  }
}
