/**
 * Central config access for Script Properties.
 * Keeps magic strings in one place.
 */

const CONFIG_KEYS = {
    ZAPIER_TOKEN: "ZAPIER_TOKEN",
    GOOGLE_MAPS_API_KEY: "GOOGLE_MAPS_API_KEY",
    GOOGLE_MAPS_GEOCODE_URL: "GOOGLE_MAPS_GEOCODE_URL"
  };
  
  /**
   * Returns a script property or throws if missing.
   *
   * @param {string} key - The key name in Script Properties.
   * @param {Object} [options] - { required: boolean, defaultValue: any }
   * @returns {string} The property value.
   */
  function getScriptProp(key, options) {
    const val = PropertiesService.getScriptProperties().getProperty(key);
  
    if (!val || val.trim() === "") {
      if (options && options.defaultValue !== undefined) return options.defaultValue;
      if (options && options.required === false) return null;
      throw new Error(`Missing required Script Property: ${key}`);
    }
  
    return val.trim();
  }
  
  /**
   * Sets a script property.
   */
  function setScriptProp(key, value) {
    PropertiesService.getScriptProperties().setProperty(key, value);
  }
  