# Inventory frontend

React + TypeScript dashboard in `frontend/`. The shell route is `/`, rendered by `InventoryPage`. Compose serves the production build from the API. A local `npm run dev` proxies `/api` and `/health` to port 8080; that proxy is documented in [frontend/README.md](../../frontend/README.md).

| Layer | Choice |
|-------|--------|
| UI | React 19.3, TypeScript 5.9 |
| Server state | TanStack Query 5 |
| Styling | Tailwind 4 |
| Routing | React Router 7. The shell route is `/` |
| API | Client to `/api/v1`, envelope unwrapped |

Compose serves the production build from the API ([ADR-0009](../adr/0009-one-origin-under-api.md)). The dev proxy is in place.

## Screens

The inventory list is the workspace. Create, correct, stock-out, and proposed action are dialogs. Labels use [GLOSSARY.md](../../GLOSSARY.md).

| Component | Role |
|-----------|------|
| `InventoryFilters` | Make, model, aging, stock state |
| `CorrectVehicleDialog` | Correct VIN, model, and entered on while in stock |
| `VehicleTable` or `VehicleCards` | List, with aging emphasis |
| `AgingBadge` | Driven by `is_aging` |
| `ProposedActionControl` | Proposed-action command |
| `StockOutDialog` | Calendar day and confirm |
