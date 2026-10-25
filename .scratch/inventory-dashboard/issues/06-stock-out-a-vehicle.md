# 06: Stock out a vehicle

**What to build:** The manager records the calendar day a vehicle stocked out. The vehicle stays on the list as stocked out, looks different from vehicles still in stock, and stops counting as aging stock. Its inventory age stops on stocked-out on. The same VIN cannot be recorded again, and the vehicle cannot return to in stock.

**Blocked by:** 03 Record a vehicle and see it on the list

**Status:** ready-for-agent

- [ ] The manager records a stocked-out on date and the vehicle remains listed as stocked out
- [ ] Stocked-out vehicles look different from vehicles that are in stock
- [ ] Inventory age stops on stocked-out on. A stocked-out vehicle is not aging stock, including one whose stay was longer than 90 days
- [ ] A stocked-out on date before entered on, or after today, is rejected. An empty date is rejected
- [ ] Stocking out an unknown vehicle is not found
- [ ] Recording the same VIN again, including the VIN of a stocked-out vehicle, is a conflict and the manager is told the VIN is already used
- [ ] The vehicle cannot be returned to in stock. An attempt to reopen it says the departure stands
- [ ] A demo vehicle that stocked out after a stay longer than 90 days is present so the retained record is visible on first open
- [ ] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed, including the derived stock state, inventory age, and aging-stock flag
