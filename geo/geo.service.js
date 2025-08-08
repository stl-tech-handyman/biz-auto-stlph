/**
 * Returns the latitude, longitude, and formatted address for a given address.
 * Uses Google Maps Geocoding API.
 *
 * @param {string} address - The address to geocode (required)
 * @param {Object} [options] - Optional settings:
 *    { region, language, components, cacheTtlSec }
 * @returns {{ lat: number, lng: number, fullAddress: string }}
 * @throws {Error} If geocoding fails or API key is missing
 */
function getLatLng(address, options) {
    // --- Validation ---
    if (!address || typeof address !== "string" || address.trim().length === 0) {
      throw new Error("Address is required and must be a non-empty string.");
    }
  
    // --- Config from script properties ---
    const apiKey = getScriptProp(CONFIG_KEYS.GOOGLE_MAPS_API_KEY);
    const endpoint = getScriptProp(
      CONFIG_KEYS.GOOGLE_MAPS_GEOCODE_URL,
      { required: false, defaultValue: "https://maps.googleapis.com/maps/api/geocode/json" }
    );
  
    // --- Simple cache ---
    const cacheKey = `geocode:${address.trim().toLowerCase()}:${options?.region || ""}:${options?.language || ""}`;
    const cache = CacheService.getScriptCache();
    const cached = cache.get(cacheKey);
    if (cached) {
      try { return JSON.parse(cached); } catch (_) {}
    }
  
    // --- Build URL ---
    const qs = [
      "address=" + encodeURIComponent(address),
      "key=" + encodeURIComponent(apiKey)
    ];
    if (options?.region) qs.push("region=" + encodeURIComponent(options.region));
    if (options?.language) qs.push("language=" + encodeURIComponent(options.language));
    if (options?.components) qs.push("components=" + encodeURIComponent(options.components));
  
    const url = `${endpoint}?${qs.join("&")}`;
  
    // --- Fetch ---
    const resp = UrlFetchApp.fetch(url, { muteHttpExceptions: true });
    const code = resp.getResponseCode();
    const json = JSON.parse(resp.getContentText() || "{}");
  
    // --- Validate response ---
    if (code === 200 && json.status === "OK" && json.results?.length > 0) {
      const location = json.results[0].geometry.location;
      const payload = {
        lat: location.lat,
        lng: location.lng,
        fullAddress: json.results[0].formatted_address
      };
  
      // Cache for 1 hour by default (max 6h)
      const ttl = Math.min(options?.cacheTtlSec || 3600, 21600);
      try { cache.put(cacheKey, JSON.stringify(payload), ttl); } catch (_) {}
  
      return payload;
    }
  
    // --- Error handling ---
    const apiErr = json.error_message || json.status || `HTTP ${code}`;
    throw new Error(`Geocoding failed: ${apiErr}`);
  }
  