function handleGetV1(action, e, initiator) {
  switch (action) {

    case GetActions.HEALTHCHECK:
      return handleHealthCheckBasicV1();

    case GetActions.GET_LAT_LNG:
      return handleGetLatLngV1(e);

    /*
    case TestGetActions.TEST_EMAIL_SENDING:
      return handleSendEmailTestV1(e);

    case TestGetActions.TEST_EMAIL_SENDING_QUOTE_EMAIL:
      return handleSendQuoteEmailTestV1(e);

    case PublicGetActions.STRIPE_GET_BOOKING_DEPOSIT_AMOUNT:
      return handleStripeDepositGetV1(e);

    case PublicGetActions.CALCULATE_ESTIMATE:
      return handleCalculateEstimateV1(e);
    */
    
    default:
      throw new Error("Unknown GET action: " + action);
  }
}