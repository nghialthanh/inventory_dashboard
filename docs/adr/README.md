# Decisions

Each file in this folder is one hard-to-reverse choice. The linked issue is the acceptance checklist that carried that choice into the dashboard. Read the decision here. Use the issue only to see how that decision was checked.

| ADR | Decision | Accepted in |
|-----|----------|-------------|
| [0001](0001-retain-stocked-out-vehicles.md) | A vehicle that leaves stays on record, with the day it stocked out | [06](../../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md), [07](../../.scratch/inventory-dashboard/issues/07-review-departures.md) |
| [0002](0002-model-catalog.md) | Make and model come from a seeded catalog | [02](../../.scratch/inventory-dashboard/issues/02-choose-make-and-model.md) |
| [0003](0003-primary-image-in-minio.md) | Superseded by [0010](0010-no-object-storage.md) | [10](../../.scratch/inventory-dashboard/issues/10-attach-a-primary-image.md) (withdrawn) |
| [0004](0004-stocked-out-is-terminal.md) | Stocked out does not return to in stock, and the VIN is never reused | [06](../../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md), [07](../../.scratch/inventory-dashboard/issues/07-review-departures.md), [08](../../.scratch/inventory-dashboard/issues/08-reject-corrections-after-stock-out.md) |
| [0005](0005-stock-state-from-stocked-out-on.md) | Stock state is whether stocked-out on is set | [03](../../.scratch/inventory-dashboard/issues/03-record-a-vehicle.md), [06](../../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md) |
| [0006](0006-uuid-vehicle-identity.md) | Superseded in part by [0011](0011-integer-addresses.md). The URL uses an address, not the VIN | [03](../../.scratch/inventory-dashboard/issues/03-record-a-vehicle.md), [05](../../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md), [08](../../.scratch/inventory-dashboard/issues/08-reject-corrections-after-stock-out.md) |
| [0007](0007-separate-commands.md) | Stock-out and proposed action are their own commands | [05](../../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md), [06](../../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md), [09](../../.scratch/inventory-dashboard/issues/09-record-a-proposed-action.md) |
| [0008](0008-api-streams-images.md) | Superseded by [0010](0010-no-object-storage.md) | [10](../../.scratch/inventory-dashboard/issues/10-attach-a-primary-image.md) (withdrawn) |
| [0009](0009-one-origin-under-api.md) | The API serves the screen on one origin | [01](../../.scratch/inventory-dashboard/issues/01-open-the-empty-dashboard.md) |
| [0010](0010-no-object-storage.md) | No vehicle images and no object store | [10](../../.scratch/inventory-dashboard/issues/10-attach-a-primary-image.md) (withdrawn) |
| [0011](0011-integer-addresses.md) | Make and model addresses are identity integers. The vehicle address is a bigint from one PostgreSQL sequence | [spec](../../.scratch/integer-addresses/spec.md) |

[04 See aging stock](../../.scratch/inventory-dashboard/issues/04-see-aging-stock.md) checks the aging-stock rule (in stock and inventory age greater than 90). That rule is an inventory invariant, recorded in the architecture rules and the glossary, and it has no ADR of its own.
