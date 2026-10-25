# Inventory Frontend

Notes for the React + TypeScript dashboard (`frontend/`).

**Status:** planned until the app is scaffolded — update routes and stack pins when code exists.

## Stack (target)

| Layer | Choice |
|-------|--------|
| UI library | React + TypeScript |
| Data | TanStack Query |
| Styling | Tailwind (+ small shared component set) |
| Routing | File or declarative router — one style only |
| API | `fetch`/client → `/api/v1`, envelope unwrap |
| Images | Same-origin `GET …/primary-image` |

Compose serves the production build from the API process ([ADR-0009](../adr/0009-one-origin-under-api.md)). Local `npm run dev` may proxy `/api` to the API for faster UI iteration — document the proxy in `frontend/README` when added.

## Information architecture

Keep the first useful screen focused on the brief:

1. **Inventory list** — primary workspace (filters + table/cards).
2. **Dialogs / drawers** — create, stock-out, proposed action, image — not separate mini-apps.
3. Avoid a marketing landing page; this is an ops dashboard.

## Copy

Use [GLOSSARY.md](../../GLOSSARY.md) terms in labels and empty states:

- Aging stock, proposed action, stocked out, entered on, inventory age, make, model, VIN, primary image.

Avoid: status (for stock state), brand, aged unit, gallery.

## Components (suggested)

| Component | Role |
|-----------|------|
| `InventoryFilters` | Make, model, aging toggle / stock state |
| `VehicleTable` / `VehicleCards` | List with aging emphasis |
| `AgingBadge` | Driven by `is_aging` |
| `ProposedActionControl` | Sets action via command mutation |
| `PrimaryImage` | Placeholder or streamed image |
| `StockOutDialog` | Calendar day + confirm |

## i18n

Not required by the brief. If added later, both languages must stay in sync — do not hardcode one-off strings in JSX without a plan.
