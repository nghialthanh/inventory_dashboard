# Frontend ↔ Backend Mapping

Living map of dashboard surfaces → API → hooks. Update a row when a page is wired.

**Status:** planned — fill **Hook** / **Status** columns as implementation proceeds.

| UI surface | User intent | API | Hook (planned) | Status |
|------------|-------------|-----|----------------|--------|
| Inventory list | Filter make / model / aging / stock | `GET /api/v1/vehicles` | `useVehicles(filters)` | planned |
| Make filter options | Choose make | `GET /api/v1/catalog/makes` | `useMakes()` | planned |
| Model filter options | Choose model for make | `GET /api/v1/catalog/models?make_id=` | `useModels(makeId)` | planned |
| Vehicle detail / row expand | Inspect one vehicle | `GET /api/v1/vehicles/:id` | `useVehicle(id)` | planned |
| Create vehicle | Add in-stock vehicle | `POST /api/v1/vehicles` | `useCreateVehicle()` | planned |
| Correct vehicle | Edit allowed fields while in stock | `PATCH /api/v1/vehicles/:id` | `useUpdateVehicle()` | planned |
| Stock out | Record departure day | `POST /api/v1/vehicles/:id/stock-out` | `useStockOutVehicle()` | planned |
| Proposed action | Set/replace action on aging (or eligible) vehicle | `PUT /api/v1/vehicles/:id/proposed-action` | `useSetProposedAction()` | planned |
| Primary image display | Thumbnail / detail image | `GET /api/v1/vehicles/:id/primary-image` | URL helper `vehiclePrimaryImageUrl(id)` | planned |
| Upload / replace image | Set primary image | `PUT /api/v1/vehicles/:id/primary-image` | `useUploadPrimaryImage()` | planned |
| Remove image | Clear primary image | `DELETE /api/v1/vehicles/:id/primary-image` | `useDeletePrimaryImage()` | planned |

## Derived fields the UI should prefer from the API

| Field | Why |
|-------|-----|
| `is_aging` | Single rule (>90 and in stock) — do not recompute with a divergent calendar |
| `inventory_age` | Display age consistently with stocked-out stop date |
| `stock_state` | Derived from `stocked_out_on` |
| `has_primary_image` | Decide placeholder vs `<img src=…>` |
| Embedded `make` / `model` | Avoid extra round-trips on the list |

## Visual distinction (brief)

- **Aging stock:** use `is_aging` for badge/row emphasis.
- **Normal in-stock:** `stock_state === in_stock` && `!is_aging`.
- **Stocked out:** retained rows; not aging; filters may isolate `stocked_out_on` day.
