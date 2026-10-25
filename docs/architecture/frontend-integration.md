# Frontend integration

`frontend/` owns presentation. The Go API owns aging, stock-out, VIN, and image rules.

Shape list and detail DTOs to the screen: embedded make and model, `is_aging`, `inventory_age`, `stock_state`, `has_primary_image`. Change React only when the backend cannot express the feature, and then reuse the existing components and layout.

```
pages → thin hooks / adapters → /api/v1 → HTTP adapter → services interface → logic → repository → PostgreSQL and private MinIO
```

Calls and primary-image loads stay on the Compose origin. After a mutation, update the React Query cache from the response when it is complete; refetch when derived fields are not.

Surfaces: [frontend-backend-mapping.md](frontend-backend-mapping.md). Stack: [inventory-frontend.md](inventory-frontend.md).
