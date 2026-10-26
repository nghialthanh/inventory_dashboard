# 09: Record a proposed action on aging stock

**Decision:** [ADR-0007](../../../docs/adr/0007-separate-commands.md). This file is the acceptance checklist for that decision. Aging stock is the only vehicle that accepts the command.

**Blocked by:** 04 See aging stock, 06 Stock out a vehicle

**Status:** ready-for-human

- [x] On an aging vehicle the manager can set price reduction planned, wholesale planned, or hold, and sees that value on the vehicle
- [x] Hold is stored and shown as a proposed action, distinct from having none
- [x] A new proposed action replaces the previous one. An empty value clears it
- [x] Setting, replacing, or clearing a proposed action is rejected when the vehicle is not aging stock, including a vehicle in stock for exactly 90 days
- [x] Any value outside those three intentions is rejected
- [x] After stock-out the last proposed action remains visible, and a further change is rejected
- [x] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed
