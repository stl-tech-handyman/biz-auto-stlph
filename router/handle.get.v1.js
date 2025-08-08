function handleGetV1(action, e, initiator) {
  switch (action) {
    case GetActions.HEALTHCHECK:          return handleHealthCheckBasicV1();
    case GetActions.GET_LAT_LNG:          return handleGetLatLngV1(e);
    case GetActions.CALCULATE_ESTIMATE:   return handleCalculateEstimateV1(e);
    default: throw new Error("Unknown GET action: " + action);
  }
}