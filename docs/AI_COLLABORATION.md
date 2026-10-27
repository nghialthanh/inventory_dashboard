# AI Collaboration — Step by Step

How I used GenAI (Cursor) on the **inventory dashboard**. The model drafted. I decided. Nothing was implemented until the decision, the issue, and the tests were in place.

**Related:** [SYSTEM_DESIGN.md §7](SYSTEM_DESIGN.md#7-genai-assistance-in-the-design-phase) · [README — AI Collaboration Narrative](../README.md#ai-collaboration-narrative)

---

## The loop

Every feature went through the same loop. Later work (logs, metrics, partitions, cache) used the same loop. It did not skip ahead of the domain.

1. **Set the rails** — rules, architecture, technical choices, glossary.
2. **Grill the feature** — the model asks until I accept each decision.
3. **File the work** — spec and issues under `.scratch/`, with a checklist I watch.
4. **Build key features first** — catalog, record, list, aging, stock-out, correct, proposed action.
5. **Add operability** — logs, trace id, metrics, health, shutdown.
6. **Improve the running system** — identity, monthly partitions, Redis before the database, read replica.
7. **Test, then check the box** — HTTP tests against Postgres; the checklist moves only when they pass.
8. **Keep the record** — ADR, architecture doc, and API outline change in the same slice.
9. **Hang the UI on the finished API** — the screen reads derived fields; it does not invent a second aging rule.
10. **Review the diff myself** — glossary, ADRs, and the issue checklist are the gate.

---

## 1. Set the rails before any feature

I started by fixing the words, the shape, and the stack. Later chats inherited that through Cursor rules, so the model could not invent a new layer or a new domain word halfway through a feature.

| Rail | What I locked | Where it lives |
|------|----------------|----------------|
| Language | Glossary terms only: aging stock, stocked out, proposed action, entered on, vehicle address | [GLOSSARY.md](../GLOSSARY.md) |
| Shape | Modular monolith. HTTP adapter binds and maps errors. Services are interfaces. Logic owns rules. Repositories own SQL. `server.go` wires them | [docs/architecture/](architecture/) · `.cursor/rules/architecture.mdc` |
| Technical choices | Go, Echo, PostgreSQL, pgx, Docker Compose, one origin on `:8080` | [SYSTEM_DESIGN.md §4](SYSTEM_DESIGN.md#4-technology-choices-and-justifications) · [ADR-0009](adr/0009-one-origin-under-api.md) |
| Invariants | In stock when `stocked_out_on` is null. Aging means in stock and inventory age > 90. Stocked out stays. VIN is unique and immutable. Stock-out and proposed action are separate routes | `.cursor/rules/architecture.mdc` · [docs/adr/](adr/) |
| Agent behavior | One feature per chat. Update the owning doc in the same change. Tests cover the aging boundaries | `.cursor/rules/` · [AGENTS.md](../AGENTS.md) |

Hard choices were written as ADRs before the code that depends on them. Each ADR states the trade-off and the consequence. A later suggestion that contradicted an ADR was rejected.

---

## 2. Grill each feature before writing code

For each feature I had the model **grill** me: a design tree, one round of questions at a time, a recommended answer on every question, and no implementation until I had answered.

A round looks like this:

- The question is something only I can decide (where the aging rule lives, whether stock-out shares a route with correct, which column is the partition key).
- The model recommends an answer and words the question so “yes” accepts it.
- Facts in the repo are the model’s job to look up. Decisions stay mine.
- The next round opens only after this round is settled.
- Code starts when the frontier is empty: nothing important left assumed.

Examples I settled this way:

| Feature | Decision the grill forced |
|---------|---------------------------|
| Stock state | Derived from `stocked_out_on`. No status column ([ADR-0005](adr/0005-stock-state-from-stocked-out-on.md)) |
| Commands | Stock-out and proposed action are their own routes ([ADR-0007](adr/0007-separate-commands.md)) |
| Images | Drafted, grilled, then withdrawn. No object storage ([ADR-0010](adr/0010-no-object-storage.md)) |
| Vehicle address | App-minted snowflake; `created_at` matches the mint instant ([ADR-0013](adr/0013-snowflake-vehicle-addresses.md)) |
| Partition key | `created_at`, because `entered_on` is still correctable while the vehicle is in stock |
| Catalog cache | Redis first. The replica runs only on a miss. Vehicle lists stay uncached ([ADR-0014](adr/0014-catalog-cache-and-read-replica.md)) |

---

## 3. Turn the accepted feature into issues, and watch `.scratch`

After a grill I had the model write the work down. I tracked it. The model did not close a line on its own.

Layout ([docs/agents/issue-tracker.md](agents/issue-tracker.md)):

- One feature directory: `.scratch/<feature-slug>/`
- The spec: `.scratch/<feature-slug>/spec.md`
- One file per issue: `.scratch/<feature-slug>/issues/<NN>-<slug>.md`
- A `Status:` line (`ready-for-agent`, `ready-for-human`, `withdrawn`, …)
- A checklist of acceptance lines inside the issue
- The ADR index points at the issue that carried that decision: [docs/adr/README.md](adr/README.md)

The inventory feature lives at [.scratch/inventory-dashboard/](../.scratch/inventory-dashboard/). Issues `01`–`09` are the domain. Issue `10` (primary image) is `withdrawn`. Later improvements have their own specs (integer addresses, build-for-the-future).

I read the checklist before and after the slice. A box is checked when the HTTP test for that line is green and the diff still matches the ADR.

---

## 4. Build the key features first

Operability and scale waited until the inventory behavior was real. The issue order is the build order:

| Issue | Feature |
|-------|---------|
| [01](../.scratch/inventory-dashboard/issues/01-open-the-empty-dashboard.md) | Empty dashboard, one origin, health |
| [02](../.scratch/inventory-dashboard/issues/02-choose-make-and-model.md) | Seeded catalog; model belongs to one make |
| [03](../.scratch/inventory-dashboard/issues/03-record-a-vehicle.md) | Record and list; derived stock state, inventory age, aging |
| [04](../.scratch/inventory-dashboard/issues/04-see-aging-stock.md) | Aging filter; 0 and 90 are not aging; 91 is |
| [05](../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md) | Correct model and entered on while in stock |
| [06](../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md) | Stock-out keeps the row; VIN is never reused |
| [07](../.scratch/inventory-dashboard/issues/07-review-departures.md) | Filter stocked-out departures |
| [08](../.scratch/inventory-dashboard/issues/08-reject-corrections-after-stock-out.md) | Correct after stock-out conflicts |
| [09](../.scratch/inventory-dashboard/issues/09-record-a-proposed-action.md) | Proposed action only while aging |
| [10](../.scratch/inventory-dashboard/issues/10-attach-a-primary-image.md) | Withdrawn after the grill |

Each slice was one chat: the issue file, the ADR it cites, and “do not start the next issue.”

---

## 5. Add logs and metrics after the domain works

With record, list, stock-out, and proposed action in place, I opened a separate slice for operability. Mixing it into feature `03` would have hidden domain mistakes inside logging diffs.

| Signal | What landed |
|--------|-------------|
| Logging | zap. Access log per request. Domain 4xx stays quiet. Unexpected SQL and `INTERNAL_ERROR` are logged |
| Trace | `X-Trace-ID` on the context, echoed on the response, attached from handler through logic to repository |
| Metrics | `GET /metrics`: HTTP count and latency, catalog cache hit / miss / error, database query count by role |
| Process | `GET /health`. Graceful shutdown drains in-flight HTTP, then closes Redis and the pools |

Detail: [SYSTEM_DESIGN.md §5](SYSTEM_DESIGN.md#5-observability-strategy).

---

## 6. Improve the system on its own issues

Scale and identity were new features, so they were grilled, written up, and built after the key paths existed.

| Improvement | Decision | Result |
|-------------|----------|--------|
| Catalog addresses | Integer identity ([ADR-0011](adr/0011-integer-addresses.md)) | Make and model ids are integers |
| VIN | Immutable after record ([ADR-0012](adr/0012-vin-immutable.md)) | Correct cannot change the VIN |
| Vehicle address and storage | Snowflake plus monthly partitions ([ADR-0013](adr/0013-snowflake-vehicle-addresses.md)) | `vehicles` is RANGE-partitioned on `created_at` (`vehicles_2026_01` … `vehicles_2027_12`). Point reads pass the month decoded from the id. Global VIN uniqueness stays on non-partitioned `vehicle_vins` |
| Catalog read path | Redis first, replica on miss ([ADR-0014](adr/0014-catalog-cache-and-read-replica.md)) | Logic reads Redis. A hit returns and skips SQL. A miss, a Redis error, or a bad payload reads the replica and then fills the key (~1h). Vehicle lists stay on the primary |
| Proof under load | Bulk seed, one script | `scripts/seed_bulk_vehicles.sql` loads 510 vehicles into every month partition |

---

## 7. Test, then move the checklist

Tests are the close of a slice, and they are also how I checked the model.

- HTTP tests hit the real envelope and PostgreSQL. The clock is fixed in UTC so age 0, 90, and 91 are exact.
- Assertions cover stock state, inventory age, the aging flag, duplicate VIN, terminal stock-out, and command isolation.
- Command for the suite: `cd backend && go test ./...`
- A UI or Compose change also gets a smoke pass (`docker compose`, curl, or the screen).
- I check the issue box after that run. A red test keeps the line open.

The testing rule for the model is `.cursor/rules/testing.mdc`.

---

## 8. Update the owning doc in the same change

A behavior change updates the doc that owns it before the slice is done:

| Change | Doc |
|--------|-----|
| Domain word | `GLOSSARY.md` |
| Hard trade-off | `docs/adr/` |
| Endpoint or envelope | `docs/architecture/api-outline.md` |
| Table, column, partition | `docs/architecture/database-schema.md` |
| Page, hook, route | `docs/architecture/frontend-backend-mapping.md` |
| This collaboration story | this file and the README narrative |

---

## 9. Put the UI on the API after the contract is stable

The assessed layer is the backend. The React screen is a client.

- Filters, badges, and aging come from API fields. The UI does not recompute the 90-day rule.
- Hooks map one-to-one onto routes.
- The mapping table is [frontend-backend-mapping.md](architecture/frontend-backend-mapping.md).
- Reviewers who skip the UI can follow [api-curl-examples.md](api-curl-examples.md).

---

## 10. Review the diff before the slice is done

```text
1. Read the diff against the glossary, the ADR, and the issue checklist
2. Run:  cd backend && go test ./...
3. Smoke Compose or the screen when the change is visible
4. Update the owning doc if behavior moved
5. Check the .scratch boxes, and only then leave the chat
```

| Check | Pass means |
|-------|------------|
| Domain | Age, VIN, stock-out, and proposed action match the issue lines |
| Structure | Handlers have no inventory rules. Logic has no SQL |
| Contract | Snake_case envelope matches the API outline |
| Language | Copy and code use glossary terms |
| Operability | Logs and metrics describe the path that was just built, including Redis hit versus miss |
| Scale | A new month partition, cache key, or read path is named in the ADR and the schema doc |

---

## What GenAI did vs what I owned

| GenAI accelerated | I owned |
|-------------------|---------|
| Grill questions and a recommended answer | Every yes / no, and the ADR wording |
| Issue files and checklists under `.scratch/` | Status, and which box is actually done |
| First drafts of handlers, SQL, and tests | Whether the slice matches the issue |
| Logs, metrics, partition DDL, cache-aside | The order: domain first, then operability, then scale |
| Doc packaging | Backend as the assessed layer; UI as a client |

**One sentence for the interviewer:**  
*I set the rules and the architecture first, grilled each feature until the decisions were mine, tracked the work as checklists in `.scratch`, built the inventory features before logs and scale, and checked a box only after the tests passed.*
