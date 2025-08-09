/**
 * Tests for handleGetLatLngV1(bodyOrEvent, context)
 *
 * Covers:
 *  1) POST body: success
 *  2) POST body: trims whitespace
 *  3) POST body: missing/empty address → 400
 *  4) GET event: success (e.parameter.address)
 *  5) GET event: missing address → 400
 *  6) getLatLng throws → handler rethrows (so router can catch)
 *
 * How to run:
 *  - Open Apps Script IDE → Run → test_handleGetLatLngV1_all
 */

//////////////////// Minimal test helpers ////////////////////

function _t_log(ok, name, extra) {
    const mark = ok ? "✅" : "❌";
    const line = `${mark} ${name}${extra ? " — " + extra : ""}`;
    console.log(line);
    if (!ok) throw new Error(line);
  }
  
  function _t_expectJson(out) {
    // handler returns a ContentService.TextOutput via jsonOk/jsonError
    if (!out || typeof out.getContent !== "function") {
      throw new Error("Output is not a TextOutput");
    }
    const txt = out.getContent();
    try { return JSON.parse(txt); }
    catch (e) { throw new Error("Response is not valid JSON: " + txt); }
  }
  
  function _t_makeGetEvent(address) {
    return {
      parameter: address != null ? { address: address } : {},
      parameters: address != null ? { address: [String(address)] } : {},
      postData: null
    };
  }
  
  // Keep originals to restore after stubs
  var __orig = {};
  function _t_stub(name, fn) { __orig[name] = this[name]; this[name] = fn; }
  function _t_restore(name)   { if (name in __orig) this[name] = __orig[name]; }
  function _t_restoreAll()    { Object.keys(__orig).forEach(_t_restore); __orig = {}; }
  
  //////////////////// Local stand-ins for jsonOk/jsonError ////////////////////
  // (We stub the project’s versions to keep tests hermetic.)
  
  function _t_jsonOk(data, meta) {
    return ContentService.createTextOutput(
      JSON.stringify({ ok: true, data: data || null, meta: meta || undefined })
    ).setMimeType(ContentService.MimeType.JSON);
  }
  
  function _t_jsonError(message, code, meta) {
    const msg = (message && message.message) || String(message || "Unknown error");
    const body = { ok: false, error: { message: msg } };
    if (code !== undefined) body.error.code = code;
    if (meta) body.meta = meta;
    return ContentService.createTextOutput(JSON.stringify(body))
      .setMimeType(ContentService.MimeType.JSON);
  }
  
  //////////////////// Tests ////////////////////
  
  function test_handleGetLatLngV1_post_success() {
    // Arrange
    const body = { address: "1600 Amphitheatre Parkway, Mountain View, CA" };
    const ctx = { requestId: "r1" };
  
    // Stub jsonOk/jsonError/getLatLng
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
    _t_stub("getLatLng", function(addr) {
      // verify the address passed through unchanged here
      _t_log(addr === body.address, "POST success: address passed to getLatLng");
      return { lat: 37.422, lng: -122.084, fullAddress: "Googleplex, Mountain View, CA" };
    });
  
    // Act
    const out = handleGetLatLngV1(body, ctx);
    const json = _t_expectJson(out);
  
    // Assert
    _t_log(json.ok === true, "POST success: ok flag");
    _t_log(json.data && json.data.lat === 37.422, "POST success: lat");
    _t_log(json.data && json.data.lng === -122.084, "POST success: lng");
    _t_log(json.data && json.data.fullAddress === "Googleplex, Mountain View, CA", "POST success: fullAddress");
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_post_trims() {
    // Arrange
    const body = { address: "   1600 Amphitheatre Pkwy   " };
    const trimmed = "1600 Amphitheatre Pkwy";
    const ctx = {};
  
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
  
    let received;
    _t_stub("getLatLng", function(addr) {
      received = addr;
      return { lat: 1, lng: 2, fullAddress: "X" };
    });
  
    // Act
    const out = handleGetLatLngV1(body, ctx);
    const json = _t_expectJson(out);
  
    // Assert
    _t_log(json.ok === true, "POST trims: ok");
    _t_log(received === trimmed, "POST trims: handler trimmed before calling getLatLng", "got=" + received);
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_post_missing_address() {
    // Arrange
    const body = { address: "   " }; // empty after trim
    const ctx = {};
  
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
    _t_stub("getLatLng", function(){ throw new Error("should not be called"); });
  
    // Act
    const out = handleGetLatLngV1(body, ctx);
    const json = _t_expectJson(out);
  
    // Assert
    _t_log(json.ok === false, "POST missing: ok=false");
    _t_log(json.error && json.error.code === 400, "POST missing: 400 code");
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_get_success() {
    // Arrange
    const e = _t_makeGetEvent("10 Downing St, London");
    const ctx = {};
  
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
    _t_stub("getLatLng", function(addr) {
      _t_log(addr === "10 Downing St, London", "GET success: param picked from e.parameter.address");
      return { lat: 51.5034, lng: -0.1276, fullAddress: "10 Downing St, Westminster, London" };
    });
  
    // Act
    const out = handleGetLatLngV1(e, ctx);
    const json = _t_expectJson(out);
  
    // Assert
    _t_log(json.ok === true, "GET success: ok");
    _t_log(json.data && json.data.lat === 51.5034, "GET success: lat");
    _t_log(json.data && json.data.lng === -0.1276, "GET success: lng");
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_get_missing_address() {
    // Arrange
    const e = _t_makeGetEvent(undefined); // no address
    const ctx = {};
  
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
    _t_stub("getLatLng", function(){ throw new Error("should not be called"); });
  
    // Act
    const out = handleGetLatLngV1(e, ctx);
    const json = _t_expectJson(out);
  
    // Assert
    _t_log(json.ok === false, "GET missing: ok=false");
    _t_log(json.error && json.error.code === 400, "GET missing: 400 code");
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_getLatLng_throws() {
    // Arrange
    const body = { address: "Some Place" };
    const ctx = {};
  
    _t_stub("jsonOk", _t_jsonOk);
    _t_stub("jsonError", _t_jsonError);
    _t_stub("getLatLng", function(){ throw new Error("Geocoding failed: OVER_QUERY_LIMIT"); });
  
    // Act + Assert: since handler doesn't catch, it should throw
    var threw = false;
    try {
      handleGetLatLngV1(body, ctx);
    } catch (e) {
      threw = /OVER_QUERY_LIMIT/.test(String(e && e.message || e));
    }
    _t_log(threw, "getLatLng throws: handler rethrows for router to catch");
  
    _t_restoreAll();
  }
  
  function test_handleGetLatLngV1_all() {
    try {
      test_handleGetLatLngV1_post_success();
      test_handleGetLatLngV1_post_trims();
      test_handleGetLatLngV1_post_missing_address();
      test_handleGetLatLngV1_get_success();
      test_handleGetLatLngV1_get_missing_address();
      test_handleGetLatLngV1_getLatLng_throws();
      console.log("🎉 All handleGetLatLngV1 tests passed");
    } catch (e) {
      console.error("❌ handleGetLatLngV1 tests failed:", e && e.message || e);
      throw e;
    }
  }
  