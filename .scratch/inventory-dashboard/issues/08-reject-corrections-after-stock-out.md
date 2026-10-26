# 08: Reject corrections after stock-out

**What to build:** Once a vehicle is stocked out, the manager can no longer correct its VIN, its model, or entered on. The dashboard tells them the correction is rejected because the vehicle has stocked out.

**Blocked by:** 05 Correct an in-stock vehicle, 06 Stock out a vehicle

**Status:** ready-for-human

- [x] Correcting the VIN, the model, or entered on of a stocked-out vehicle is a conflict
- [x] The manager is told the correction is rejected because the vehicle is stocked out
- [x] The VIN, model, and entered on stay as they were
- [x] This behavior is verified through the HTTP API against PostgreSQL
