# 04: See aging stock

**Decision:** Aging stock means in stock and inventory age greater than 90. That rule lives in the architecture invariants and the glossary. This file is the acceptance checklist.

**Blocked by:** 03 Record a vehicle and see it on the list

**Status:** ready-for-agent

- [x] Aging stock looks different from normal in-stock vehicles, using the aging-stock flag returned by the API
- [x] The screen does not apply a second aging rule of its own
- [x] The manager can filter the list to aging stock only
- [x] Leaving the aging filter off does not hide non-aging vehicles. `aging=true` is aging stock only; `aging=false` is not aging stock; any other value is a validation error
- [x] On first open, demo vehicles cover inventory age 0, exactly 90 days, and 91 days
- [x] The aging-only filter and the 0-, 90-, and 91-day outcomes are verified through the HTTP API against PostgreSQL with the clock fixed
