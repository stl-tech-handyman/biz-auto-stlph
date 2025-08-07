function handleGetV1(action, e, initiator) {
  switch (action) {
    case PublicGetActions.TEST_HEALTHCHECK:
      return handleHealthCheckBasicV1();

    case PublicGetActions.TEST_HEALTHCHECK_V1:
      return handleHealthCheckV1();

    case PublicGetActions.TEST_EMAIL_SENDING:
      return handleSendEmailTestV1(e);

    case PublicGetActions.TEST_EMAIL_SENDING_QUOTE_EMAIL:
      return handleSendQuoteEmailTestV1(e);

    case PublicGetActions.STRIPE_GET_BOOKING_DEPOSIT_AMOUNT:
      return handleStripeDepositGetV1(e);

    case PublicGetActions.CALCULATE_ESTIMATE:
      return handleCalculateEstimateV1(e);

    case PublicGetActions.GET_LAT_LNG:
      return handleGetLatLngV1(e);

    default:
      throw new Error("Unknown GET action: " + action);
  }
}