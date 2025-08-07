/**
 * Unit Test: Tests the raw return value from healthCheckBasic
 */
function test_healthCheckBasic_unit() {
  const response = healthCheckBasic();
  const actual = response?.getContent();
  const expected = TEST_HEALTHCHECK_RESPONSE;

  logTestResult("healthCheckBasic", 1, "should return basic healthcheck string", actual === expected, actual, expected);
}
