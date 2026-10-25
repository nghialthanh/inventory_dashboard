# 07: Review departures

**What to build:** The manager filters the list to vehicles in stock or stocked out, and to vehicles stocked out on one calendar day. They can correct a wrong stocked-out on date. They cannot clear that date or move it outside the vehicle's stay.

**Blocked by:** 06 Stock out a vehicle

**Status:** ready-for-agent

- [ ] The manager can filter the list to in stock or to stocked out
- [ ] The manager can filter the list to vehicles stocked out on one calendar day
- [ ] Those filters combine with make, model, and aging stock
- [ ] A later stock-out records a new stocked-out on date and replaces the previous one
- [ ] Clearing stocked-out on is rejected, and the vehicle stays stocked out
- [ ] A replacement date before entered on, or after today, is rejected
- [ ] Stocked-out on changes only through the stock-out command
- [ ] This behavior is verified through the HTTP API against PostgreSQL with the clock fixed
