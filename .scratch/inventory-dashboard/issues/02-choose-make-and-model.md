# 02: Choose make and model from the catalog

**Decision:** [ADR-0002](../../../docs/adr/0002-model-catalog.md). This file is the acceptance checklist for that decision.

**Blocked by:** 01 Open the empty dashboard

**Status:** ready-for-agent

- [x] The make filter lists the seeded makes
- [x] Choosing a make limits the model filter to models of that make
- [x] Models can also be read without choosing a make
- [x] Each model belongs to one make
- [x] The dashboard and the API offer no way to create, rename, or delete a make or model
- [x] Catalog responses use the shared envelope and snake_case fields
- [x] The catalog reads are verified through the HTTP API against PostgreSQL
