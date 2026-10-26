# Vehicle URLs use an opaque id

Superseded in part by [ADR-0011](0011-integer-addresses.md). The vehicle address is a bigint assigned by PostgreSQL, not a UUID. The URL still uses that address rather than the VIN, and the address stays the same when the VIN is corrected.

The API identified a vehicle by a UUID. VIN stays unique, and it can be corrected while the vehicle is in stock. A stocked-out VIN is never reused. The URL stays the same when the VIN is corrected.

## Considered options

- VIN in the URL. A correction changes the identity clients already stored, and the URL mixes a mutable attribute with the route.

Accepted in [03](../../.scratch/inventory-dashboard/issues/03-record-a-vehicle.md), [05](../../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md), and [08](../../.scratch/inventory-dashboard/issues/08-reject-corrections-after-stock-out.md).
