# 05: Correct an in-stock vehicle

**What to build:** While a vehicle is in stock, the manager corrects its VIN, its model, or entered on. The vehicle keeps the same address when the VIN changes. The correction cannot set a proposed action or stocked-out on.

**Blocked by:** 03 Record a vehicle and see it on the list

**Status:** ready-for-agent

- [x] While the vehicle is in stock, the manager can correct the VIN, the model, entered on, or any non-empty subset of those
- [x] After a VIN correction the vehicle opens at the same address as before
- [x] A correction that includes a proposed action or stocked-out on is rejected, and those values stay unchanged
- [x] An empty correction is rejected
- [x] A corrected entered on still cannot be after today, and a corrected VIN still follows the VIN rules and stays unique
- [x] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed
