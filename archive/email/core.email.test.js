function test_sendEmail_valid_withDefaultSubject() {
  const result = sendEmail(
    "test-sre@stlpartyhelpers.com",
    DEFAULT_TEST_RECIPIENT,
    "<p>Hello World</p>"
  );

  logTestResult("sendEmail", 1, "should send with default subject", result.status === "sent", result.status, "sent");
}

function test_sendEmail_valid_withCustomSubject() {
  const customSubject = "Custom Subject";
  const result = sendEmail(
    "test-sre@stlpartyhelpers.com",
    DEFAULT_TEST_RECIPIENT,
    "<p>Hello again</p>",
    customSubject
  );

  logTestResult("sendEmail", 2, "should send with custom subject", result.subject === customSubject, result.subject, customSubject);
}

function test_sendEmail_rejectsUnauthorizedSender() {
  try {
    sendEmail("hacker@notallowed.com", DEFAULT_TEST_RECIPIENT, "<p>Oops</p>");
    logTestResult("sendEmail", 3, "should reject unauthorized sender", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 3, "should reject unauthorized sender", true, err.message, "Error expected");
  }
}

function test_sendEmail_rejectsEmptyRecipient() {
  try {
    sendEmail("test-sre@stlpartyhelpers.com", "", "<p>Oops</p>");
    logTestResult("sendEmail", 4, "should reject empty recipient", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 4, "should reject empty recipient", true, err.message, "Error expected");
  }
}

function test_sendEmail_rejectsBadEmailFormat() {
  try {
    sendEmail("test-sre@stlpartyhelpers.com", "invalidemail", "<p>Oops</p>");
    logTestResult("sendEmail", 5, "should reject badly formatted recipient", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 5, "should reject badly formatted recipient", true, err.message, "Error expected");
  }
}

function test_sendEmail_rejectsEmptyHtmlBody() {
  try {
    sendEmail("test-sre@stlpartyhelpers.com", DEFAULT_TEST_RECIPIENT, "");
    logTestResult("sendEmail", 6, "should reject empty htmlBody", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 6, "should reject empty htmlBody", true, err.message, "Error expected");
  }
}

function test_sendEmail_rejectsWhitespaceOnlyHtmlBody() {
  try {
    sendEmail("test-sre@stlpartyhelpers.com", DEFAULT_TEST_RECIPIENT, "    ");
    logTestResult("sendEmail", 7, "should reject whitespace-only htmlBody", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 7, "should reject whitespace-only htmlBody", true, err.message, "Error expected");
  }
}

function test_sendEmail_rejectsMissingHtmlBody() {
  try {
    sendEmail("test-sre@stlpartyhelpers.com", DEFAULT_TEST_RECIPIENT);
    logTestResult("sendEmail", 8, "should reject missing htmlBody", false, "no error", "Expected error");
  } catch (err) {
    logTestResult("sendEmail", 8, "should reject missing htmlBody", true, err.message, "Error expected");
  }
}
