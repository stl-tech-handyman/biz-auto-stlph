// Domain used for email sender construction
const DOMAIN_NAME = "stlpartyhelpers.com";

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
