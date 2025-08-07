/**
 * Main test runner for all util tests.
 */
function test_all() {
  console.log("🚀 Running all unit tests...");

  test_getApiVersion_valid();
  test_getApiVersion_missing();
  test_getApiVersion_invalidString();
  test_getApiVersion_negative();
  test_getApiVersion_float();
  test_getApiVersion_nullEvent();

  test_getAction();
  test_getVersion();
  test_jsonResponseHelper();
  test_textResponseHelper();
  test_textResponseHelper();

  console.log("✅ All tests executed.");
}

function test_getApiVersion_valid() {
  const e = { parameter: { v: "2" } };
  const actual = getApiVersion(e);
  const expected = 2;

  logTestResult("getApiVersion", 1, "should return valid version", actual === expected, actual, expected);
}

function test_getApiVersion_missing() {
  const e = { parameter: {} };
  const actual = getApiVersion(e);
  const expected = 1;

  logTestResult("getApiVersion", 2, "should default to version 1", actual === expected, actual, expected);
}

function test_getApiVersion_invalidString() {
  const e = { parameter: { v: "banana" } };
  const actual = getApiVersion(e);
  const expected = 1;

  logTestResult("getApiVersion", 3, "should fallback for invalid string", actual === expected, actual, expected);
}

function test_getApiVersion_negative() {
  const e = { parameter: { v: "-5" } };
  const actual = getApiVersion(e);
  const expected = 1;

  logTestResult("getApiVersion", 4, "should fallback for negative version", actual === expected, actual, expected);
}

function test_getApiVersion_float() {
  const e = { parameter: { v: "2.7" } };
  const actual = getApiVersion(e);
  const expected = 1;

  logTestResult("getApiVersion", 5, "should fallback for float input", actual === expected, actual, expected);
}

function test_getApiVersion_nullEvent() {
  const actual = getApiVersion(null);
  const expected = 1;

  logTestResult("getApiVersion", 6, "should fallback when event is null", actual === expected, actual, expected);
}

/**
 * Unit test for getAction()
 */
function test_getAction() {
  const tests = [
    {
      description: "✅ should return action if present",
      input: { parameter: { action: "sendQuote" } },
      expected: "sendQuote"
    },
    {
      description: "❌ should throw if action is missing",
      input: { parameter: {} },
      expectError: true
    },
    {
      description: "❌ should throw if parameter is missing entirely",
      input: {},
      expectError: true
    }
  ];

  runParamTests(tests, getAction, "getAction");
}

/**
 * Unit test for getVersion()
 */
function test_getVersion() {
  const tests = [
    {
      description: "✅ should return version if present",
      input: { parameter: { version: "v2" } },
      expected: "v2"
    },
    {
      description: "✅ should default to 'v1' if version is missing",
      input: { parameter: {} },
      expected: "v1"
    },
    {
      description: "✅ should default to 'v1' if parameter is missing",
      input: {},
      expected: "v1"
    }
  ];

  runParamTests(tests, getVersion, "getVersion");
}

/**
 * Unit test for json() response helper.
 */
function test_jsonResponseHelper() {
  const input = { msg: "OK", code: 200 };
  const output = json(input);

  const content = output.getContent();
  const mimeType = output.getMimeType();

  logTestResult("json()", 1, "returns valid JSON", content === JSON.stringify(input) && mimeType === ContentService.MimeType.JSON, content, JSON.stringify(input));
}

/**
 * Unit test for text() response helper.
 */
function test_textResponseHelper() {
  const message = "Hello World";
  const output = text(message);

  const content = output.getContent();
  const mimeType = output.getMimeType();

  logTestResult("text()", 1, "returns valid plain text", content === message && mimeType === ContentService.MimeType.TEXT, content, message);
}

function test_extractInitiator() {
  const tests = [
    {
      description: "✅ should extract initiator if present",
      input: { parameter: { "req-initiator": "ZapierBot" } },
      expected: "ZapierBot"
    },
    {
      description: "❌ should return 'Unknown' if req-initiator is missing",
      input: { parameter: {} },
      expected: "Unknown"
    },
    {
      description: "❌ should return 'Unknown' if parameter is missing",
      input: {},
      expected: "Unknown"
    },
    {
      description: "❌ should return 'Unknown' if e is null",
      input: null,
      expected: "Unknown"
    },
    {
      description: "❌ should return 'Unknown' if e is undefined",
      input: undefined,
      expected: "Unknown"
    },
  ];

  runParamTests(tests, extractInitiator);
}

function runParamTests(tests, fn) {
  const fnName = fn.name || "anonymous";
  let passed = 0;

  tests.forEach((t, i) => {
    try {
      const result = fn(t.input);
      const success = t.expectError ? false : result === t.expected;
      logTestResult(fnName, i + 1, t.description, success, result, t.expected);
      if (success) passed++;
    } catch (err) {
      const success = t.expectError;
      logTestResult(fnName, i + 1, t.description, success, err.message, "Expected error");
      if (success) passed++;
    }
  });

  console.log(`🎯 ${fnName}: ${passed}/${tests.length} tests passed\n`);
}