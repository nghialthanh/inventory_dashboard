Status: ready-for-agent

# Inventory dashboard

## Problem Statement

A dealership manager can see which vehicles the dealership holds, which of those have been in stock too long, and what the manager intends to do about them. The manager needs one inventory dashboard that filters vehicles by make, model, and inventory age, marks aging stock (in stock more than 90 days), keeps a visible difference between normal in-stock vehicles and aging stock, and records a proposed action for each aging vehicle. Vehicles that have stocked out stay on record, including the calendar day they left.

## Solution

A single local application: a React inventory screen and a Go REST API on one origin, with PostgreSQL as the source of truth and a private MinIO bucket for an optional primary image. The manager works from an inventory list. Filters, aging stock, stock state, and inventory age all come from the API. Stock-out, proposed action, and the primary image are separate commands. Stock state is derived from whether stocked-out on is set. A stocked-out vehicle stays stocked out, and its VIN is never recorded again.

## User Stories

1. As a dealership manager, I want to see every vehicle the dealership has recorded, including vehicles that have stocked out, so that departure does not erase the inventory.
2. As a dealership manager, I want each vehicle to show its VIN, make, model, entered on, stocked-out on, stock state, inventory age, whether it is aging stock, its proposed action, and whether it has a primary image, so that I can judge it from the list.
3. As a dealership manager, I want make and model names on each row, so that I do not look them up separately.
4. As a dealership manager, I want to filter the list by make, so that I can focus on one manufacturer.
5. As a dealership manager, I want to filter the list by model, so that I can focus on one vehicle line.
6. As a dealership manager, I want model choices to belong to the make I selected, so that I cannot pair a model with the wrong manufacturer.
7. As a dealership manager, I want to filter the list to aging stock only, so that I can work the vehicles that have been in stock more than 90 days.
8. As a dealership manager, I want to filter the list by stock state, in stock or stocked out, so that I can separate what we hold from what has left.
9. As a dealership manager, I want to filter the list to vehicles stocked out on one calendar day, so that I can review that day's departures.
10. As a dealership manager, I want to combine make, model, stock state, aging stock, and stocked-out on, so that I can narrow to one slice of inventory.
11. As a dealership manager, I want the list ordered with the earliest entered on first, so that the longest-held vehicles are nearest the top.
12. As a dealership manager, I want to page through the list, so that a large inventory stays usable.
13. As a dealership manager, I want a page past the last vehicle to be empty and still report the true total, so that I know I have passed the end.
14. As a dealership manager, I want an empty inventory to be an empty list, so that I can tell "no vehicles" from a failed load.
15. As a dealership manager, I want a failed load to be reported as an error, so that I know to retry.
16. As a dealership manager, I want aging stock to look different from normal in-stock vehicles, so that the 90-day problem is obvious on the list.
17. As a dealership manager, I want stocked-out vehicles to look different from vehicles still in stock, so that retained departures are not mistaken for cars we hold.
18. As a dealership manager, I want a vehicle entered today to show inventory age 0 and not be aging stock, so that a new arrival is not flagged.
19. As a dealership manager, I want a vehicle in stock for exactly 90 days to not be aging stock, so that the cutoff is more than 90 days.
20. As a dealership manager, I want a vehicle in stock for 91 days to be aging stock, so that the day after the cutoff is flagged.
21. As a dealership manager, I want a stocked-out vehicle whose stay was longer than 90 days to not be aging stock, so that only vehicles we still hold are flagged.
22. As a dealership manager, I want inventory age for a stocked-out vehicle to stop on stocked-out on, so that days after departure do not count.
23. As a dealership manager, I want to choose make and model from a seeded catalog, so that a vehicle cannot be stored under a mistyped or impossible pair.
24. As a dealership manager, I want to record a new in-stock vehicle with a VIN, a catalog model, and entered on, so that the car enters inventory.
25. As a dealership manager, I want a new vehicle to start with no proposed action and no primary image, so that I record those only when I decide to.
26. As a dealership manager, I want a duplicate VIN to be rejected, including a VIN that belongs to a stocked-out vehicle, so that a VIN is never reused.
27. As a dealership manager, I want two attempts to record the same VIN at the same time to create only one vehicle, so that the VIN stays unique under concurrency.
28. As a dealership manager, I want a future entered on to be rejected, so that inventory age cannot start in the future.
29. As a dealership manager, I want a vehicle with a missing VIN, model, or entered on to be rejected, so that every vehicle is identifiable and ageable.
30. As a dealership manager, I want an unknown model to be rejected when I record a vehicle, so that the vehicle always points at the catalog.
31. As a dealership manager, I want the vehicle's address to stay the same when I correct its VIN, so that links I already have still open that vehicle.
32. As a dealership manager, I want to correct the VIN while the vehicle is in stock, so that a typo does not become a second vehicle.
33. As a dealership manager, I want to correct the model while the vehicle is in stock, so that a catalog mistake can be fixed.
34. As a dealership manager, I want to correct entered on while the vehicle is in stock, so that inventory age uses the real arrival day.
35. As a dealership manager, I want a correction to leave proposed action, primary image, and stocked-out on unchanged, so that those stay on their own commands.
36. As a dealership manager, I want a correction of VIN, model, or entered on to be rejected once the vehicle is stocked out, so that a departed vehicle's identity and tenure stay fixed.
37. As a dealership manager, I want to record the calendar day a vehicle stocked out, so that it leaves in-stock inventory and stays on record.
38. As a dealership manager, I want that vehicle to remain listed as stocked out, so that stock-out is not a deletion.
39. As a dealership manager, I want a stocked-out vehicle to stop being aging stock, so that the aging list is only what we still hold.
40. As a dealership manager, I want a stocked-out vehicle to stay stocked out, so that a fallen-through deal is not modeled as a return to in stock.
41. As a dealership manager, I want to correct the stocked-out on date after departure, so that a wrong departure day can be fixed.
42. As a dealership manager, I want clearing stocked-out on to be rejected, so that erasing the date cannot reopen the vehicle.
43. As a dealership manager, I want a stocked-out on date before entered on, or after today, to be rejected, so that the stay is a real span of calendar days.
44. As a dealership manager, I want stocking out an unknown vehicle to be reported as not found, so that I cannot record a departure that does not exist.
45. As a dealership manager, I want to set a proposed action on an aging vehicle to price reduction planned, wholesale planned, or hold, so that my intention is recorded.
46. As a dealership manager, I want hold to be a real proposed action, distinct from having none, so that "keep at the current price" is not the same as "not yet decided".
47. As a dealership manager, I want a new proposed action to replace the previous one, so that only the latest intention is stored.
48. As a dealership manager, I want to clear the proposed action on an aging vehicle, so that I can return it to having no intention recorded.
49. As a dealership manager, I want to see the current proposed action on the vehicle, so that I know what was last decided.
50. As a dealership manager, I want setting, replacing, or clearing a proposed action to be rejected when the vehicle is not aging stock, so that the action is only for vehicles in stock more than 90 days.
51. As a dealership manager, I want the last proposed action to remain visible after the vehicle stocks out, so that I can still see what had been intended.
52. As a dealership manager, I want changing the proposed action after stock-out to be rejected, so that the intention stays as it was when the vehicle left.
53. As a dealership manager, I want any proposed action outside the three named intentions to be rejected, so that the vocabulary stays closed.
54. As a dealership manager, I want to attach one primary image to an in-stock vehicle, so that the list is easier to scan.
55. As a dealership manager, I want to replace that primary image while the vehicle is in stock, so that an updated photo becomes the one picture.
56. As a dealership manager, I want to remove the primary image while the vehicle is in stock, so that the vehicle can go back to having none.
57. As a dealership manager, I want a vehicle with no primary image to remain valid, so that a photo is optional.
58. As a dealership manager, I want a placeholder when a vehicle has no primary image, so that the list still lays out cleanly.
59. As a dealership manager, I want the picture to load from the dashboard's own address, so that I never follow a storage link.
60. As a dealership manager, I want a request for a missing primary image to be not found, so that the screen can keep the placeholder.
61. As a dealership manager, I want attaching, replacing, or removing a primary image to be rejected once the vehicle is stocked out, so that the only correction after departure is the stocked-out on date.
62. As a dealership manager, I want a primary image that was already attached to keep streaming after stock-out, so that the retained record still shows the picture.
63. As a dealership manager, I want a duplicate VIN to say that the VIN is already used, so that I know why the vehicle was not recorded.
64. As a dealership manager, I want an attempt to reopen a stocked-out vehicle to say that the departure stands, so that I know the stock state cannot be reversed.
65. As a dealership manager, I want an unknown vehicle address to say the vehicle was not found, so that I know the link is stale.
66. As a dealership manager, I want a model that does not belong to the selected make to be rejected as a bad filter, so that I do not mistake that mistake for an empty inventory.
67. As a dealership manager, I want the dashboard and the inventory API on one address, so that I open a single local site.

## Implementation Decisions

Routes, bodies, VIN, paging, and errors: `docs/architecture/api-outline.md`. Tables and derived fields: `docs/architecture/database-schema.md`. UI stack: `docs/architecture/inventory-frontend.md`. Keep those aligned with this spec.

- One modular monolith: catalog, vehicle, media, and the HTTP shell. Handlers bind input, call a service, and map errors. Inventory age, aging stock, terminal stocked-out, and VIN uniqueness live in the vehicle module.
- The screen renders `stock_state`, `inventory_age`, and `is_aging` from the API and does not apply a second aging rule.
- Local Compose runs PostgreSQL, unpublished MinIO, and the API on port 8080, serving `/api` and the built screen. No accounts, roles, or TLS. A local UI dev server may proxy `/api`.
- Seed a small make and model catalog, plus demo vehicles at inventory age 0, exactly 90 days, 91 days, and a stocked-out vehicle whose stay was longer than 90 days.

## Testing Decisions

- A good test exercises the HTTP API as a black box. It sends the request a manager's screen would send and asserts the status, the envelope, the derived fields (`stock_state`, `inventory_age`, `is_aging`), and `error_code`. It does not assert SQL, handler structure, or React output.
- The module under test is the API. Catalog reads, vehicle commands, aging stock, and primary-image behavior are all observed through `/api/v1`. PostgreSQL is the real store. The test sets the clock to a fixed UTC date. Image bytes are uploaded and fetched through the image routes. The test does not open MinIO itself and does not drive the screen.
- There are no existing tests to copy. Cover the external behavior of the stories above, including these boundaries: entered today is age 0 and not aging stock; exactly 90 days in stock is not aging stock; 91 days in stock is aging stock; a stocked-out vehicle with a stay over 90 days is not aging stock and its inventory age stops on stocked-out on; a second create with the same VIN conflicts; a correction after stock-out conflicts; clearing stocked-out on conflicts; a proposed action is accepted only for aging stock; a mixed correction that includes proposed action, image, or stock-out fields is rejected; a missing primary image fetches as not found; vehicle JSON has `has_primary_image` and no storage location.
- Do not add a second automated seam in the vehicle logic packages or in the React app for the same rules. The screen is checked by using it against this API.

## Out of Scope

- Returning a stocked-out vehicle to in stock, including a deal that falls through.
- A gallery of more than one image, a pasted image URL, or publishing MinIO to the browser.
- Editing the make and model catalog from the screen, or free-text make and model on the vehicle.
- Prices, executing a price reduction, or executing a wholesale sale. Proposed action records intention only.
- Accounts, roles, multi-tenant access, TLS, and other production hardening.
- Microservices, message queues, caching, GraphQL, and gRPC.
- Internationalization.
- VIN check-digit validation or a decoder.
- A browser test suite or a second test seam that re-asserts the same inventory rules.
- Hard-deleting a vehicle.

## Further Notes

Glossary, ADRs 0001–0009, and the architecture docs are accepted. Application code is not scaffolded yet. The closed picks (no auth, `GET /health`, stock-out date correction, aging-only proposed action, image upload limits, UTC clock, list sort) are recorded in `docs/architecture/api-outline.md`.
