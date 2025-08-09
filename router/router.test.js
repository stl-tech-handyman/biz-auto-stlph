/**
 * Tests for router.doPost(e)
 * - No external libs; pure Apps Script-style tests.
 * - We stub ONLY the things doPost depends on.
 *
 * What we cover:
 *  1) invalid JSON body
 *  2) missing token
 *  3) invalid token
 *  4) missing action
 *  5) unknown action
 *  6) success path (HEALTHCHECK via alias and via explicit V1)
 *
 * How to run:
 *  - Open Apps Script IDE → Run → test_doPost_all
 */

//////////////////// Test Harness ////////////////////

function tp_logResult_(name, ok, details) {
  const mark = ok ? "✅" : "❌";
  const line = `${mark} ${name}${details ? " — " + details : ""}`;
  console.log(line);
  if (!ok) throw new Error(line);
}

function tp_expectJson_(output) {
  // doPost returns a ContentService.TextOutput (JSON)
  if (!output || typeof output.getContent !== "function") {
    throw new Error("Output is not a TextOutput");
  }
  const text = output.getContent();
  try { return JSON.parse(text); }
  catch (e) { throw new Error("Response is not valid JSON: " + text); }
}

function tp_makeEvent_(obj) {
  const json = JSON.stringify(obj || {});
  return {
    postData: {
      contents: json,
      type: "application/json",
      length: json.length
    },
    // keep e.parameter/e.parameters empty for POST tests
    parameter: {},
    parameters: {}
  };
}

//////////////////// Stubs & Restore ////////////////////

var __orig__ = {};

/** Replace a global with a stub and remember original for restore. */
function tp_stub_(name, fn) {
  __orig__[name] = this[name];
  this[name] = fn;
}

/** Restore a previously stubbed global. */
function tp_restore_(name) {
  if (name in __orig__) this[name] = __orig__[name];
}

function tp_restoreAll_() {
  Object.keys(__orig__).forEach(tp_restore_);
  __orig__ = {};
}

//////////////////// Shared helpers used by doPost ////////////////////

// Minimal jsonOk/jsonError used by handlers/router
function tp_jsonOk_(data, meta) {
  return ContentService.createTextOutput(
    JSON.stringify({ ok: true, data: data || null, meta: meta || undefined })
  ).setMimeType(ContentService.MimeType.JSON);
}
function tp_jsonError_(message, code, meta) {
  const msg = (message && message.message) || String(message || "Unknown error");
  const body = { ok: false, error: { message: msg } };
  if (code !== undefined) body.error.code = code;
  if (meta) body.meta = meta;
  return ContentService.createTextOutput(JSON.stringify(body))
    .setMimeType(ContentService.MimeType.JSON);
}

//////////////////// Test Suite ////////////////////

function test_doPost_invalid_json_body() {
  // Arrange: event with broken JSON
  const e = {
    postData: { contents: "{not_json", type: "application/json", length: 9 },
    parameter: {}, parameters: {}
  };

  // Stubs (only what doPost touches)
  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  // Act
  const out = doPost(e);
  const json = tp_expectJson_(out);

  // Assert
  tp_logResult_("invalid JSON → 400", json.ok === false && json.error.code === 400);

  tp_restoreAll_();
}

function test_doPost_missing_token() {
  const e = tp_makeEvent_({ action: "HEALTHCHECK" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  // Token validator returns invalid
  tp_stub_("isValidApiToken", function(tok){ return { valid:false, preview: tok ? String(tok).slice(0,4) : null }; });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  tp_logResult_("missing token → 401", json.ok === false && json.error.code === 401);

  tp_restoreAll_();
}

function test_doPost_invalid_token() {
  const e = tp_makeEvent_({ token: "bad", action: "HEALTHCHECK" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  tp_stub_("isValidApiToken", function(){ return { valid:false, preview:"bad" }; });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  tp_logResult_("invalid token → 401", json.ok === false && json.error.code === 401);

  tp_restoreAll_();
}

function test_doPost_missing_action() {
  const e = tp_makeEvent_({ token: "ok" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  tp_stub_("isValidApiToken", function(){ return { valid:true, preview:"ok" }; });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  tp_logResult_("missing action → 400", json.ok === false && json.error.code === 400);

  tp_restoreAll_();
}

function test_doPost_unknown_action() {
  const e = tp_makeEvent_({ token: "ok", action: "NOT_EXISTS" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  tp_stub_("isValidApiToken", function(){ return { valid:true, preview:"ok" }; });

  // Public map returns nothing; fallback is the same "NOT_EXISTS"
  tp_stub_("getPublicPostActions_", function(){ return {}; });
  // Registry also has nothing
  tp_stub_("getPostActionsRegistry_", function(){ return {}; });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  tp_logResult_("unknown action → 404", json.ok === false && json.error.code === 404);

  tp_restoreAll_();
}

function test_doPost_success_healthcheck_alias() {
  const e = tp_makeEvent_({ token: "ok", action: "HEALTHCHECK", initiator: "Test" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  tp_stub_("isValidApiToken", function(){ return { valid:true, preview:"ok" }; });

  // Alias mapping HEALTHCHECK -> HEALTHCHECK_V1
  tp_stub_("getPublicPostActions_", function(){
    return { HEALTHCHECK: "HEALTHCHECK_V1" };
  });

  // Registry points V1 to handler
  tp_stub_("getPostActionsRegistry_", function(){
    return { "HEALTHCHECK_V1": function(body, ctx){
      return tp_jsonOk_({ status: "ok" }, { actionId: ctx.actionId, initiator: ctx.initiator, requestId: ctx.requestId });
    }};
  });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  const ok = json.ok === true && json.data && json.data.status === "ok" &&
             json.meta && json.meta.actionId === "HEALTHCHECK_V1";
  tp_logResult_("success (alias HEALTHCHECK → V1)", ok);

  tp_restoreAll_();
}

function test_doPost_success_healthcheck_v1() {
  const e = tp_makeEvent_({ token: "ok", action: "HEALTHCHECK_V1", initiator: "Test" });

  tp_stub_("jsonOk", tp_jsonOk_);
  tp_stub_("jsonError", tp_jsonError_);
  tp_stub_("logRequestSummary_", function(){});
  tp_stub_("logInfo", function(){});
  tp_stub_("logWarn", function(){});
  tp_stub_("logError", function(){});
  tp_stub_("spanStart_", function(){ return 0; });
  tp_stub_("spanEndMs_", function(){ return 1; });

  tp_stub_("isValidApiToken", function(){ return { valid:true, preview:"ok" }; });

  // Map returns identity for explicit V1
  tp_stub_("getPublicPostActions_", function(){
    return { HEALTHCHECK_V1: "HEALTHCHECK_V1" };
  });

  tp_stub_("getPostActionsRegistry_", function(){
    return { "HEALTHCHECK_V1": function(body, ctx){
      return tp_jsonOk_({ status: "ok" }, { actionId: ctx.actionId, initiator: ctx.initiator, requestId: ctx.requestId });
    }};
  });

  const out = doPost(e);
  const json = tp_expectJson_(out);

  const ok = json.ok === true && json.data && json.data.status === "ok" &&
             json.meta && json.meta.actionId === "HEALTHCHECK_V1";
  tp_logResult_("success (explicit V1)", ok);

  tp_restoreAll_();
}

/** Run the whole suite */
function test_doPost_all() {
  try {
    test_doPost_invalid_json_body();
    test_doPost_missing_token();
    test_doPost_invalid_token();
    test_doPost_missing_action();
    test_doPost_unknown_action();
    test_doPost_success_healthcheck_alias();
    test_doPost_success_healthcheck_v1();
    console.log("🎉 All doPost tests passed");
  } catch (e) {
    console.error("❌ doPost tests failed:", e && e.message || e);
    throw e;
  }
}
