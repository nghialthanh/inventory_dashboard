# Retain vehicles after they leave inventory

Managers need to see which vehicles stocked out on a given day, including today. A vehicle that leaves stays on record as stocked out, with the calendar day it left. Aging stock is only vehicles still in stock. Inventory age for a stocked-out vehicle stops on that day.

## Considered options

- Drop the vehicle on departure. Smaller model, and the stocked-out-on-a-day filter is impossible.
- Keep a separate departure event log. More machinery than a state and a date on the vehicle.
