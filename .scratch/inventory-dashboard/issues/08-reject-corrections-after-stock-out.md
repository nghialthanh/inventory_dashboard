# 08: Reject corrections after stock-out

**Decision:** [ADR-0004](../../../docs/adr/0004-stocked-out-is-terminal.md) and [ADR-0006](../../../docs/adr/0006-uuid-vehicle-identity.md). This file is the acceptance checklist for those decisions.

**Blocked by:** 05 Correct an in-stock vehicle, 06 Stock out a vehicle

**Status:** ready-for-human

- [x] Correcting the VIN, the model, or entered on of a stocked-out vehicle is a conflict
- [x] The manager is told the correction is rejected because the vehicle is stocked out
- [x] The VIN, model, and entered on stay as they were
- [x] This behavior is verified through the HTTP API against PostgreSQL
