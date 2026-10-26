# 03: Record a vehicle and see it on the list

**What to build:** The manager records an in-stock vehicle with a VIN, a catalog model, and entered on, then sees it on the inventory list and on its own view. The row shows make, model, stock state, inventory age, and whether it is aging stock. The manager can filter by make and model, page through the list, and understand an empty list, a failed load, and a rejected VIN.

**Blocked by:** 02 Choose make and model from the catalog

**Status:** ready-for-agent

- [x] The manager records a vehicle with a VIN, a catalog model, and entered on, and it starts in stock with no proposed action
- [x] The list and the vehicle view show VIN, make, model, entered on, an empty stocked-out on, stock state, inventory age, whether it is aging stock, and proposed action
- [x] Make and model names appear on the row
- [x] The list is ordered by entered on, earliest first, then by vehicle id
- [x] The manager can page the list. The default page is 1 and the default page size is 20. Page size is capped at 100. A page past the end is empty and still reports the true total
- [x] Filtering by make, by model, or by both narrows the list. A model that does not belong to the selected make is a validation error
- [x] An inventory with no vehicles is an empty list. A failed load is shown as an error
- [x] A vehicle entered today has inventory age 0 and is not aging stock. Exactly 90 days in stock is not aging stock. 91 days in stock is aging stock
- [x] Today is the UTC calendar date of a clock the test can set. Entered on after today is rejected
- [x] A missing VIN, model, or entered on is rejected. An unknown model is rejected
- [x] The VIN is stored trimmed and uppercase, 17 characters, using digits and letters except I, O, and Q
- [x] A duplicate VIN is a conflict and the manager is told the VIN is already used. Two simultaneous records of the same VIN create only one vehicle
- [x] An unknown vehicle address is not found
- [x] Responses use the shared envelope and snake_case fields. Errors carry an error code for validation, not found, and conflict
- [x] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed, asserting the envelope and the derived stock state, inventory age, and aging-stock flag
