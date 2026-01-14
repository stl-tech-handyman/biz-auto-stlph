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
// Backwards-compatible single spreadsheet ID:
const SPREADSHEET_ID_PROPERTY = 'EMAIL_ANALYZER_SPREADSHEET_ID';
// New multi-spreadsheet mode (JSON array of IDs OR comma-separated):
const SPREADSHEET_IDS_PROPERTY = 'EMAIL_ANALYZER_SPREADSHEET_IDS';

// How long after Expires At we consider the lock "certainly not in use".
// (Gives some buffer for clock skew / transient failures.)
const EXPIRE_GRACE_SECONDS = 15;

// If true, delete expired rows instead of marking them EXPIRED.
// Marking keeps an audit trail and is safer by default.
const DELETE_EXPIRED_ROWS = false;

const LOCKS_SHEET_NAME = 'Locks';

const REQUIRED_HEADERS = [
  'Agent ID',
  'Created At',
  'Expires At',
  'Status',
  // Optional audit columns (added if missing):
  'Cleaned At',
  'Cleanup Reason',
];

function _now() {
  return new Date();
}

function _parseSpreadsheetIds_(properties) {
  const multi = properties.getProperty(SPREADSHEET_IDS_PROPERTY);
  if (multi && multi.trim()) {
    // Accept either JSON array or comma-separated.
    const trimmed = multi.trim();
    if (trimmed.startsWith('[')) {
      try {
        const ids = JSON.parse(trimmed);
        if (Array.isArray(ids)) return ids.map(String).map(s => s.trim()).filter(Boolean);
      } catch (e) {
        // fall through to comma parsing
      }
    }
    return trimmed.split(',').map(s => s.trim()).filter(Boolean);
  }

  const single = properties.getProperty(SPREADSHEET_ID_PROPERTY);
  if (single && single.trim()) return [single.trim()];
  return [];
}

function _ensureLocksHeaders_(sheet) {
  const headerRange = sheet.getRange(1, 1, 1, Math.max(REQUIRED_HEADERS.length, sheet.getLastColumn() || 1));
  const current = headerRange.getValues()[0];

  // Build a map from header name -> col index (1-based)
  const headerMap = {};
  for (let i = 0; i < current.length; i++) {
    const name = String(current[i] || '').trim();
    if (name) headerMap[name] = i + 1;
  }

  // Ensure required headers exist (append missing)
  let lastCol = sheet.getLastColumn();
  REQUIRED_HEADERS.forEach(h => {
    if (!headerMap[h]) {
      lastCol++;
      sheet.getRange(1, lastCol).setValue(h);
      headerMap[h] = lastCol;
    }
  });

  return headerMap;
}

function _parseDateSafe_(value) {
  if (!value) return { ok: false, date: null };
  // If it's already a Date object (Sheets often returns Date for ISO-ish strings)
  if (Object.prototype.toString.call(value) === '[object Date]') {
    if (!isNaN(value.getTime())) return { ok: true, date: value };
    return { ok: false, date: null };
  }
  const d = new Date(String(value));
  if (!isNaN(d.getTime())) return { ok: true, date: d };
  return { ok: false, date: null };
}

/**
 * Main cleanup function - runs on schedule
 */
function cleanupExpiredLocks() {
  const properties = PropertiesService.getScriptProperties();
  const spreadsheetIds = _parseSpreadsheetIds_(properties);
  
  if (!spreadsheetIds.length) {
    Logger.log(`No spreadsheet IDs configured. Set ${SPREADSHEET_IDS_PROPERTY} (JSON array or comma-separated) or ${SPREADSHEET_ID_PROPERTY}.`);
    return;
  }
  
  const now = _now();
  const graceMs = EXPIRE_GRACE_SECONDS * 1000;

  spreadsheetIds.forEach(spreadsheetId => {
    try {
      const ss = SpreadsheetApp.openById(spreadsheetId);
      const locksSheet = ss.getSheetByName(LOCKS_SHEET_NAME);
      
      if (!locksSheet) {
        Logger.log(`[${spreadsheetId}] Locks sheet not found`);
        return;
      }

      const headerMap = _ensureLocksHeaders_(locksSheet);
      const lastRow = locksSheet.getLastRow();
      if (lastRow <= 1) return; // only header

      const lastCol = locksSheet.getLastColumn();
      const data = locksSheet.getRange(2, 1, lastRow - 1, lastCol).getValues();

      const idxAgent = headerMap['Agent ID'] - 1;
      const idxCreated = headerMap['Created At'] - 1;
      const idxExpires = headerMap['Expires At'] - 1;
      const idxStatus = headerMap['Status'] - 1;
      const idxCleanedAt = headerMap['Cleaned At'] - 1;
      const idxReason = headerMap['Cleanup Reason'] - 1;

      const rowsToExpire = [];
      const rowsToDelete = [];

      for (let r = 0; r < data.length; r++) {
        const row = data[r];
        const status = String(row[idxStatus] || '').trim();
        if (status !== 'ACTIVE') continue;

        const parsedExpires = _parseDateSafe_(row[idxExpires]);
        const agentId = String(row[idxAgent] || '').trim();

        if (!parsedExpires.ok) {
          // If we cannot parse Expires At, it's unsafe to keep as ACTIVE.
          rowsToExpire.push({
            rowNumber: r + 2,
            agentId,
            reason: 'invalid_expires_at',
          });
          continue;
        }

        const expiresAt = parsedExpires.date;
        // "Certainly not in use" = ExpiresAt passed + grace window
        if (now.getTime() > expiresAt.getTime() + graceMs) {
          rowsToExpire.push({
            rowNumber: r + 2,
            agentId,
            reason: 'expires_at_passed',
          });
        }
      }

      // Apply updates
      if (DELETE_EXPIRED_ROWS) {
        // Delete in reverse order to keep row numbers stable
        rowsToExpire.sort((a, b) => b.rowNumber - a.rowNumber).forEach(item => {
          locksSheet.deleteRow(item.rowNumber);
          Logger.log(`[${spreadsheetId}] Deleted expired lock row: Agent ${item.agentId} (${item.reason})`);
        });
      } else {
        rowsToExpire.forEach(item => {
          locksSheet.getRange(item.rowNumber, headerMap['Status']).setValue('EXPIRED');
          locksSheet.getRange(item.rowNumber, headerMap['Cleaned At']).setValue(new Date().toISOString());
          locksSheet.getRange(item.rowNumber, headerMap['Cleanup Reason']).setValue(item.reason);
          Logger.log(`[${spreadsheetId}] Marked lock EXPIRED: Agent ${item.agentId} (${item.reason})`);
        });
      }

      if (rowsToExpire.length) {
        Logger.log(`[${spreadsheetId}] Cleaned up ${rowsToExpire.length} stale locks`);
      }

    } catch (error) {
      Logger.log(`[${spreadsheetId}] Error cleaning up locks: ${error}`);
    }
  });
}

/**
 * Setup function - run once to configure
 */
function setupLockCleanup() {
  const properties = PropertiesService.getScriptProperties();

  const ids = _parseSpreadsheetIds_(properties);
  if (!ids.length) {
    Logger.log(`Please set ${SPREADSHEET_IDS_PROPERTY} (JSON array or comma-separated) OR ${SPREADSHEET_ID_PROPERTY} in script properties`);
    Logger.log('Go to: Project settings > Script properties');
    Logger.log(`Add: ${SPREADSHEET_IDS_PROPERTY} = ["spreadsheet-id-1","spreadsheet-id-2"]`);
    return;
  }

  Logger.log(`Lock cleanup configured for ${ids.length} spreadsheet(s):`);
  ids.forEach(id => Logger.log(`  - ${id}`));

  Logger.log('Now set up trigger: Triggers (clock icon) > Add Trigger');
  Logger.log('Function: cleanupExpiredLocks | Event source: Time-driven | Type: Minutes timer | Every minute');
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
  const spreadsheetIds = _parseSpreadsheetIds_(properties);
  
  if (!spreadsheetIds.length) {
    Logger.log('No spreadsheet IDs configured');
    return;
  }
  
  const now = _now();
  spreadsheetIds.forEach(spreadsheetId => {
    try {
      const ss = SpreadsheetApp.openById(spreadsheetId);
      const locksSheet = ss.getSheetByName(LOCKS_SHEET_NAME);
      
      if (!locksSheet) {
        Logger.log(`[${spreadsheetId}] Locks sheet not found`);
        return;
      }
      
      const headerMap = _ensureLocksHeaders_(locksSheet);
      const lastRow = locksSheet.getLastRow();
      if (lastRow <= 1) {
        Logger.log(`[${spreadsheetId}] No locks found`);
        return;
      }
      
      const lastCol = locksSheet.getLastColumn();
      const data = locksSheet.getRange(2, 1, lastRow - 1, lastCol).getValues();
      const idxAgent = headerMap['Agent ID'] - 1;
      const idxExpires = headerMap['Expires At'] - 1;
      const idxStatus = headerMap['Status'] - 1;
      
      Logger.log(`[${spreadsheetId}] Current lock status:`);
      data.forEach(row => {
        const agentId = String(row[idxAgent] || '').trim();
        const status = String(row[idxStatus] || '').trim();
        const parsedExpires = _parseDateSafe_(row[idxExpires]);
        let timeLeft = 'unknown';
        if (parsedExpires.ok) {
          const ms = parsedExpires.date.getTime() - now.getTime();
          timeLeft = ms < 0 ? 'EXPIRED' : `${Math.round(ms / 1000)}s`;
        }
        Logger.log(`  Agent: ${agentId} | Status: ${status} | Time left: ${timeLeft}`);
      });
      
    } catch (error) {
      Logger.log(`[${spreadsheetId}] Error getting lock status: ${error}`);
    }
  });
}
