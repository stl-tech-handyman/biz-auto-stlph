function getExpectedZapierToken() {
  return PropertiesService.getScriptProperties().getProperty(ZAPIER_TOKEN_NAME);
}

function extractTokenFromRequest(e) {
  return e?.parameter?.token || null;
}

function isValidApiToken(token) {
  const props = PropertiesService.getScriptProperties();
  const knownTokens = {
    ZAPIER_TOKEN: props.getProperty("ZAPIER_TOKEN"),
    POSTMAN_TOKEN: props.getProperty("POSTMAN_TOKEN"),
    FRONTEND_TOKEN: props.getProperty("FRONTEND_TOKEN")
  };

  for (let [key, value] of Object.entries(knownTokens)) {
    if (token === value) {
      return { valid: true, source: key, preview: value?.slice(0, 4) };
    }
  }

  return { valid: false, source: null, preview: token?.slice(0, 4) || "N/A" };
}
