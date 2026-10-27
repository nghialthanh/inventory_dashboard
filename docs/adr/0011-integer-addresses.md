# Catalog and vehicle addresses are database integers

**Superseded in part by [ADR-0013](0013-snowflake-vehicle-addresses.md).** Make and model addresses remain PostgreSQL identity integers. The vehicle address is no longer a database sequence: it is an app-minted snowflake (see ADR-0013).

Make and model addresses are PostgreSQL identity integers. The VIN stays the business identity and is immutable after record ([ADR-0012](0012-vin-immutable.md)).

This supersedes the UUID representation in [ADR-0006](0006-uuid-vehicle-identity.md). The URL still does not use the VIN.

## Considered options

- UUIDs minted in the app. That is what ADR-0006 chose. The catalog and the inventory list do not need an unguessable address, and the extra width is carried in every foreign key and URL.
- A snowflake with a timestamp and a node id, minted in the pod. Chosen later in ADR-0013.
- The VIN in the URL. Rejected in ADR-0006. The URL must stay an opaque address, not the business identity.

## Consequences

- Catalog: the API uses JSON numbers for make and model addresses. Path and query addresses are decimal text, such as `make_id=2`. A string, a fraction, zero, or a negative number is a validation error. A make or model address has to fit in a 32-bit integer.
- Vehicle address consequences moved to ADR-0013.
