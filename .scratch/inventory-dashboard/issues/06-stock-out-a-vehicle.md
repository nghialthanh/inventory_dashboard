# 06: Stock out a vehicle

**Decision:** [ADR-0001](../../../docs/adr/0001-retain-stocked-out-vehicles.md), [ADR-0004](../../../docs/adr/0004-stocked-out-is-terminal.md), [ADR-0005](../../../docs/adr/0005-stock-state-from-stocked-out-on.md), and [ADR-0007](../../../docs/adr/0007-separate-commands.md). This file is the acceptance checklist for those decisions.

**Blocked by:** 03 Record a vehicle and see it on the list

**Status:** ready-for-human

- [x] The manager records a stocked-out on date and the vehicle remains listed as stocked out
- [x] Stocked-out vehicles look different from vehicles that are in stock
- [x] Inventory age stops on stocked-out on. A stocked-out vehicle is not aging stock, including one whose stay was longer than 90 days
- [x] A stocked-out on date before entered on, or after today, is rejected. An empty date is rejected
- [x] Stocking out an unknown vehicle is not found
- [x] Recording the same VIN again, including the VIN of a stocked-out vehicle, is a conflict and the manager is told the VIN is already used
- [x] The vehicle cannot be returned to in stock. An attempt to reopen it says the departure stands
- [x] A demo vehicle that stocked out after a stay longer than 90 days is present so the retained record is visible on first open
- [x] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed, including the derived stock state, inventory age, and aging-stock flag
