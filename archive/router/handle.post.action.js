// Ленивый реестр: actionId → функция
function getPostActionsRegistry_() {
  if (this.__POST_ACTIONS_REGISTRY) return this.__POST_ACTIONS_REGISTRY;

  var A = getPostActions_();
  this.__POST_ACTIONS_REGISTRY = {};

  // Примеры регистраций:
  this.__POST_ACTIONS_REGISTRY[A.HEALTHCHECK_V1] = handleHealthCheckBasicV1;
  this.__POST_ACTIONS_REGISTRY[A.GET_LAT_LNG_V1] = handleGetLatLngV1;
  this.__POST_ACTIONS_REGISTRY[A.STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1] = handleStripeDepositGetV1;
  this.__POST_ACTIONS_REGISTRY[A.CALCULATE_ESTIMATE_V1] = handleCalculateEstimateV1;

  return this.__POST_ACTIONS_REGISTRY;
}
