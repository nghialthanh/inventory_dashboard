# Frontend Integration Rule

`frontend/` is the **official frontend foundation** and the **source of truth for the presentation layer**.

**Related:** [Frontend ↔ Backend mapping](frontend-backend-mapping.md) · [Inventory frontend](inventory-frontend.md)

## Principle

Integrate the React UI with the Golang backend by adapting **backend** APIs, contracts, and business logic to match what the UI needs. Prefer evolving the backend to support the screen — not redesigning the screen to match raw storage — unless data integrity or ADR invariants require otherwise.

## Requirements

1. **Preserve presentation** — Keep the UI library, component structure, routing, layouts, and UX stable once scaffolded.
2. **No unnecessary FE refactors** — Do not replace the UI kit mid-feature without a critical reason.
3. **Reuse components** — Prefer existing inventory/list/filter components over new ones.
4. **Backend fills gaps** — Missing filters, derived aging fields, or command endpoints → extend Go first.
5. **Contract follows UI** — If the list needs `is_aging` / embedded `make`, add them to the list DTO rather than forcing N+1 client calls.
6. **Minimal FE edits** — Change React only when the feature cannot be supported from the backend; keep diffs tiny.
7. **Same design language** — New controls match existing dashboard patterns.
8. **Same origin** — Call `/api/v1/...` and image routes on the Compose host; no CORS and no MinIO URLs ([ADR-0009](../adr/0009-one-origin-under-api.md), [ADR-0008](../adr/0008-api-streams-images.md)).

## Integration pattern (preferred)

```
frontend pages (stable JSX)
        ↓ thin hooks / adapters (map API → view models)
        ↓
   /api/v1  (Go handlers + DTOs shaped for the dashboard)
        ↓
   domain services / PostgreSQL (+ MinIO behind media logic)
```

- Keep page components consuming stable view-model shapes (`adapters.ts`).
- Prefer dashboard-friendly payloads from the backend; use FE mappers when domain enums must stay internal.
- After successful mutations, patch React Query cache from the mutation payload when practical; refetch when derived aggregates are incomplete.

## Goal

The final system keeps a clear inventory UI while the Golang backend owns aging, stock-out, VIN, and image rules.
