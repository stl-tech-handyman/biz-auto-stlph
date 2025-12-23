/**
 * Read expected API token from Script Properties.
 */
function getExpectedApiToken_() {
  return getScriptProp(CONFIG_KEYS.ZAPIER_TOKEN); // single source of truth
}

/**
 * Extract token from query (?token=...), or Authorization header (Bearer X).
 * Returns both token and the source used (for logs).
 */
function extractTokenFromRequest(e) {
  // Query param
  const qp = e?.parameter?.token;
  if (qp) return { token: String(qp), source: "query" };

  // Header (webapp executions expose headers via e?.headers in some contexts; fallback if not present)
  const headers = (e && (e.headers || e.parameter)) || {};
  const auth = headers.Authorization || headers.authorization;
  if (auth && /^Bearer\s+/i.test(auth)) {
    return { token: auth.replace(/^Bearer\s+/i, "").trim(), source: "header" };
  }

  return { token: null, source: "none" };
}

function isValidApiToken(token) {
  var props = PropertiesService.getScriptProperties();

  // normalize the incoming token
  var tokenNorm = String(token == null ? "" : token).trim();

  // load & normalize known tokens
  var knownTokens = {
    ZAPIER_TOKEN:   props.getProperty("ZAPIER_TOKEN"),
    POSTMAN_TOKEN:  props.getProperty("POSTMAN_TOKEN"),
    FRONTEND_TOKEN: props.getProperty("FRONTEND_TOKEN")
  };

  // optional: one-time debug preview (safe)
  try {
    logDebug("Auth props preview", {
      ZAPIER_TOKEN:   (knownTokens.ZAPIER_TOKEN   || "").slice(0,4),
      POSTMAN_TOKEN:  (knownTokens.POSTMAN_TOKEN  || "").slice(0,4),
      FRONTEND_TOKEN: (knownTokens.FRONTEND_TOKEN || "").slice(0,4),
      incomingPreview: tokenNorm.slice(0,4)
    });
  } catch (_) {}

  for (var key in knownTokens) {
    var valueNorm = String(knownTokens[key] == null ? "" : knownTokens[key]).trim();
    if (tokenNorm && valueNorm && tokenNorm === valueNorm) {
      return { valid: true, source: key, preview: valueNorm.slice(0, 4) };
    }
  }

  return { valid: false, source: null, preview: tokenNorm.slice(0, 4) || "N/A" };
}
