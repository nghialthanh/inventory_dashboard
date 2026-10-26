# AI Collaboration — Step by Step

How I used GenAI (Cursor) to design and build the **inventory dashboard** assessment, and how I kept ownership of architecture and quality.

**Related:** [SYSTEM_DESIGN.md §7](SYSTEM_DESIGN.md#7-genai-assistance-in-the-design-phase) (design phase) · [README — AI Collaboration Narrative](../README.md#ai-collaboration-narrative) (summary for submission)

---

## 0. Ground rules I set before prompting

I treated the AI as a **fast junior engineer**, not as the decision-maker.

| Rule | Why |
|------|-----|
| One vocabulary | Only [GLOSSARY.md](../GLOSSARY.md) terms (`aging stock`, `stocked out`, `proposed action`) |
| Docs before large code | Spec → ADRs → API outline → then implementation |
| One slice per chat | One issue or invariant at a time |
| Always verify | `go test ./...`, Compose smoke, doc sync |
| Reject ADR conflicts | Even if the model “helps” with a shortcut |

Tooling: **Cursor** (chat + agents + rules under `.cursor/rules/`).

---

## Phase A — Understand the brief (design)

### Step 1 — Turn the challenge into a single scenario

**What I asked the AI to do**

- Restate the Keyloop-style brief as one problem: dealership inventory, aging stock (> 90 days), proposed action.
- Extract user stories and acceptance boundaries (age 0 / 90 / 91, stocked-out not aging).

**What I kept / changed**

- Kept modular monolith + one origin (API + UI on one port).
- Rejected ideas that split into many microservices for an assessment.

**Output in the repo**

- `.scratch/inventory-dashboard/spec.md`
- Issue breakdown under `.scratch/inventory-dashboard/issues/`

---

### Step 2 — Lock the domain language

**What I asked**

- Draft a glossary from the spec; flag synonyms to avoid.

**What I reviewed**

- Every term has a preferred word and an “Avoid” list.
- No “status” column for stock; stock state is derived.

**Output**

- [GLOSSARY.md](../GLOSSARY.md)

---

### Step 3 — Capture hard decisions as ADRs

**What I asked**

- Propose ADRs for irreversible choices (retain stocked-out rows, catalog via seed, UUID in URL, separate commands, one origin). Object storage was considered and later dropped ([ADR-0010](adr/0010-no-object-storage.md)).

**How I refined**

- Edited each ADR until the **trade-off and consequence** were clear.
- Used ADRs as a stop-list: if a later AI suggestion contradicted an ADR, I rejected it.

**Output**

- [docs/adr/](adr/) (0001–0011)

---

### Step 4 — Shape the system (layers, API, schema)

**What I asked**

- Architecture diagram and package layout: HTTP adapter → services interfaces → logic → repositories.
- API outline (envelope, filters, aging, commands).
- Database schema matching derived fields (no `is_aging` column).

**What I decided myself**

- Backend is the **assessed** layer; React UI is a bonus client.
- Aging threshold **> 90** owned only by API logic.

**Output**

- [docs/architecture/](architecture/)
- Later packaged as [SYSTEM_DESIGN.md](SYSTEM_DESIGN.md) for the submission

---

### Step 5 — Install persistent AI guidance

**What I asked / set up**

- Cursor rules so every future prompt follows the same architecture and glossary.

**Why**

- Stops the model from reinventing layers or inventing new domain words mid-feature.

**Output**

- `.cursor/rules/` (`architecture.mdc`, `backend.mdc`, `api.mdc`, `testing.mdc`, …)
- [AGENTS.md](../AGENTS.md)

---

## Phase B — Implement the backend (chosen service layer)

I implemented features in the same order as the issue list, with the AI writing first drafts and me verifying each slice.

### Step 6 — Bootstrap runtime

**Prompt pattern**

> “Compose: Postgres + Go API serving `/api` and static UI on one origin. Health outside `/api/v1`. Follow ADR-0009.”

**Verify**

- `docker compose up --build`
- `GET /health`

---

### Step 7 — Catalog (seeded makes / models)

**Prompt pattern**

> “Read-only catalog from seed. List makes and models. Model belongs to one make. No free-text make/model.”

**Verify**

- Seed SQL + HTTP catalog tests
- ADR-0002 respected

---

### Step 8 — Record and list vehicles

**Prompt pattern**

> “POST vehicle with VIN, model_id, entered_on. Derive stock_state, inventory_age, is_aging. List with paging and filters. Fixed UTC clock in tests.”

**Verify**

- Age 0 not aging; 90 not aging; 91 aging
- Envelope + `error_code` assertions

---

### Step 9 — Aging filter and stock-out

**Prompt pattern**

> “`aging=true` only. Stock-out is a separate route. Row stays; VIN never reused; clearing stocked_out_on conflicts.”

**Verify**

- Issue 04 + 06 style HTTP tests
- ADR-0001, ADR-0004, ADR-0007

---

### Step 10 — Correct in stock; reject after stock-out

**Prompt pattern**

> “PATCH only vin / model_id / entered_on while in stock. UUID in path stays. After stock-out → CONFLICT.”

**Verify**

- Correction and post-stock-out conflict tests

---

### Step 11 — Proposed action (aging only)

**Prompt pattern**

> “PUT proposed-action only when is_aging. Values: price_reduction_planned | wholesale_planned | hold | empty to clear.”

**Verify**

- Conflict when not aging; value remains readable after later stock-out

---

### Step 12 — Primary image (removed)

An optional image slice was drafted against MinIO. [ADR-0010](adr/0010-no-object-storage.md) removed object storage, the image routes, and the `primary_image_key` column. Compose is Postgres and the API only.

---

### Step 13 — Observability (zap + trace_id)

**Prompt pattern**

> “Use zap. Put trace_id on context (X-Trace-ID). Pass through handler → logic → repository. Access log + internal error logs.”

**Verify**

- Middleware unit tests
- Full suite still green: `cd backend && go test ./...`
- Docs updated in SYSTEM_DESIGN §5

---

## Phase C — Frontend (bonus, not the graded layer)

### Step 14 — Thin UI on the real API

**Prompt pattern**

> “React inventory page. Filters and badges from API fields only — do not recompute aging. Hooks map 1:1 to routes.”

**Verify**

- Manual Compose / Vite smoke
- Mapping table in [frontend-backend-mapping.md](architecture/frontend-backend-mapping.md)

I still submit **backend** as the chosen service layer; the UI is a live client. cURL examples stub the client for reviewers who do not open the UI: [api-curl-examples.md](api-curl-examples.md).

---

## Phase D — Submission artifacts

### Step 15 — System design document

**Prompt pattern**

> “Produce SDD covering diagram, components, data flow, tech justifications, observability, GenAI in design — from existing architecture docs.”

**What I checked**

- Matches Keyloop Part 1 checklist
- Honest “implemented vs near-term” observability

**Output**

- [SYSTEM_DESIGN.md](SYSTEM_DESIGN.md)

### Step 16 — README narrative + this document

**Prompt pattern**

> “README: how to build, run, test; AI Collaboration Narrative; link to full step-by-step.”

**Output**

- [README.md](../README.md)
- This file

---

## How I verified AI output (every phase)

```text
1. Read the diff against GLOSSARY + ADRs
2. Run:  cd backend && go test ./...
3. Smoke: docker compose / curl (and UI if the change touched it)
4. Update owning doc in the same change if behavior moved
5. Only then mark the issue / slice done
```

| Check | Pass means |
|-------|------------|
| Domain | Age / VIN / stock-out / proposed-action rules match tests |
| Structure | Handlers have no inventory rules; logic has no SQL |
| Contract | JSON snake_case envelope matches api-outline |
| Language | No new synonyms in code or UI copy |

---

## Example prompts that worked well

**Good (constrained)**

> “Implement stock-out as POST `/api/v1/vehicles/:id/stock-out` only. Follow ADR-0007. Add HTTP tests with fixed clock. Do not put stocked_out_on on PATCH.”

**Bad (too open)**

> “Add vehicle status and make the API nice.”  
> → invites a status column and mixed commands; I would reject and re-prompt with glossary + ADR references.

---

## What GenAI did vs what I owned

| GenAI accelerated | I owned |
|-------------------|---------|
| First drafts of ADRs, API outline, SQL, handlers, tests | Final ADR wording and tech choice (Go, Postgres, modular monolith). MinIO was dropped in ADR-0010 |
| Boilerplate and test scaffolding | Pass/fail on acceptance criteria |
| Doc packaging (SDD, cURL examples) | Whether FE is “bonus” vs graded layer |
| Observability wiring (zap, middleware) | What must be logged vs kept quiet (4xx) |

---

## Timeline view (for a 5–10 minute interview)

1. Brief → glossary + issues  
2. ADRs + architecture  
3. Cursor rules so AI stays on rails  
4. Backend slices 01→10 with tests each time  
5. Optional UI as client  
6. SDD + README AI narrative + this step-by-step  

**One sentence for the interviewer:**  
*I used GenAI to move faster on drafts and scaffolding; I used glossary, ADRs, and automated tests so the final system still reflects my architectural and domain decisions.*
