// Постоянная таблица внутренних действий (ленивая инициализация)
function getPostActions_() {
  if (this.__POST_ACTIONS) return this.__POST_ACTIONS;

  this.__POST_ACTIONS = {
    // --- Health ---
    HEALTHCHECK_V1: "HEALTHCHECK_V1",
    HEALTHCHECK:    "HEALTHCHECK_V1", // алиас на актуальную

    // --- Geo ---
    GET_LAT_LNG_V1: "GET_LAT_LNG_V1",
    GET_LAT_LNG:    "GET_LAT_LNG_V1", // алиас

    // --- Stripe ---
    STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1: "STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1",

    // --- Estimates ---
    CALCULATE_ESTIMATE_V1: "CALCULATE_ESTIMATE_V1",
  };
  return this.__POST_ACTIONS;
}

// Публичные имена → внутренние ID (лениво)
function getPublicPostActions_() {
  if (this.__PUBLIC_POST_ACTIONS) return this.__PUBLIC_POST_ACTIONS;

  var A = getPostActions_();
  this.__PUBLIC_POST_ACTIONS = {
    // Health
    HEALTHCHECK_V1: A.HEALTHCHECK_V1,
    HEALTHCHECK:    A.HEALTHCHECK,

    // Geo
    GET_LAT_LNG_V1: A.GET_LAT_LNG_V1,
    GET_LAT_LNG:    A.GET_LAT_LNG,

    // Stripe
    STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1: A.STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1,

    // Estimates
    CALCULATE_ESTIMATE_V1: A.CALCULATE_ESTIMATE_V1,
  };
  return this.__PUBLIC_POST_ACTIONS;
}
