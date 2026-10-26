# Stock-out and proposed action are separate commands

PATCH corrects only the fields the current stock state allows. Departure and proposed action each have their own route, so each rule lives with one command.

## Considered options

- One PATCH that accepts stock state, dates, and action together. The client can send a mix the rules forbid, and the handler has to pick which field failed.

Accepted in [05](../../.scratch/inventory-dashboard/issues/05-correct-an-in-stock-vehicle.md), [06](../../.scratch/inventory-dashboard/issues/06-stock-out-a-vehicle.md), and [09](../../.scratch/inventory-dashboard/issues/09-record-a-proposed-action.md).
