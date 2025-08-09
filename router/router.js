function doPost(e) {
  var reqId = newRequestId_();
  var t0 = spanStart_(); // start timing

  try {
    if (LOGGING_CONFIG.enableRequestSnapshot) logRequestSummary_(e, reqId);

    var body = {};
    try {
      body = JSON.parse(e.postData && e.postData.contents || "{}");
    } catch (parseErr) {
      logWarn("Invalid JSON body", {}, reqId);
      return jsonError("Invalid JSON body", 400, { requestId: reqId });
    }

    var token = body.token || null;
    var auth = isValidApiToken(token);
    logInfo("Auth check", { tokenPreview: auth.preview, valid: auth.valid, source: "body.token" }, reqId);
    if (!auth.valid) {
      logWarn("Unauthorized", {}, reqId);
      return jsonError("Unauthorized", 401, { requestId: reqId });
    }

    var publicAction = body.action;
    if (!publicAction) {
      return jsonError("Missing required parameter: 'action'", 400, { requestId: reqId });
    }

    var actionId = (getPublicPostActions_()[publicAction] || publicAction);
    var fn = getPostActionsRegistry_()[actionId];
    if (!fn) {
      logWarn("Unknown action", { actionPublic: publicAction, actionId: actionId }, reqId);
      return jsonError("Unknown action", 404, { action: actionId, requestId: reqId });
    }

    var initiator = body.initiator || "Unknown";
    var res = fn(body, { initiator: initiator, actionId: actionId, publicAction: publicAction, requestId: reqId });

    logInfo("POST completed", {
      actionId: actionId,
      initiator: initiator,
      status: "ok",
      durationMs: spanEndMs_(t0) // end timing
    }, reqId);

    return res;

  } catch (err) {
    logError("POST failed", {
      message: (err && err.message) || String(err),
      stack: err && err.stack,
      durationMs: spanEndMs_(t0) // end timing
    }, reqId);
    return jsonError(err, "POST_HANDLER_ERROR", { requestId: reqId });
  }
}
