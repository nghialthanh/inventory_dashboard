# Stock state is the presence of stocked-out on

A vehicle is in stock when stocked-out on is empty, and stocked out when that date is set. Stock state is derived from the date. There is no status column that can disagree with it.

## Considered options

- A status value plus a date. Those two can diverge: stocked out with no date, or in stock with a departure date. Every filter would then need a tie-break.
