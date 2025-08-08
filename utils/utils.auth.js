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

/**
 * Validate token and return a safe object for logging + checks.
 */
function isValidApiToken(token) {
  const expected = getExpectedApiToken_();
  const valid = Boolean(token && expected && token === expected);
  return {
    valid,
    preview: token ? String(token).slice(0, 4) : null,
    source: null // filled by caller from extractTokenFromRequest
  };
}
