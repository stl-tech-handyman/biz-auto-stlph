function handleHealthCheckBasicV1(e, ctx) {
  return jsonOk(
    { status: "ok" },
    { actionId: ctx?.actionId || "HEALTHCHECK_V1", requestId: ctx?.requestId, initiator: ctx?.initiator }
  );
}
