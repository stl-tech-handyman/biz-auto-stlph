// Routing and response utilities
/**
 * Extracts the API version number from the request query parameter (?v=2)
 * Defaults to 1 if not provided or invalid.
 *
 * @param {Object} e - Apps Script event object
 * @returns {number} - Parsed version (default: 1)
 */
/**
 * Extracts the API version from request (e.g., ?v=2), defaults to 1.
 * Ensures it's a whole positive integer, otherwise falls back to 1.
 * @param {Object} e - Apps Script event object
 * @returns {number}
 */
function getApiVersion(e) {
  const raw = e?.parameter?.v;
  const version = Number(raw);

  // Accept only whole positive integers
  if (Number.isInteger(version) && version > 0) {
    return version;
  }

  return 1;
}

/**
 * Extracts the "req-initiator" parameter from a request event object (e).
 * If missing, returns "Unknown" and logs a warning.
 * 
 * @param {Object} e - The request event object passed to doGet/doPost.
 * @returns {string} - The initiator identifier (usually email, service name, etc).
 * 
 * Usage:
 * const initiator = extractInitiator(e);
 */
function extractInitiator(e) {
  const initiator = e?.parameter?.["req-initiator"];
  if (!initiator) {
    console.warn("⚠️ Missing req-initiator param.");
    return "Unknown";
  }
  console.log(`🔍 Initiator: ${initiator}`);
  return initiator;
}

function logTestResult(fnName, index, description, success, actual, expected) {
  const status = success ? "✅ PASSED" : "❌ FAILED";
  console.log(`${status} [${fnName} #${index}] ${description}
→ Expected: ${expected}
→ Got: ${actual}`);
}

/**
 * Extracts the `action` parameter from the request object (e).
 * Throws an error if missing.
 * 
 * @param {Object} e - Apps Script event object
 * @returns {string} - The requested action (e.g. 'sendQuote')
 */
function getAction(e) {
  const action = e?.parameter?.action;
  if (!action) {
    throw new Error("Missing required parameter: 'action'");
  }
  return action;
}

/**
 * Extracts the API version (e.g., 'v1', 'v2') from the request object (e).
 * Defaults to 'v1' if not provided.
 * 
 * @param {Object} e - Apps Script event object
 * @returns {string} - The version string
 */
function getVersion(e) {
  const version = e?.parameter?.version || 'v1';
  return version;
}

/**
 * Wraps a JS object or string in a JSON response.
 * @param {Object|string} obj - Object to return as JSON.
 * @returns {ContentService.Output}
 */
function json(obj) {
  return ContentService.createTextOutput(JSON.stringify(obj))
    .setMimeType(ContentService.MimeType.JSON);
}

/**
 * Returns plain text output.
 * @param {string} msg - Message string
 * @returns {ContentService.Output}
 */
function text(msg) {
  return ContentService.createTextOutput(msg)
    .setMimeType(ContentService.MimeType.TEXT);
}
