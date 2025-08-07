function getExpectedZapierToken() {
  return PropertiesService.getScriptProperties().getProperty(ZAPIER_TOKEN_NAME);
}

function extractTokenFromRequest(e) {
  return e?.parameter?.token || null;
}

function isValidZapierToken(token) {
  const expected = getExpectedZapierToken();
  return Boolean(token && token === expected);
}