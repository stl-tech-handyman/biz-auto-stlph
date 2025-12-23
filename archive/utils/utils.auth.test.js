/**
 * Unit tests for auth-related utils: token extraction and validation.
 */
function test_authUtils_all() {
  console.log("🔐 Running auth utils tests...");

  test_getExpectedZapierToken();
  test_extractTokenFromRequest_valid();
  test_extractTokenFromRequest_missing();
  test_isValidZapierToken_valid();
  test_isValidZapierToken_invalid();
  test_isValidZapierToken_empty();

  console.log("✅ Auth utils tests complete.");
}

/**
 * Test: getExpectedZapierToken returns the correct script property.
 */
function test_getExpectedZapierToken() {
  const expected = "secret123"; // 👈 Set this manually for test
  PropertiesService.getScriptProperties().setProperty(ZAPIER_TOKEN_NAME, expected);

  const actual = getExpectedZapierToken();
  logTestResult("getExpectedZapierToken", 1, "should retrieve stored token", actual === expected, actual, expected);
}

/**
 * Test: extractTokenFromRequest should return token if present
 */
function test_extractTokenFromRequest_valid() {
  const e = { parameter: { token: "secret123" } };
  const actual = extractTokenFromRequest(e);
  const expected = "secret123";

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
 * Test: isValidZapierToken should return true for matching token
 */
function test_isValidZapierToken_valid() {
  PropertiesService.getScriptProperties().setProperty(ZAPIER_TOKEN_NAME, "secret123");
  const actual = isValidZapierToken("secret123");

  logTestResult("isValidZapierToken", 1, "should return true for correct token", actual === true, actual, true);
}

/**
 * Test: isValidZapierToken should return false for incorrect token
 */
function test_isValidZapierToken_invalid() {
  PropertiesService.getScriptProperties().setProperty(ZAPIER_TOKEN_NAME, "secret123");
  const actual = isValidZapierToken("wrong-token");

  logTestResult("isValidZapierToken", 2, "should return false for wrong token", actual === false, actual, false);
}

/**
 * Test: isValidZapierToken should return false if token is empty
 */
function test_isValidZapierToken_empty() {
  PropertiesService.getScriptProperties().setProperty(ZAPIER_TOKEN_NAME, "secret123");
  const actual = isValidZapierToken("");
  const expected = false;

  logTestResult("isValidZapierToken", 3, "should reject empty token", actual === expected, actual, expected);
}