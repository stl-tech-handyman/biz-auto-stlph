// ──────────────────────────────────────────────────────────────────────────────
// HTTP / Routing utilities
// ──────────────────────────────────────────────────────────────────────────────

/**
 * Extracts API version from request (?v=2 or ?version=2 or ?version=v2).
 * Returns a positive integer; defaults to 1.
 * @param {Object} e - Apps Script event object
 * @returns {number}
 */
function getApiVersion(e) {
  var raw = (e && e.parameter && (e.parameter.v || e.parameter.version)) || "";
  if (typeof raw === "string" && /^v\d+$/i.test(raw)) {
    raw = raw.slice(1); // strip leading 'v'
  }
  var version = Number(raw);
  return (Number.isInteger(version) && version > 0) ? version : 1;
}

/**
 * Extracts the action parameter. Throws if missing.
 * @param {Object} e
 * @returns {string}
 */
function getAction(e) {
  var action = e && e.parameter && e.parameter.action;
  if (!action) throw new Error("Missing required parameter: 'action'");
  return String(action);
}

/**
 * Extracts "req-initiator" parameter. Returns "Unknown" and logs a warning if not present.
 * @param {Object} e
 * @returns {string}
 */
function extractInitiator(e) {
  var initiator = e && e.parameter && e.parameter["req-initiator"];
  if (!initiator) {
    console.warn("⚠️ Missing req-initiator param.");
    return "Unknown";
  }
  console.log("🔍 Initiator: " + initiator);
  return String(initiator);
}

/**
 * Generic parameter accessor with optional default/required/coercion.
 * @param {Object} e - Apps Script event
 * @param {string} name - Param name
 * @param {Object} [opts] - { required?: boolean, defaultValue?: any, coerce?: (val:string)=>any }
 * @returns {*}
 */
function getParam(e, name, opts) {
  var val = e && e.parameter && e.parameter[name];
  if ((val === undefined || val === null || val === "") && opts && opts.hasOwnProperty("defaultValue")) {
    return opts.defaultValue;
  }
  if ((val === undefined || val === null || val === "") && opts && opts.required) {
    throw new Error("Missing required parameter: '" + name + "'");
  }
  if (opts && typeof opts.coerce === "function") {
    return opts.coerce(val);
  }
  return val;
}

/**
 * Convenience wrapper for required string params (trimmed).
 * @param {Object} e
 * @param {string} name
 * @returns {string}
 */
function getRequiredParam(e, name) {
  var v = getParam(e, name, { required: true });
  v = String(v).trim();
  if (!v) throw new Error("Parameter '" + name + "' must be a non-empty string.");
  return v;
}
/**
 * Creates a standardized JSON **success** response.
 *
 * Structure:
 * {
 *   ok: true,             // indicates success
 *   data: <your data>,    // the payload (object, array, string, etc.)
 *   meta: { ... }         // optional metadata about the request
 * }
 *
 * @param {*} data - The payload to send back. Will be `null` if falsy (except 0/false).
 * @param {Object} [meta] - Optional metadata (e.g., { action: "GET_LAT_LNG", v: 1 })
 * @returns {ContentService.Output} - Apps Script HTTP JSON response
 *
 * Usage:
 *   return jsonOk({ lat: 38.63, lng: -90.2 }, { v: 1 });
 */
function jsonOk(data, meta) {
  var body = { ok: true, data: data || null };
  if (meta) body.meta = meta;

  return ContentService.createTextOutput(JSON.stringify(body))
    .setMimeType(ContentService.MimeType.JSON);
}

/**
 * Creates a standardized JSON **error** response.
 *
 * Structure:
 * {
 *   ok: false,             // indicates failure
 *   error: {
 *     message: "Readable explanation of the error",
 *     code: "OPTIONAL_CODE" // optional short code or HTTP status
 *   },
 *   meta: { ... }          // optional metadata about the request
 * }
 *
 * @param {string|Error} message - Error message (or Error object)
 * @param {string|number} [code] - Optional error code (can be string or numeric)
 * @param {Object} [meta] - Optional metadata (e.g., { action: "GET_LAT_LNG", v: 1, errorId: "abc123" })
 * @returns {ContentService.Output} - Apps Script HTTP JSON response
 *
 * Usage:
 *   return jsonError("Unauthorized", 401, { v: 1 });
 *   return jsonError(new Error("Geocoding failed"), "GEOCODE_ERROR");
 */
function jsonError(message, code, meta) {
  // Use the message property if an Error object was passed
  var msg = (message && message.message) || String(message || "Unknown error");

  var body = { ok: false, error: { message: msg } };

  // Include code if provided (HTTP status code or custom code)
  if (code !== undefined && code !== null) {
    body.error.code = code;
  }

  // Include any additional metadata
  if (meta) body.meta = meta;

  return ContentService.createTextOutput(JSON.stringify(body))
    .setMimeType(ContentService.MimeType.JSON);
}

/**
 * Raw JSON passthrough (kept for back-compat).
 * @param {Object|string} obj
 */
function json(obj) {
  return ContentService.createTextOutput(JSON.stringify(obj))
    .setMimeType(ContentService.MimeType.JSON);
}

/**
 * Plain text response.
 * @param {string} msg
 */
function text(msg) {
  return ContentService.createTextOutput(String(msg))
    .setMimeType(ContentService.MimeType.TEXT);
}

/**
 * Test logging helper.
 */
function logTestResult(fnName, index, description, success, actual, expected) {
  var status = success ? "✅ PASSED" : "❌ FAILED";
  console.log(
    status + " [" + fnName + " #" + index + "] " + description + "\n" +
    "→ Expected: " + expected + "\n" +
    "→ Got: " + actual
  );
}
