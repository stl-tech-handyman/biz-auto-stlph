/**
 * Unit tests for auth-related utils: token extraction and validation.
 */
function test_authUtils_all() {
  console.log("🔐 Running auth utils tests...");

  test_getExpectedZapierToken();
  test_extractTokenFromRequest_valid();
  test_extractTokenFromRequest_missing();
  test_isValidApiToken_validZapier();
  test_isValidApiToken_validPostman();
  test_isValidApiToken_invalid();
  test_isValidApiToken_empty();

  console.log("✅ Auth utils tests complete.");
}


/**
 * Test: getExpectedZapierToken returns the correct script property.
 */
function test_getExpectedZapierToken() {
  const expected = "zapier-abc123";
  PropertiesService.getScriptProperties().setProperty("ZAPIER_TOKEN", expected);

  const actual = getExpectedZapierToken();
  logTestResult("getExpectedZapierToken", 1, "should retrieve stored token", actual === expected, actual, expected);
}


/**
 * Test: extractTokenFromRequest should return token if present
 */
function test_extractTokenFromRequest_valid() {
  const e = { parameter: { token: "zapier-abc123" } };
  const actual = extractTokenFromRequest(e);
  const expected = "zapier-abc123";

  logTestResult("extractTokenFromRequest", 1, "should extract token from event", actual === expected, actual, expected);
}


/**
 * Test: extractTokenFromRequest should return null if token missing
 */
function test_extractTokenFromRequest_missing() {
  const e = { parameter: {} };
  const actual = extractTokenFromRequest(e);
  const expected = null;

  logTestResult("extractTokenFromRequest", 2, "should return null when token is missing", actual === expected, actual, expected);
}


/**
 * Test: isValidApiToken should validate zapier token correctly
 */
function test_isValidApiToken_validZapier() {
  const token = "zapier-abc123";
  PropertiesService.getScriptProperties().setProperty("ZAPIER_TOKEN", token);

  const result = isValidApiToken(token);
  const passed = result.valid && result.source === "ZAPIER_TOKEN" && result.preview === "zapi";

  logTestResult("isValidApiToken", 1, "should recognize valid zapier token", passed, JSON.stringify(result));
}


/**
 * Test: isValidApiToken should validate postman token correctly
 */
function test_isValidApiToken_validPostman() {
  const token = "post-xyz789";
  PropertiesService.getScriptProperties().setProperty("POSTMAN_TOKEN", token);

  const result = isValidApiToken(token);
  const passed = result.valid && result.source === "POSTMAN_TOKEN" && result.preview === "post";

  logTestResult("isValidApiToken", 2, "should recognize valid postman token", passed, JSON.stringify(result));
}


/**
 * Test: isValidApiToken should return false for unknown token
 */
function test_isValidApiToken_invalid() {
  const result = isValidApiToken("invalid-111");
  const passed = !result.valid && result.source === null;

  logTestResult("isValidApiToken", 3, "should reject unknown token", passed, JSON.stringify(result));
}


/**
 * Test: isValidApiToken should return false for empty input
 */
function test_isValidApiToken_empty() {
  const result = isValidApiToken("");
  const passed = !result.valid && result.source === null;

  logTestResult("isValidApiToken", 4, "should reject empty token", passed, JSON.stringify(result));
}
