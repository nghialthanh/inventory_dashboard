# Vehicle addresses are app-minted snowflakes; vehicles are monthly-partitioned

The API pod mints the vehicle address as a Twitter-style snowflake (41-bit time, 10-bit worker, 12-bit sequence, fixed epoch 2024-01-01 UTC). `created_at` is set to that same instant. The address is a JSON string and decimal path text so values above `2^53 − 1` stay exact in JavaScript. Make and model addresses remain PostgreSQL identity integers and JSON numbers.

`vehicles` is RANGE-partitioned by month on `created_at` (2026-01 through 2027-12 in migration; no DEFAULT partition). The primary key is `(id, created_at)`. Point reads prune with the month decoded from the snowflake. Global VIN uniqueness lives in non-partitioned `vehicle_vins` (vin → vehicle address), written in the same transaction as the vehicle row.

This supersedes the sequence-minted vehicle address and `2^53 − 1` cap in [ADR-0011](0011-integer-addresses.md). Catalog addresses in ADR-0011 still stand. Worker id comes from config (`WORKER_ID`, default 0); colliding workers can mint colliding addresses.

## Considered options

- Keep a PostgreSQL identity sequence (ADR-0011). Rejected for this change: the operator wants pod-minted ids and month partitions now.
- Partition on `entered_on`. Rejected: `entered_on` is correctable while in stock; mint/`created_at` is not.
- Skip `vehicle_vins` and rely on per-partition uniqueness. Rejected: VIN must stay unique across months ([ADR-0004](0004-stocked-out-is-terminal.md)).

## Consequences

- Drop the Compose volume (or otherwise rebuild `vehicles`) before starting; there is no in-place rewrite of sequence rows.
- Demo seed pins four fixed snowflake literals for worker 0.
- A missing month partition makes insert fail until a partition is added.
