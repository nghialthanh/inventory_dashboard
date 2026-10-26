# Vehicles reference a make and model catalog

Make and model are chosen from a seeded catalog so a vehicle cannot be stored under a mistyped or impossible pair. A model belongs to one make. The vehicle references the model. Daily use does not edit the catalog; a missing model is added to the seed.

## Considered options

- Free-text make and model on the vehicle. Simpler to enter, and it breaks “choose from a list.”
- A model table with the make name copied onto each row. The same make is then spelled many ways.

Accepted in [02](../../.scratch/inventory-dashboard/issues/02-choose-make-and-model.md).
