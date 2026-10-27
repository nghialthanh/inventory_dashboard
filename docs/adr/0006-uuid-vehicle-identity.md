# Vehicle URLs use an opaque id

Superseded in part by [ADR-0011](0011-integer-addresses.md) and [ADR-0013](0013-snowflake-vehicle-addresses.md). The vehicle address is a snowflake bigint minted by the API, not a UUID. The URL uses that address rather than the VIN. VIN is immutable after record ([ADR-0012](0012-vin-immutable.md)).

The API identified a vehicle by a UUID. VIN stays unique. A stocked-out VIN is never reused. The URL is the vehicle address, not the VIN.

## Considered options

- VIN in the URL. The URL would mix the business identity with the route, and the distributor does not own or change that identity.

Accepted in [03](../../.scratch/inventory-dashboard/issues/03-record-a-vehicle.md), [05](../../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md), and [08](../../.scratch/inventory-dashboard/issues/08-reject-corrections-after-stock-out.md).
