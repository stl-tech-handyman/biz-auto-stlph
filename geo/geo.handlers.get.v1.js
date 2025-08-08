function handleGetLatLngV1(e) {
    const address = e.parameter?.address;
    if (!address) return json({ error: "Missing 'address' parameter" });
    const loc = GeoService.getLatLng(address);
    return json({ ...loc, v: 1 });
}  