# API Actions – GET Requests

This document explains the **`GetActions`** and **`PublicGetActions`** pattern used in our Google Apps Script backend for routing `doGet()` requests.

---

## Overview

We separate **internal** action identifiers from the **public** API surface.

- **`GetActions`** → Internal constants used inside our router and handlers.
- **`PublicGetActions`** → External constants exposed to API consumers (Zapier, Postman, integrations).

This mapping layer allows us to **change internal code without breaking the public API** and to **version endpoints safely**.

---

## File Location

| File                                 | Purpose                                          |
|--------------------------------------|--------------------------------------------------|
| `src/shared/actions.get.gs`          | Defines `GetActions` and `PublicGetActions` for all GET endpoints. |
| `src/shared/actions.post.gs` *(TBD)* | Defines actions for POST endpoints.              |

---

## Structure

### Internal Actions (`GetActions`)

- Represent the **actual** work our backend can perform.
- Used by the router `switch` statements (`handleGetV1`, `handleGetV2`, etc.).
- Can be renamed, reorganized, or replaced without affecting public clients.

```js
// Internal GET action identifiers
const GetActions = {
  HEALTHCHECK: "HEALTHCHECK",
  GET_LAT_LNG: "GET_LAT_LNG",
  GET_LAT_LNG_V2: "GET_LAT_LNG_V2", // New internal action for v2
  CALCULATE_ESTIMATE: "CALCULATE_ESTIMATE",
};
```

---

### Public Actions (`PublicGetActions`)

- The **public API contract**: what external callers pass via `?action=...`.
- Map 1:1 to an internal `GetActions` constant.
- Can include **version suffixes** (e.g., `_V1`, `_V2`) or convenience aliases.
- Old versions can coexist with new ones during migrations.

```js
// Public-facing GET action names → mapped to internal IDs
const PublicGetActions = {
  HEALTHCHECK_V1: GetActions.HEALTHCHECK,

  GET_LAT_LNG_V1: GetActions.GET_LAT_LNG,
  GET_LAT_LNG_V2: GetActions.GET_LAT_LNG_V2,

  // Optional alias pointing to the latest version
  GET_LAT_LNG: GetActions.GET_LAT_LNG_V2
};
```

---

## How Routing Uses This

When a request comes in:

```
GET /?action=GET_LAT_LNG_V1&v=1&token=...
```

1. `doGet()` reads `action = "GET_LAT_LNG_V1"`.
2. Look up `PublicGetActions[action]` to find the matching internal constant:
   ```js
   const internalAction = PublicGetActions[action];
   ```
3. Router calls `handleGetV1()` and uses a `switch` on **`GetActions`**:
   ```js
   switch (internalAction) {
     case GetActions.GET_LAT_LNG:
       return handleGetLatLngV1(e);
   }
   ```

---

## Versioning Example

| Timeline | `GetActions` (internal)               | `PublicGetActions` (public)                   | Effect |
|----------|---------------------------------------|------------------------------------------------|--------|
| Launch V1 | `GET_LAT_LNG`                         | `GET_LAT_LNG_V1`                               | Basic lat/lng lookup |
| Add V2   | `GET_LAT_LNG`, `GET_LAT_LNG_V2`        | `GET_LAT_LNG_V1`, `GET_LAT_LNG_V2`, `GET_LAT_LNG` | V1 + V2 both work; alias points to V2 |
| Drop V1  | `GET_LAT_LNG_V2`                       | `GET_LAT_LNG_V2`, `GET_LAT_LNG`                | Only V2 remains |

---

## Best Practices

- **Always** add new internal actions to `GetActions` first.
- Add matching entries in `PublicGetActions` for **each version** you want exposed.
- Use `_V1`, `_V2`, etc. in public names for versioning.
- Use plain uppercase without version suffix in public actions **only** as an alias to the latest version.
- Never remove a public action without checking that no clients depend on it.
- Keep **all GET action constants in one file** for discoverability and maintainability.

---

## Benefits of This Pattern

1. **Separation of Concerns**  
   Internal refactors don't break the public API.

2. **Safe Versioning**  
   You can run old and new versions side-by-side.

3. **Easier Migrations**  
   Aliases let you redirect traffic to new handlers without forcing all clients to change immediately.

4. **Clear Documentation**  
   `PublicGetActions` is the single source of truth for the public API surface.

---

## Related SOP

- All new GET endpoints:
  1. Add to `GetActions`.
  2. Add to `PublicGetActions` with correct version suffix.
  3. Implement handler in `src/<domain>/<domain>.handlers.get.vX.gs`.
  4. Wire up in `handleGetVX()` inside `src/router/handle.get.vX.gs`.

- All new POST endpoints: *(same pattern but use `PostActions` / `PublicPostActions`)*

---