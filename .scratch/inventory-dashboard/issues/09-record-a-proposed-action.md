# 09: Record a proposed action on aging stock

**What to build:** On an aging vehicle, the manager sets a proposed action of price reduction planned, wholesale planned, or hold, replaces it, or clears it. Hold is different from having no proposed action. The same change is rejected when the vehicle is not aging stock. After the vehicle stocks out, the last proposed action stays visible and cannot be changed.

**Blocked by:** 04 See aging stock, 06 Stock out a vehicle

**Status:** ready-for-agent

- [ ] On an aging vehicle the manager can set price reduction planned, wholesale planned, or hold, and sees that value on the vehicle
- [ ] Hold is stored and shown as a proposed action, distinct from having none
- [ ] A new proposed action replaces the previous one. An empty value clears it
- [ ] Setting, replacing, or clearing a proposed action is rejected when the vehicle is not aging stock, including a vehicle in stock for exactly 90 days
- [ ] Any value outside those three intentions is rejected
- [ ] After stock-out the last proposed action remains visible, and a further change is rejected
- [ ] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed
