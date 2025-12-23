function test_getLatLng_validAddress() {
    const result = getLatLng("1600 Amphitheatre Parkway, Mountain View, CA");
    const passed = result && result.lat && result.lng && result.fullAddress;
    logTestResult("getLatLng", 1, "should return lat/lng for valid address", passed, result, "Non-empty geocode result");
  }
  
function test_getLatLng_emptyAddress() {
    try {
      getLatLng("");
      logTestResult("getLatLng", 2, "should fail on empty address", false, "no error", "Expected error");
    } catch (err) {
      logTestResult("getLatLng", 2, "should fail on empty address", true, err.message, "Error expected");
    }
}
  
function test_getLatLng_nullAddress() {
    try {
      getLatLng(null);
      logTestResult("getLatLng", 3, "should fail on null address", false, "no error", "Expected error");
    } catch (err) {
      logTestResult("getLatLng", 3, "should fail on null address", true, err.message, "Error expected");
    }
}
  
function test_getLatLng_noApiKey() {
    // Temporarily clear key
    const originalKey = PropertiesService.getScriptProperties().getProperty(GOOGLE_MAPS_API_KEY_VAR_NAME);
    PropertiesService.getScriptProperties().deleteProperty(GOOGLE_MAPS_API_KEY_VAR_NAME);
  
    try {
      getLatLng("123 Main St");
      logTestResult("getLatLng", 4, "should fail without API key", false, "no error", "Expected error");
    } catch (err) {
      logTestResult("getLatLng", 4, "should fail without API key", true, err.message, "Error expected");
    } finally {
      // Restore key
      if (originalKey) {
        PropertiesService.getScriptProperties().setProperty(GOOGLE_MAPS_API_KEY_VAR_NAME, originalKey);
      }
    }
}