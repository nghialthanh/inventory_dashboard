# Vehicle inventory

How a dealership records the cars it holds, how long they have been held, and what the manager intends to do about aging stock.

## Language

### Identity

**Vehicle**:
A single physical car the dealership records for sale, identified by its VIN.
_Avoid_: unit, stock item, car

**VIN**:
The business identity of a vehicle. No two vehicles share a VIN.
_Avoid_: stock number, serial number

**In stock**:
A vehicle the dealership currently holds for sale.
_Avoid_: on hand, available, active

**Stocked out**:
A vehicle that has left inventory and is still on record.
_Avoid_: sold, deleted, archived, departed

**Stocked-out on**:
The calendar day a vehicle became stocked out.
_Avoid_: sold date, exit timestamp

**Stock state**:
Whether a vehicle is in stock or stocked out.
_Avoid_: status

**Primary image**:
The single optional picture of a vehicle.
_Avoid_: gallery, photo set, attachment

### Tenure

**Entered on**:
The calendar day a vehicle entered inventory. Inventory age starts on this day.
_Avoid_: received at, stocked date, arrival timestamp

**Inventory age**:
Calendar days from entered on. The entered-on day counts as 0. For a stocked-out vehicle the count stops on stocked-out on.
_Avoid_: age, vehicle age, model year, days on lot

**Aging stock**:
A vehicle that is in stock and whose inventory age is more than 90 days.
_Avoid_: aged unit, old stock, stale inventory

### Catalog

**Make**:
A manufacturer in the vehicle catalog.
_Avoid_: brand, free-text make

**Model**:
A vehicle line that belongs to one make.
_Avoid_: trim, variant, free-text model

### Manager's decision

**Proposed action**:
The manager's current intended action for an aging vehicle: price reduction planned, wholesale planned, or hold. One optional value; a new choice replaces the previous one.
_Avoid_: status, note, action log

**Price reduction planned**:
A proposed action to reduce the vehicle's price.
_Avoid_: discount, markdown

**Wholesale planned**:
A proposed action to sell the vehicle wholesale.
_Avoid_: auction, trade-out

**Hold**:
A proposed action to keep the vehicle for sale at its current price.
_Avoid_: no action, undecided, empty
