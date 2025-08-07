// Domain used for email sender construction
const DOMAIN_NAME = "stlpartyhelpers.com";

// STL Tech Handyman Slack
const DEFAULT_TEST_RECIPIENT = "test-stlph-app-script-aaaaq6jpm4tpkl7764inkz7hpu@stlph.slack.com";


/**
 * Allowed sender identities and associated metadata.
 */
const ALLOWED_SENDERS = {
  [`test-sre@${DOMAIN_NAME}`]: {
    name: "Site Reliability Team",
    prefix: "sre"
  },
  [`team@${DOMAIN_NAME}`]: {
    name: "STL Party Helpers Team",
    prefix: "ops"
  },
  [`team@${DOMAIN_NAME}`]: {
    name: "STL Party Helpers Sales",
    prefix: ""
  }
};

/**
 * Default test sender email, used in testing scenarios.
 */
const TEST_SRE_EMAIL = `test-sre@${DOMAIN_NAME}`;

/**
 * Static response for basic ping/healthcheck endpoint.
 */
const TEST_HEALTHCHECK_RESPONSE = "PONG";

/**
 * Prefix map for different test actions.
 */
const TestGetActions = {
  TEST_HEALTHCHECK: "TEST_HEALTHCHECK",
  TEST_HEALTHCHECK_V1: "TEST_HEALTHCHECK_V1",
  TEST_EMAIL_SENDING: "TEST_EMAIL_SENDING",
  TEST_EMAIL_SENDING_QUOTE_EMAIL: "TEST_EMAIL_SENDING_QUOTE_EMAIL",
  TEST_EMAIL_QUOTESENDING_WITHFORWARDING: "TEST_EMAIL_QUOTESENDING_WITHFORWARDING"
};
