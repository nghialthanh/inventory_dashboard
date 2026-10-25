# Vehicle URLs use an opaque id

The API identifies a vehicle by a UUID. VIN stays unique, and it can be corrected while the vehicle is in stock. A stocked-out VIN is never reused. The URL stays the same when the VIN is corrected.

## Considered options

- VIN in the URL. A correction changes the identity clients already stored, and the URL mixes a mutable attribute with the route.
