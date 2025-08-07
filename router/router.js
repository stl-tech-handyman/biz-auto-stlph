function doGet(e) {
  try {
    const action = getAction(e);                  // from utils.http
    const initiator = extractInitiator(e);        // from utils.http
    const token = extractTokenFromRequest(e);     // from utils.auth
    const version = getApiVersion(e);             // from utils.http

    const authCheck = isValidApiToken(token);

    console.log("📥 Incoming GET request");
    console.log("• Initiator:", initiator);
    console.log("• Action:", action);
    console.log("• API Version:", version);
    console.log("• Token Preview:", authCheck.preview || "none");
    console.log("• Token Source:", authCheck.source || "Unknown");

    if (!authCheck.valid) {
      console.warn("❌ Invalid API token attempt");
      return text("Unauthorized");
    }

    switch (version) {
      case 2:
        //return handleGetV2(action, e, initiator);
      case 1:
      default:
        return handleGetV1(action, e, initiator);
    }

  } catch (err) {
    console.error("🚨 Error in doGet:", err);
    return json({ error: err.message });
  }
}
