function cacheGet_(key) {
    var c = CacheService.getScriptCache();
    var raw = c.get(key);
    if (!raw) return null;
    try { return JSON.parse(raw); } catch (_) { return null; }
  }
  
  function cachePut_(key, value, ttlSec) {
    try {
      var c = CacheService.getScriptCache();
      c.put(key, JSON.stringify(value), Math.min(ttlSec || 3600, 21600));
    } catch (_) {}
  }
  