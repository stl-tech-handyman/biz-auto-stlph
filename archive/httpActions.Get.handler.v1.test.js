function test_handleGetV1_healthcheck() {
  const response = handleGetV1(TestGetActions.TEST_HEALTHCHECK, {}, "tester");
  const actual = response?.getContent();
  const expected = TEST_HEALTHCHECK_RESPONSE;

  logTestResult("handleGetV1", 1, "should return healthcheck from route", actual === expected, actual, expected);
}