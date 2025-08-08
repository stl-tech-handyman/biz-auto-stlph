function doGet(e) {
  try {
    const action    = getAction(e);            // utils/http
    const initiator = extractInitiator(e);     // utils/http
    const version   = getApiVersion(e);        // utils/http

    // Auth
    const { token, source } = extractTokenFromRequest(e);      // utils/auth
    const authCheck = isValidApiToken(token);                  // utils/auth

    // Request log (safe previews only)
    console.log("📥 Incoming GET request");
    console.log("• Initiator:", initiator);
    console.log("• Action:", action);
    console.log("• API Version:", version);
    console.log("• Token Preview:", authCheck.preview || "none");
    console.log("• Token Source:", source || "Unknown");

    if (!authCheck.valid) {
      console.warn("❌ Unauthorized request (invalid API token)");
      return jsonError("Unauthorized", 401, { action, v: version });
    }

    switch (version) {
      case 2:
        // return handleGetV2(action, e, initiator);
        // break; // keep commented until v2 exists
      case 1:
      default:
        return handleGetV1(action, e, initiator);
    }
  } catch (err) {
    const id = Utilities.getUuid().slice(0, 8);
    console.error(`🚨 Error in doGet [${id}]:`, err && err.stack || err);
    return jsonError(err, "GET_HANDLER_ERROR", { errorId: id });
  }
}
