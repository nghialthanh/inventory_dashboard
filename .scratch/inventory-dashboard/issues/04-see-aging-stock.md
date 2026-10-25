# 04: See aging stock

**What to build:** The manager can tell aging stock from normal in-stock vehicles on the list, and can filter the list down to aging stock only. Demo vehicles make the 90-day boundary visible on first open.

**Blocked by:** 03 Record a vehicle and see it on the list

**Status:** ready-for-agent

- [x] Aging stock looks different from normal in-stock vehicles, using the aging-stock flag returned by the API
- [x] The screen does not apply a second aging rule of its own
- [x] The manager can filter the list to aging stock only
- [x] Leaving the aging filter off does not hide non-aging vehicles. Any aging filter other than the aging-only value is a validation error
- [x] On first open, demo vehicles cover inventory age 0, exactly 90 days, and 91 days
- [x] The aging-only filter and the 0-, 90-, and 91-day outcomes are verified through the HTTP API against PostgreSQL with the clock fixed
