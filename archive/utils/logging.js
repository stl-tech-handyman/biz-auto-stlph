/**
 * Structured logging helpers
 * - One-line JSON logs via console.* so Cloud Logging can filter easily
 * - Safe request snapshot with masking of sensitive fields
 */

const LOGGING_CONFIG = {
    enableDebug: true,        // turn off in prod to silence DEBUG
    defaultLevel: "INFO",
    // keys (case-insensitive substrings) to mask when logging
    redactHints: ["token", "authorization", "apikey", "password", "secret", "key"],
    // include request snapshot log? (call logRequestSummary_ yourself)
    enableRequestSnapshot: true,
    // if the POST JSON body contains { headers: {...} }, log a masked preview of it
    allowBodyHeadersPreview: true
  };
  
  function newRequestId_() {
    return Utilities.getUuid().slice(0, 8);
  }
  
  /**
   * Core logger
   * @param {"INFO"|"WARN"|"ERROR"|"DEBUG"} level
   * @param {string} message - short, human-readable
   * @param {Object=} data   - attached structured fields
   * @param {string=} requestId
   */
  function log_(level, message, data, requestId) {
    if (level === "DEBUG" && !LOGGING_CONFIG.enableDebug) return;
  
    var payload = {
      event: "api_request",
      level: level || LOGGING_CONFIG.defaultLevel,
      ts: new Date().toISOString(),
      requestId: requestId || null,
      message: message
    };
    if (data && typeof data === "object") {
      for (var k in data) payload[k] = data[k];
    }
  
    var line = JSON.stringify(payload);
    switch (level) {
      case "ERROR": console.error(line); break;
      case "WARN":  console.warn(line);  break;
      default:      console.log(line);
    }
  }
  
  function logInfo(msg, data, reqId){ log_("INFO",  msg, data, reqId); }
  function logWarn(msg, data, reqId){ log_("WARN",  msg, data, reqId); }
  function logError(msg, data, reqId){log_("ERROR", msg, data, reqId); }
  function logDebug(msg, data, reqId){log_("DEBUG", msg, data, reqId); }
  
  /* =========================
     Small span helpers
     ========================= */
  function spanStart_(){ return Date.now(); }
  function spanEndMs_(t0){ return Date.now() - t0; }
  
  /* =========================
     Safe request snapshot
     ========================= */
  
  /** Mask sensitive values based on key name. */
  function maskSensitive_(key, value) {
    if (value === null || value === undefined) return value;
    var k = String(key).toLowerCase();
    for (var i = 0; i < LOGGING_CONFIG.redactHints.length; i++) {
      if (k.indexOf(LOGGING_CONFIG.redactHints[i]) !== -1) {
        var s = String(value);
        return (s.length <= 4) ? "****" : (s.slice(0, 4) + "...");
      }
    }
    return value;
  }
  
  function safeJsonParse_(text) {
    try { return JSON.parse(text || "{}"); } catch (_){ return null; }
  }
  
  /**
   * Build a safe, compact summary of the Apps Script event object.
   * - Masks sensitive param keys
   * - Shows param counts (for e.parameters) to avoid dumping big arrays
   * - Shows POST body type/length and masked JSON preview if JSON
   * - Optionally shows masked preview of body.headers if present
   */
  function buildRequestSnapshot_(e) {
    var isPost = !!(e && e.postData);
    var paramObj = (e && e.parameter) ? e.parameter : {};
    var params = {};
    Object.keys(paramObj || {}).forEach(function(k){
      params[k] = maskSensitive_(k, paramObj[k]);
    });
  
    var parameters = {};
    if (e && e.parameters) {
      Object.keys(e.parameters).forEach(function(k){
        var arr = e.parameters[k];
        parameters[k] = { count: Array.isArray(arr) ? arr.length : 0 };
      });
    }
  
    var bodyInfo = {};
    if (isPost && e.postData) {
      bodyInfo.type = e.postData.type || null; // e.g. "application/json"
      var contents = e.postData.contents || "";
      bodyInfo.length = (typeof e.postData.length === "number") ? e.postData.length : contents.length;
  
      if (bodyInfo.type && bodyInfo.type.indexOf("json") !== -1) {
        var obj = safeJsonParse_(contents);
        if (obj) {
          var masked = {};
          Object.keys(obj).forEach(function(k){ masked[k] = maskSensitive_(k, obj[k]); });
          // Optional masked headers preview if user includes them in body
          if (LOGGING_CONFIG.allowBodyHeadersPreview && obj.headers && typeof obj.headers === "object") {
            var maskedHeaders = {};
            Object.keys(obj.headers).forEach(function(h){ maskedHeaders[h] = maskSensitive_(h, obj.headers[h]); });
            masked.headers = maskedHeaders;
          }
          bodyInfo.jsonPreview = masked;
        } else {
          bodyInfo.jsonPreview = "parse_error";
        }
      } else {
        bodyInfo.preview = contents ? (contents.slice(0, 40) + (contents.length > 40 ? "..." : "")) : "";
      }
    }
  
    return {
      method: isPost ? "POST" : "GET",
      params: params,          // masked values
      parameters: parameters,  // counts only
      body: bodyInfo           // type/length/preview only
    };
  }
  
  /** Log a one-line safe snapshot of the request. */
  function logRequestSummary_(e, reqId) {
    try {
      var snap = buildRequestSnapshot_(e);
      logInfo("Request snapshot", snap, reqId);
    } catch (err) {
      logWarn("Request snapshot failed", { reason: String(err && err.message || err) }, reqId);
    }
  }
  