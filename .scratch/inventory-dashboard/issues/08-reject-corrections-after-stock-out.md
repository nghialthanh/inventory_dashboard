# 08: Reject corrections after stock-out

**What to build:** Once a vehicle is stocked out, the manager can no longer correct its VIN, its model, or entered on. The dashboard tells them the correction is rejected because the vehicle has stocked out.

**Blocked by:** 05 Correct an in-stock vehicle, 06 Stock out a vehicle

**Status:** ready-for-agent

- [ ] Correcting the VIN, the model, or entered on of a stocked-out vehicle is a conflict
- [ ] The manager is told the correction is rejected because the vehicle is stocked out
- [ ] The VIN, model, and entered on stay as they were
- [ ] This behavior is verified through the HTTP API against PostgreSQL
