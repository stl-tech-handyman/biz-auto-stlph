function handleGetLatLngV1(bodyOrEvent, context) {
    // For POST requests, address comes from body; for GET requests, from e.parameter
    const address = (bodyOrEvent.address || bodyOrEvent.parameter?.address || "").trim();
    if (!address) return jsonError("Missing 'address' parameter", 400);

    const loc = getLatLng(address);
    return jsonOk({ ...loc});
}  