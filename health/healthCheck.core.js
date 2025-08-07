/**
 * Returns a simple health check string for availability monitoring.
 * @returns {ContentService.Output}
 */
function handleHealthCheckBasicV1() {
  return text(TEST_HEALTHCHECK_RESPONSE);
}