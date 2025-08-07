function test_doGet_validToken_validAction_v1() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      v: "1",
      "req-initiator": "zapier",
      api_token: VALID_TEST_TOKEN
    }
  };
  const result = doGet(e).getContent();
  const expected = "TESTING ROUTE ALIVE";
  logTestResult("doGet", 1, "should return alive response for v1", result === expected, result, expected);
}

function test_doGet_validToken_validAction_v2() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK_V1", // Or any v2 route you test
      v: "2",
      "req-initiator": "zapier",
      api_token: VALID_TEST_TOKEN
    }
  };
  const result = doGet(e).getContent();
  const expected = JSON.stringify(healthCheckResponse()); // Or stub it
  logTestResult("doGet", 2, "should return healthcheck v2 response", result === expected, result, expected);
}

function test_doGet_invalidToken() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      "req-initiator": "zapier",
      api_token: "WRONG_TOKEN"
    }
  };

  const result = doGet(e).getContent();
  const expected = "Unauthorized";
  logTestResult("doGet", 2, "should reject bad token", result === expected, result, expected);
}

function test_doGet_missingToken() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      "req-initiator": "zapier"
      // No token
    }
  };

  const result = doGet(e).getContent();
  const expected = "Unauthorized";
  logTestResult("doGet", 4, "should reject missing token", result === expected, result, expected);
}

function test_doGet_missingAction() {
  const e = {
    parameter: {
      api_token: VALID_TEST_TOKEN,
      "req-initiator": "zapier"
    }
  };
  const result = doGet(e).getContent();
  logTestResult("doGet", 5, "should fail on missing action", result.includes("Missing 'action'"), result, "Expected error");
}

function test_doGet_unknownAction() {
  const e = {
    parameter: {
      action: "UNKNOWN_ACTION",
      v: "1",
      api_token: VALID_TEST_TOKEN,
      "req-initiator": "zapier"
    }
  };
  const parsed = JSON.parse(doGet(e).getContent());
  logTestResult("doGet", 6, "should return error for unknown action", parsed.error.includes("Unknown GET action"), parsed.error, "Expected error message");
}

function test_doGet_floatVersionFallbackToV1() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      v: "1.5",
      api_token: VALID_TEST_TOKEN,
      "req-initiator": "zapier"
    }
  };
  const result = doGet(e).getContent();
  const expected = "TESTING ROUTE ALIVE";
  logTestResult("doGet", 7, "should fallback to v1 for float version", result === expected, result, expected);
}

function test_doGet_invalidVersion_fallback() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      v: "banana",
      api_token: VALID_TEST_TOKEN,
      "req-initiator": "zapier"
    }
  };
  const result = doGet(e).getContent();
  const expected = "TESTING ROUTE ALIVE";
  logTestResult("doGet", 8, "should fallback to v1 for invalid version string", result === expected, result, expected);
}

function test_doGet_missingVersion_defaultsToV1() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      api_token: VALID_TEST_TOKEN,
      "req-initiator": "zapier"
    }
  };
  const result = doGet(e).getContent();
  const expected = "TESTING ROUTE ALIVE";
  logTestResult("doGet", 9, "should default to v1 when version is missing", result === expected, result, expected);
}


function test_doGet_errorInHandler() {
  const e = {
    parameter: {
      action: "UNKNOWN_ACTION",
      "req-initiator": "zapier",
      api_token: VALID_TEST_TOKEN
    }
  };

  const response = doGet(e);
  const content = JSON.parse(response.getContent());
  const hasError = content?.error?.includes("Unknown GET action");

  logTestResult("doGet", 5, "should catch thrown errors", hasError, content.error, "should contain 'Unknown GET action'");
}