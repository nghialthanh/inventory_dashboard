# Frontend ↔ Backend Mapping

Living map of dashboard surfaces → API → hooks. Update a row when a page is wired.

**Status:** planned — fill **Hook** / **Status** columns as implementation proceeds.

| UI surface | User intent | API | Hook (planned) | Status |
|------------|-------------|-----|----------------|--------|
| Inventory list | Filter by make and model, and page | `GET /api/v1/vehicles` | `useVehicles(filters)` | implemented |
| Inventory list | Filter to aging stock only | `GET /api/v1/vehicles?aging=true` | `useVehicles(filters)` | implemented |
| Inventory list | Filter by stock state | `GET /api/v1/vehicles` | `useVehicles(filters)` | planned |
| Make filter options | Choose make | `GET /api/v1/catalog/makes` | `useMakes()` | implemented |
| Model filter options | Choose model for make | `GET /api/v1/catalog/models?make_id=` | `useModels(makeId)` | implemented |
| Vehicle detail / row expand | Inspect one vehicle | `GET /api/v1/vehicles/:id` | `useVehicle(id)` | implemented |
| Create vehicle | Add in-stock vehicle | `POST /api/v1/vehicles` | `useCreateVehicle()` | implemented |
| Correct vehicle | Edit allowed fields while in stock | `PATCH /api/v1/vehicles/:id` | `useUpdateVehicle()` | planned |
| Stock out | Record departure day | `POST /api/v1/vehicles/:id/stock-out` | `useStockOutVehicle()` | planned |
| Proposed action | Set, replace, or clear on aging stock | `PUT /api/v1/vehicles/:id/proposed-action` | `useSetProposedAction()` | planned |
| Primary image display | Thumbnail / detail image | `GET /api/v1/vehicles/:id/primary-image` | URL helper `vehiclePrimaryImageUrl(id)` | planned |
| Upload / replace image | Set primary image | `PUT /api/v1/vehicles/:id/primary-image` | `useUploadPrimaryImage()` | planned |
| Remove image | Clear primary image | `DELETE /api/v1/vehicles/:id/primary-image` | `useDeletePrimaryImage()` | planned |

Render `is_aging`, `inventory_age`, `stock_state`, `has_primary_image`, and embedded `make` / `model` from the list payload in [api-outline.md](api-outline.md). Aging rows follow `is_aging`. Stocked-out rows stay distinct and are not aging stock.
