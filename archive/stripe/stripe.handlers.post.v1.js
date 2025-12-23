/**
 * Handler for STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1
 * Returns the deposit amount for a given booking value.
 * 
 * @param {Object} bodyOrEvent - Request body (POST) or event object (GET)
 * @param {Object} context - Context object with actionId, requestId, initiator, etc.
 * @returns {ContentService.Output} JSON response
 */
function handleStripeDepositGetV1(bodyOrEvent, context) {
  // For POST requests, amount comes from body; for GET requests, from e.parameter
  const amountStr = (bodyOrEvent.amount || bodyOrEvent.parameter?.amount || "").trim();
  
  if (!amountStr) {
    return jsonError("Missing 'amount' parameter", 400);
  }

  const amount = Number(amountStr);
  if (isNaN(amount) || amount <= 0) {
    return jsonError("'amount' must be a positive number", 400);
  }

  // Find the closest matching price ID from BOOKING_DEPOSITS_PRICE_IDS
  // We'll find the price ID that matches the amount, or the closest one below
  let selectedPrice = null;
  let closestDiff = Infinity;

  for (let i = 0; i < BOOKING_DEPOSITS_PRICE_IDS.length; i++) {
    const price = BOOKING_DEPOSITS_PRICE_IDS[i];
    if (price.id === "N/A") continue;
    
    const diff = Math.abs(price.value - amount);
    if (diff < closestDiff) {
      closestDiff = diff;
      selectedPrice = price;
    }
  }

  if (!selectedPrice || selectedPrice.id === "N/A") {
    return jsonError("No matching deposit price found for amount: " + amount, 404);
  }

  return jsonOk({
    amount: selectedPrice.value,
    priceId: selectedPrice.id,
    requestedAmount: amount
  }, {
    actionId: context?.actionId || "STRIPE_GET_BOOKING_DEPOSIT_AMOUNT_V1",
    requestId: context?.requestId,
    initiator: context?.initiator
  });
}

