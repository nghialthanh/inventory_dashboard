# VIN is immutable after record

The operator is a vehicle distributor, not a manufacturer. The VIN is the manufacturer's identity for the physical car; the distributor records it as-is and does not change it. After a vehicle is recorded, `PATCH /api/v1/vehicles/:id` accepts only `model_id` and `entered_on` while the vehicle is in stock. A body that includes `vin` is a validation error.

This supersedes the VIN-correction language in [ADR-0006](0006-uuid-vehicle-identity.md) and [ADR-0011](0011-integer-addresses.md). The vehicle address remains the URL handle and is still not the VIN. VIN uniqueness, including against stocked-out rows, is unchanged ([ADR-0004](0004-stocked-out-is-terminal.md)).
