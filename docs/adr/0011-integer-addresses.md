# Catalog and vehicle addresses are database integers

Make and model addresses are PostgreSQL identity integers. The vehicle address is a bigint from one sequence in this database, and that integer is the URL. The VIN stays the business identity. Correcting the VIN does not change the vehicle address. Several API pods share that sequence. A snowflake node id waits until a pod must mint an address without asking PostgreSQL.

This supersedes the UUID representation in [ADR-0006](0006-uuid-vehicle-identity.md). The URL still does not use the VIN.

## Considered options

- UUIDs minted in the app. That is what ADR-0006 chose. The catalog and the inventory list do not need an unguessable address, and the extra width is carried in every foreign key and URL.
- A snowflake with a timestamp and a node id, minted in the pod. That pays off when a pod cannot ask PostgreSQL. On Kubernetes the node id has to stay unique across restarts. Two pods that start with the same node id, or a clock step backward, can mint the same vehicle address.
- The VIN in the URL. Rejected in ADR-0006. A correction would change links already stored.

## Consequences

- The operator drops the database or the Compose volume. The existing `CREATE TABLE` and seed files are rewritten in place. There is no migration that rewrites UUID rows.
- The API uses JSON numbers for make, model, and vehicle addresses in bodies and responses. Path and query addresses are decimal text, such as `/vehicles/1` and `make_id=2`. A string, a fraction, zero, or a negative number is a validation error. A make or model address has to fit in a 32-bit integer. The vehicle column is a bigint, and a vehicle address above 9007199254740991 (`2^53 − 1`) is a validation error so a JavaScript number stays exact. Request trace ids stay UUIDs.
- The seed pins Ford, Honda, and Toyota as 1, 2, and 3. Models follow today's insert order from 1 (F-150, Mustang, Accord, Civic, Camry, Corolla). The four demo vehicles are 1 through 4 in today's insert order. Sequences continue after those pinned values.
- An address that is not a positive integer is a validation error. A positive integer that matches no vehicle is not found. An unknown model is rejected the same way it is today.
- A skipped sequence value is never reused. A failed insert can leave a gap.
