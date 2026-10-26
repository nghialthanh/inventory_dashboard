Status: ready-for-agent

# Build for the future

## Problem Statement

The inventory dashboard already lets a manager filter vehicles, see aging stock, and record a proposed action. The assessment also asks the design to consider scalability, performance, reliability, maintainability, and observability as the number of vehicles grows.

Listing gets more expensive as stocked-out vehicles stay on record forever. The make and model catalog is read often and changes only when the seed changes. The process can report itself up while PostgreSQL is unreachable, and there is no measure of request rate, error rate, or latency. Splitting `vehicles` into a table per month, indexing every column, and adding Redis were considered as remedies. Those three would fight VIN uniqueness, the UUID address of a vehicle, and the decision that PostgreSQL is the system of record.

## Solution

PostgreSQL stays the only system of record, in one database, with one `vehicles` table. Aging stock stays “in stock and inventory age greater than 90”. The list keeps filtering that predicate in SQL. New indexes match the predicates the list actually runs, including a partial index of in-stock vehicles ordered by entered on, so stocked-out history can grow without becoming the aging-stock read.

The make and model catalog is cached in the API process after the first successful read and dropped on restart. Vehicle lists still join make and model in the same query as the vehicle.

Liveness stays a process check. A separate readiness check pings PostgreSQL. Request rate, error rate, and latency are exported for each route, plus counts of vehicles recorded, stock-outs, and proposed-action writes. Trace id on the request stays as it is.

Proposed action is unchanged: the manager sets, replaces, or clears price reduction planned, wholesale planned, or hold on aging stock only.

## User Stories

1. As a dealership manager, I want aging stock to remain vehicles that are in stock with inventory age greater than 90, so that a faster list does not change who I work.
2. As a dealership manager, I want a vehicle in stock for exactly 90 days to stay off the aging list, so that the boundary stays the day after 90.
3. As a dealership manager, I want a vehicle in stock for 91 days to appear on the aging list, so that the first aging day is unchanged.
4. As a dealership manager, I want a stocked-out vehicle whose stay was longer than 90 days to stay off the aging list, so that departure still ends aging.
5. As a dealership manager, I want inventory age of a stocked-out vehicle to stop on stocked-out on, so that the stored stay does not keep growing.
6. As a dealership manager, I want the aging list to be driven by entered on before the cutoff day and an empty stocked-out on, so that the database applies the same rule the screen shows.
7. As a dealership manager, I want an old entered-on month that contains both in-stock and stocked-out vehicles to contribute only its in-stock rows past 90 days, so that “that month’s table” is not the definition of aging stock.
8. As a dealership manager, I want the unfiltered list ordered by entered on, then by vehicle id, so that the longest-held vehicles stay nearest the top.
9. As a dealership manager, I want paging and the reported total to stay the same, so that a large inventory is still one page at a time.
10. As a dealership manager, I want a page past the end to be an empty list with the true total, so that I can tell the end from a failed load.
11. As a dealership manager, I want an empty inventory to be an empty list, so that “no vehicles” stays distinct from an error.
12. As a dealership manager, I want to filter by make, so that one manufacturer is unchanged.
13. As a dealership manager, I want to filter by model, so that one vehicle line is unchanged.
14. As a dealership manager, I want a make and model pair that do not belong together to be rejected, so that the catalog rule still holds.
15. As a dealership manager, I want to filter by stock state, so that in stock and stocked out stay separate.
16. As a dealership manager, I want to filter by the calendar day stocked out on, so that one day’s departures stay reviewable.
17. As a dealership manager, I want to combine make, model, stock state, aging stock, and stocked-out on, so that a narrow slice stays correct.
18. As a dealership manager, I want each list row to include make and model names from the same read as the vehicle, so that the list does not look the catalog up once per row.
19. As a dealership manager, I want to open one vehicle by its id without knowing entered on, so that the address in the URL still finds that vehicle.
20. As a dealership manager, I want to record a vehicle with a VIN, a catalog model, and entered on, so that recording still starts the vehicle in stock with no proposed action.
21. As a dealership manager, I want a VIN that already exists, including on a stocked-out vehicle or on a vehicle from another entered-on month, to be rejected, so that a VIN is never reused.
22. As a dealership manager, I want two simultaneous records of the same VIN to result in one vehicle, so that the unique VIN still holds under concurrency.
23. As a dealership manager, I want to correct VIN, model, or entered on only while the vehicle is in stock, so that a stocked-out vehicle stays terminal.
24. As a dealership manager, I want stock-out and proposed action to keep working when I know only the vehicle id, so that those commands do not depend on which month the vehicle entered.
25. As a dealership manager, I want to set, replace, or clear a proposed action only on aging stock, so that the intention stays limited to vehicles in stock more than 90 days.
26. As a dealership manager, I want hold to stay a real proposed action, distinct from having none, so that “keep the current price” is still a decision.
27. As a dealership manager, I want the last proposed action to remain visible after stock-out and to reject a later change, so that the intention freezes with the departure.
28. As a dealership manager, I want any proposed action outside price reduction planned, wholesale planned, and hold to be rejected, so that the vocabulary stays closed.
29. As a dealership manager, I want the make list to be the seeded catalog, so that I still choose a manufacturer rather than typing one.
30. As a dealership manager, I want models for a make to be only that make’s models, so that the pair stays valid.
31. As a dealership manager, I want a second load of makes and models in the same running process to return the same catalog, so that repeating the catalog call does not keep rereading PostgreSQL.
32. As a dealership manager, I want the first catalog load after a restart to come from PostgreSQL, so that a new seed is what I see after a deploy.
33. As a dealership manager, I want a failed catalog read to be an error and not an empty cached catalog, so that a database failure is not stored as “no makes”.
34. As a dealership manager, I want recording a vehicle to reject an unknown model after the catalog has been read, so that the cache cannot accept a model PostgreSQL does not have.
35. As a dealership manager, I want recording a vehicle to accept a model that exists on the first try, before any catalog list has been opened, so that the cache miss still checks PostgreSQL.
36. As a dealership manager, I want the catalog to stay read-only in daily use, so that the screen still cannot create a make or a model.
37. As a dealership manager, I want inventory responses to keep the same envelope and fields, so that the screen and the cURL examples keep working.
38. As the person running the API, I want liveness to succeed when the process is up and PostgreSQL is down, so that a database outage does not restart the process in a loop.
39. As the person running the API, I want readiness to fail when PostgreSQL cannot be pinged, so that callers can wait until vehicles can be recorded.
40. As the person running the API, I want readiness to succeed once PostgreSQL accepts a connection after migrations, so that a healthy database is visible.
41. As the person running the API, I want readiness to fail quickly when the database hangs, so that a stuck connection is not treated as ready.
42. As the person running the API, I want Compose to treat the API as healthy only when readiness succeeds, so that the stack is not called healthy on process liveness alone.
43. As the person running the API, I want request count, error count, and latency for each inventory route, so that a slow list or a rise in failures is visible.
44. As the person running the API, I want client errors and server errors counted apart, so that a validation failure is not an outage.
45. As the person running the API, I want those series labeled by the route pattern, so that one vehicle id does not become its own series.
46. As the person running the API, I want VIN and vehicle id kept out of metric labels, so that the series stay bounded as inventory grows.
47. As the person running the API, I want liveness, readiness, and the metrics scrape left out of the request series, so that scrapes do not hide inventory traffic.
48. As the person running the API, I want a count of successful vehicle records, stock-outs, and proposed-action writes, so that I can see how the commands are used.
49. As the person running the API, I want a rejected record, stock-out, or proposed action to leave those counts unchanged, so that conflicts are not counted as completed commands.
50. As the person running the API, I want the trace id to stay on the response and on the request log, so that one failed call can still be found in the logs.
51. As the person running the API, I want metrics in Prometheus text, so that a scraper can read them without the JSON envelope.
52. As a dealership manager, I want a domain error to stay a validation error, not found, or conflict, so that metrics work does not change the contract of a bad VIN or a missing vehicle.

## Implementation Decisions

- One new ADR records the scale choices below. It does not change stock state, inventory age, aging stock, terminal stock-out, VIN uniqueness, UUID identity, the catalog relationship, separate commands, or one origin. Monthly partitions and Redis are named in that ADR as options that were rejected.
- `vehicles` stays one table. It is not partitioned or split by entered-on month. Aging stock is in stock and entered on strictly before the UTC day 90 days ago. A month table would still have to be read for every month older than that cutoff, which becomes almost the whole history. PostgreSQL can enforce a unique constraint on a partitioned table only when the constraint includes the partition key, so a global unique VIN and a primary key of vehicle id alone cannot both stand if the partition key is the entered-on month. Lookup by id is how correct, stock-out, proposed action, and the vehicle view find a row. That lookup must not require the month.
- Indexes follow list predicates. They are added in a new idempotent migration, the same way the current schema migrations reapply on startup.
  - In stock, including aging stock, sorted by entered on then id: `(entered_on, id)` where stocked-out on is null. Aging adds the entered-on upper bound on that index.
  - Every stock state, sorted the same way: `(entered_on, id)`.
  - Stocked out, and an exact stocked-out on day, then the same sort: `(stocked_out_on, entered_on, id)`.
  - A model filter with the same sort: `(model_id, entered_on, id)`.
- The existing single-column indexes on entered on, stocked-out on, and model id are dropped in that migration because the composites cover them. The unique VIN constraint stays.
- Proposed action, created at, and updated at are not indexed. They are not list predicates. An index on every column would tax record, correct, and stock-out without serving a query.
- Make filtering stays a join from vehicle to model to make. The model’s unique `(make_id, name)` already supports finding a make’s models. The vehicle does not gain a copied make id.
- Vehicle list SQL stays one statement that joins make and model. The catalog cache is not used to assemble list rows.
- Catalog logic keeps the full make list and model list in memory after the first successful repository read. List makes, list models, and the model-to-make check used when recording or correcting a vehicle all read through that cache. A failed read is returned and not stored. There is no TTL. Restart drops the cache, which is the invalidation, because daily use does not edit the catalog. The repository remains the loader. The catalog service interface that handlers call stays the same.
- Redis is not added. The catalog fits in the process, is not written by the API, and would become a second store to operate. A later move to more than one API process is the point to revisit a shared cache. Until then, PostgreSQL is what a cold process reads.
- Liveness stays `GET /health` and does not ping PostgreSQL. Readiness is `GET /ready`, outside `/api/v1`, and pings PostgreSQL with a one-second timeout. Success uses the normal success envelope. Failure is HTTP 503 with the normal error envelope and `INTERNAL_ERROR`. The health service gains a readiness method; the HTTP adapter only maps the result. Compose’s API healthcheck calls `/ready`. `restart` stays failure-on-exit, so an unready database does not by itself kill the process. Liveness remains available for a check that means “the process is up”.
- `GET /metrics`, outside `/api/v1`, returns Prometheus text and is not wrapped in the JSON envelope. That exception is written next to health in the API outline.
- HTTP middleware records, for each request that is not liveness, readiness, or the metrics scrape: route pattern (not the raw path), status class (`2xx`, `4xx`, `5xx`), and latency. Labels do not include VIN, vehicle id, make id, or model id.
- Vehicle logic increments a counter only after a successful record, stock-out, or proposed-action write. The handler does not increment them.
- Access logging and trace id stay. Domain 4xx stays quiet in the logs. Unexpected failures stay logged with the trace id.
- Docs updated with the behavior: the database schema (indexes), the API outline (`/ready`, `/metrics`), and the system-design observability and “build for the future” sections. No glossary term is added. Readiness and the metrics scrape are operational, not inventory language.

## Testing Decisions

A good test asserts what a client observes: status, envelope, which vehicles come back, catalog contents, readiness, and metric series. It does not assert a query plan, a cache hit, or which index PostgreSQL chose.

One seam: the existing HTTP tests against PostgreSQL, with the UTC clock fixed where aging is involved. Extend that suite. Do not add a second test seam in the repository or in catalog logic.

Prior art is that suite’s coverage of aging boundaries, list filters, VIN conflicts, catalog reads, and liveness.

Cover at least:

- Aging boundaries stay 90 off the list and 91 on it, including a stocked-out stay longer than 90 days.
- Make, model, stock state, stocked-out on, aging, paging, and sort still return the same vehicles.
- The same VIN cannot be recorded twice, including when entered on falls in different months.
- Get, correct, stock-out, and proposed action still address the vehicle by id alone.
- Proposed action stories stay green without new rules.
- Two catalog reads in one process return the seeded makes and models. An unknown model is still rejected on record. The catalog still has no write route.
- `/health` stays successful. `/ready` succeeds against the test database. `/ready` returns 503 when the API is built against a database that cannot be pinged, using that same HTTP harness.
- After one successful list, record, stock-out, and proposed-action write, `/metrics` shows the route series and the three command counts. A rejected proposed action does not increase the proposed-action count. The series do not contain a vehicle id or a VIN.

## Out of Scope

- A vehicle table, partition, or shard per month.
- Redis, or any cache outside the API process.
- An index on every column.
- A new proposed-action command, screen, or rule. That command already exists.
- Read replicas, more than one API process, and a shared cache for that topology.
- OpenTelemetry, or replacing the current trace id.
- Changing the aging threshold, VIN shape, stock state, or the decision that stocked-out rows stay.
- Authentication, catalog editing, and a new inventory screen.
- Query-plan tests and load tests. Correctness of the list is in scope. Proving a particular plan is not.

## Further Notes

Proposed action is already implemented and should be demoed as it stands. On an aging vehicle the manager sends `PUT /api/v1/vehicles/:id/proposed-action` with `price_reduction_planned`, `wholesale_planned`, `hold`, or an empty value to clear. The inventory screen has a Proposed action control on the vehicle, enabled only when the vehicle is aging stock. Hold is stored. A vehicle that is in stock for 90 days or fewer, and a stocked-out vehicle, reject the write. The HTTP suite already covers those cases.

The interview line for this work: stocked-out history is retained, so the in-stock and aging reads are partial indexes rather than a monthly split that would drop global VIN uniqueness. The catalog is immutable at runtime, so the process caches it and PostgreSQL remains the system of record. Liveness and readiness are different signals. Rate, errors, and latency use the route pattern so the series stay finite.

The test seam is the existing HTTP suite against PostgreSQL. That is the only seam this spec adds behavior to.
