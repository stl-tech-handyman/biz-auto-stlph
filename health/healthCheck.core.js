/**
 * Returns a simple health check string for availability monitoring.
 * @returns {ContentService.Output}
 */
function healthCheckBasic() {
  return text(TEST_HEALTHCHECK_RESPONSE);
}