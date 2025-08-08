/**
 * Action constants for GET requests.
 * Used by the router to match incoming ?action=... to the correct handler.
 */
const GetActions = {
    HEALTHCHECK: "HEALTHCHECK",
    GET_LAT_LNG: "GET_LAT_LNG",
  
    // Test-only actions
    TEST_EMAIL_SENDING: "TEST_EMAIL_SENDING",
    TEST_EMAIL_SENDING_QUOTE_EMAIL: "TEST_EMAIL_SENDING_QUOTE_EMAIL",
  
    // Stripe
    STRIPE_GET_BOOKING_DEPOSIT_AMOUNT: "STRIPE_GET_BOOKING_DEPOSIT_AMOUNT",
  
    // Estimates
    CALCULATE_ESTIMATE: "CALCULATE_ESTIMATE",
  };
  
  /**
   * Public-facing GET actions (external API surface).
   * These may map directly to internal GetActions or be aliases.
   */
  const PublicGetActions = {
    HEALTHCHECK_V1: GetActions.HEALTHCHECK,
    GET_LAT_LNG_V1: GetActions.GET_LAT_LNG,
  
    // Optional alias for convenience
    HEALTHCHECK: GetActions.HEALTHCHECK,
  };
  