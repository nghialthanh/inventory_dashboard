Status: ready-for-agent

# Integer addresses

## Problem Statement

A dealership manager opens a make, a model, and a vehicle by a long address that the API invents for the vehicle. Those addresses are hard to read, and each API process creates the vehicle address on its own. Several API processes will share one database. The manager wants the database to hand out the next address, and the old addresses can be thrown away.

## Solution

Make and model addresses become ordinary increasing integers from the database. The vehicle address becomes a larger integer from one sequence in that same database, and that integer is what the URL uses. The VIN stays the business identity. Correcting the VIN does not change the vehicle address. The screen and the API show and accept these integers as numbers. The person running the API drops the old database first, so there is no map from an old address to a new one. A snowflake, with a timestamp and a node id, is not part of this change.

## User Stories

1. As a dealership manager, I want each make to have a small integer address, so that I can tell the manufacturers apart without a long token.
2. As a dealership manager, I want Ford, Honda, and Toyota to be addresses 1, 2, and 3 after a fresh database, so that the seeded catalog is stable.
3. As a dealership manager, I want each model to have a small integer address, so that I can choose a vehicle line by that address.
4. As a dealership manager, I want the seeded models to be addresses 1 through 6 in the order F-150, Mustang, Accord, Civic, Camry, and Corolla, so that a known model keeps a known address.
5. As a dealership manager, I want each model's address to name the make it belongs to, so that Civic still belongs to Honda.
6. As a dealership manager, I want the make list to return those integer addresses with the make names, so that the screen can build the make filter from the catalog.
7. As a dealership manager, I want the model list to return each model's integer address and its make's integer address, so that choosing a make can limit the models.
8. As a dealership manager, I want to ask for the models of make 2 and receive only Honda's models, so that the make filter still limits the list.
9. As a dealership manager, I want to read models without choosing a make, so that the full catalog is still available.
10. As a dealership manager, I want a make address that is not a positive integer to be rejected, so that a mistyped filter is not an empty catalog.
11. As a dealership manager, I want a model address that is not a positive integer to be rejected, so that I do not record a vehicle under a garbled model.
12. As a dealership manager, I want a make or model address above 2147483647 to be rejected, so that the address fits the catalog integer.
13. As a dealership manager, I want a quoted address such as `"2"`, a fraction, zero, or a negative number to be rejected in a JSON body, so that only a JSON number is an address.
14. As a dealership manager, I want `make_id=2` and `model_id=4` in the query string to mean those integers, so that a filter URL stays readable.
15. As a dealership manager, I want a model that is not one of the selected make's models to be rejected as a bad filter, so that Honda plus an F-150 is not an empty inventory.
16. As a dealership manager, I want an unknown model address to be rejected when I record or correct a vehicle, so that the vehicle still has to point at the catalog.
17. As a dealership manager, I want to record a vehicle with a VIN, a model address, and entered on, and receive a new vehicle address from the database, so that I do not invent that address myself.
18. As a dealership manager, I want the record request to omit the vehicle address, so that two managers cannot choose the same address.
19. As a dealership manager, I want that new vehicle to start in stock with no proposed action, so that recording still only starts the vehicle.
20. As a dealership manager, I want the vehicle address in the list and on the vehicle view to be a JSON number, so that the screen and the URL use the same integer.
21. As a dealership manager, I want the make and model nested on that vehicle to use their integer addresses, so that the row matches the catalog.
22. As a dealership manager, I want to open a vehicle at `/vehicles/1`, so that the URL is the vehicle address and not the VIN.
23. As a dealership manager, I want the four demo vehicles to be addresses 1 through 4 in today's seed order, so that a fresh database has a stable demonstration inventory.
24. As a dealership manager, I want the next vehicle I record on that fresh database to receive address 5 or higher, so that a new vehicle does not take a demo vehicle's address.
25. As a dealership manager, I want two vehicles recorded at the same time to receive two different addresses, so that the shared sequence still separates them.
26. As a dealership manager, I want a duplicate VIN to be rejected, including a VIN that belongs to a stocked-out vehicle, so that a VIN is never reused.
27. As a dealership manager, I want two attempts to record the same VIN at the same time to create only one vehicle, so that the VIN stays unique under concurrency.
28. As a dealership manager, I want to correct the VIN while the vehicle is in stock and keep the same vehicle address, so that links I already have still open that vehicle.
29. As a dealership manager, I want to correct the model or entered on while the vehicle is in stock without changing the vehicle address, so that those corrections stay on the same vehicle.
30. As a dealership manager, I want a correction that includes a proposed action or stocked-out on to be rejected, so that those stay on their own commands.
31. As a dealership manager, I want a correction of VIN, model, or entered on to be rejected once the vehicle is stocked out, so that a departed vehicle's identity and tenure stay fixed.
32. As a dealership manager, I want to stock out a vehicle by its integer address, so that departure still uses the vehicle address and not the VIN.
33. As a dealership manager, I want to replace or clear a proposed action by that same integer address, so that the command still addresses one vehicle.
34. As a dealership manager, I want a vehicle address that is not a positive integer, including a string, a fraction, zero, or a negative number, to be a validation error, so that a bad link is not treated as an empty garage.
35. As a dealership manager, I want a vehicle address above 9007199254740991 to be a validation error, so that the screen's number stays exactly the address the API stored.
36. As a dealership manager, I want a positive integer that matches no vehicle to be reported as not found, so that a stale link is distinct from a malformed one.
37. As a dealership manager, I want the list to filter by integer make and model addresses, so that the filters I already use keep working.
38. As a dealership manager, I want the list to filter by stock state, aging stock, and stocked-out on together with those integer addresses, so that a narrow slice stays correct.
39. As a dealership manager, I want the list ordered by entered on, then by vehicle address, so that the longest-held vehicles stay nearest the top when they entered on the same day.
40. As a dealership manager, I want paging and the reported total to stay the same, so that a large inventory is still one page at a time.
41. As a dealership manager, I want a vehicle entered today to show inventory age 0 and not be aging stock, so that the address change does not change the aging rule.
42. As a dealership manager, I want a vehicle in stock for exactly 90 days to stay off the aging list, and one in stock for 91 days to appear on it, so that the cutoff stays the day after 90.
43. As a dealership manager, I want a stocked-out vehicle whose stay was longer than 90 days to stay off the aging list, and its inventory age to stop on stocked-out on, so that departure still ends aging.
44. As a dealership manager, I want the catalog to stay read-only in daily use, so that I still cannot create a make or a model from the screen.
45. As a dealership manager, I want the screen's make filter, model filter, record, correct, stock-out, and proposed action to send the integer addresses the API returned, so that I can do that work without pasting an address by hand.
46. As the person running the API, I want to drop the database before this change runs, so that the old addresses are gone and the new tables are created empty.
47. As the person running the API, I want a second API process to draw the next vehicle address from the same database sequence, so that several processes can record vehicles without minting addresses themselves.
48. As the person running the API, I want the trace id on a request to stay a UUID, so that request logs are unchanged by the address change.
49. As the person running the API, I want a failed record to leave a gap in the vehicle addresses rather than reuse a skipped value, so that an address is never given to two vehicles.

## Implementation Decisions

ADR-0011 is the decision. The UUID representation in ADR-0006 is superseded. The URL still uses the vehicle address, not the VIN, and that address stays the same when the VIN is corrected. Stock state, inventory age, aging stock, terminal stock-out, VIN uniqueness, the catalog relationship, separate commands, and one origin are unchanged.

- The operator drops the database or the Compose volume before the new schema runs. The existing table definitions and the seed are rewritten in place. There is no migration that rewrites old addresses, and there is no map from an old address to a new one.
- `makes.id` and `models.id` are 32-bit identity integers. `models.make_id` and `vehicles.model_id` use that same integer. `vehicles.id` is a bigint identity. The database assigns an address when an insert omits it. The seed inserts the pinned addresses explicitly, then moves each sequence past the highest pinned value so the next insert does not reuse one.
- Pinned catalog addresses: Ford 1, Honda 2, Toyota 3. Models, in order: F-150 1, Mustang 2, Accord 3, Civic 4, Camry 5, Corolla 6. The four demo vehicles keep today's seed order and receive addresses 1, 2, 3, and 4. The next recorded vehicle receives 5 or higher.
- The record body stays VIN, model address, and entered on. The client does not send the vehicle address.
- JSON bodies and JSON responses use a JSON number for every make, model, and vehicle address. A JSON string, a fraction, zero, or a negative number is a validation error. Path and query addresses are decimal text (`/vehicles/1`, `make_id=2`, `model_id=4`). A make or model address above 2147483647 is a validation error. A vehicle address above 9007199254740991 is a validation error. A positive in-range integer that matches no vehicle is not found. An unknown model, and a model that is not of the selected make, stay validation errors.
- The screen's make, model, and vehicle addresses become numbers. Filters, record, correct, stock-out, proposed action, and the vehicle view use the integer the API returned. The screen does not mint an address.
- Request trace ids stay UUIDs. They are not make, model, or vehicle addresses.
- The API outline, the database schema doc, and the cURL examples are updated in the same change so they describe integer addresses. The glossary term is already the vehicle address. No new glossary term is added.

## Testing Decisions

- A good test exercises the HTTP API as a black box. It sends the request a manager's screen would send and asserts the status, the envelope, the JSON type of each address, the pinned seed addresses, the derived fields (`stock_state`, `inventory_age`, `is_aging`), and `error_code`. It does not assert sequence names, identity internals, or React output.
- The module under test is the API. Catalog reads and vehicle commands are observed through `/api/v1`. PostgreSQL is the real store. The test sets the clock to a fixed UTC date where aging is involved.
- Prior art is the existing HTTP suite against PostgreSQL for catalog reads, aging boundaries, list filters, VIN conflicts, corrections, stock-out, and proposed action. Extend that suite. Do not add a second automated seam in the vehicle logic packages, the catalog logic packages, or the React app.
- Cover at least: seeded makes 1, 2, and 3 and models 1 through 6; a recorded vehicle whose address is a JSON number and whose VIN correction returns that same address; the next address after the demo seed is not 1, 2, 3, or 4; a JSON string, a fraction, zero, a negative number, a make or model address above 2147483647, and a vehicle address above 9007199254740991 are validation errors; an unknown positive vehicle address is not found; an unknown model is rejected; make and model filters accept the pinned integers; aging boundaries, VIN uniqueness, terminal stock-out, and proposed action stay as they are.
- The screen is checked by using it against this API. The trace id check stays a UUID check and is not reused as the address suite.

## Out of Scope

- A snowflake, a timestamp in the address, a node id, or any generator that mints a vehicle address inside an API process.
- Keeping the old UUID addresses, translating them, or altering a database that still has UUID columns. The operator drops that database.
- Changing the VIN rules, the aging threshold, stock state, terminal stock-out, proposed action, or the decision that stocked-out rows stay.
- A new inventory screen, catalog editing, authentication, or a second test suite in the React app.
- Changing request trace ids away from UUIDs.

## Further Notes

The interview line for this work: several API processes share one PostgreSQL, so the database sequence is the allocator. The catalog is small, so its addresses are 32-bit identity integers. The vehicle address is a bigint capped at the largest integer JavaScript stores exactly, because the screen reads JSON numbers. The VIN remains the business identity.

The operator drops the database or the Compose volume before running the new schema. The schema, seed, API, and screen now use integer addresses.
