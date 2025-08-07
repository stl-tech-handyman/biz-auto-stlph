function test_all() {
  console.log("🔁 Running doGet() tests...");

  test_doGet_validToken_validAction();
  test_doGet_invalidToken();
  test_doGet_missingAction();
  test_doGet_missingToken();
  test_doGet_errorInHandler();
  
  console.log("✅ All tests completed.");
}


function test_doGet_validToken_validAction() {
  const e = {
    parameter: {
      action: "TEST_HEALTHCHECK",
      "req-initiator": "zapier",
      api_token: VALID_TEST_TOKEN  // set via Properties
    }
  };

  const result = doGet(e).getContent();
  const expected = "TESTING ROUTE ALIVE";
  logTestResult("doGet", 1, "should return alive response", result === expected, result, expected);
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

function test_doGet_missingAction() {
  const e = {
    parameter: {
      "req-initiator": "zapier",
      api_token: VALID_TEST_TOKEN
    }
  };

  const result = doGet(e).getContent();
  const expected = "Missing 'action' parameter";
  logTestResult("doGet", 3, "should reject missing action param", result === expected, result, expected);
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