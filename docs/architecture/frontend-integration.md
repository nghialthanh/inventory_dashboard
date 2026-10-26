# Frontend integration

`frontend/` owns presentation. The Go API owns aging, stock-out, and VIN rules.

Shape list and detail DTOs to the screen: embedded make and model, `is_aging`, `inventory_age`, `stock_state`. Change React only when the backend cannot express the feature, and then reuse the existing components and layout.

```
pages → thin hooks / adapters → /api/v1 → HTTP adapter → services interface → logic → repository → PostgreSQL
```

Calls stay on the Compose origin. After a mutation, update the React Query cache from the response when it is complete; refetch when derived fields are not.

Surfaces: [frontend-backend-mapping.md](frontend-backend-mapping.md). Stack: [inventory-frontend.md](inventory-frontend.md).
