/**
 * Returns the latitude, longitude, and formatted address for a given address.
 * Uses Google Maps Geocoding API.
 *
 * @param {string} address - The address to geocode
 * @returns {{ lat: number, lng: number, fullAddress: string }}
 * @throws {Error} If geocoding fails or API key is missing
 */
function getLatLng(address) {
    if (!address || typeof address !== "string" || address.trim().length === 0) {
      throw new Error("Address is required and must be a non-empty string.");
    }
  
    const apiKey = PropertiesService.getScriptProperties().getProperty(GOOGLE_MAPS_API_KEY_VAR_NAME);
    if (!apiKey) {
      throw new Error("Google Maps API key is not set in script properties.");
    }
  
    const url = `https://maps.googleapis.com/maps/api/geocode/json?address=${encodeURIComponent(address)}&key=${apiKey}`;
    const response = UrlFetchApp.fetch(url);
    const json = JSON.parse(response.getContentText());
  
    if (json.status === "OK" && json.results?.length > 0) {
      const location = json.results[0].geometry.location;
      return {
        lat: location.lat,
        lng: location.lng,
        fullAddress: json.results[0].formatted_address
      };
    } else {
      const errorMessage = json.error_message || json.status || "Unknown error from Geocoding API";
      throw new Error(`Geocoding failed: ${errorMessage}`);
    }
  }
  