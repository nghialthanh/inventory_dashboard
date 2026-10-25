package aging

import (
	"time"

	"inventory_dashboard/internal/common/constant"
)

// UTCDate is the UTC calendar day of t.
func UTCDate(t time.Time) time.Time {
	year, month, day := t.UTC().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// InventoryAge is the number of calendar days from entered on to the end of the stay.
// The entered-on day is 0. A stocked-out vehicle stops on stocked-out on.
func InventoryAge(enteredOn time.Time, stockedOutOn *time.Time, today time.Time) int {
	end := UTCDate(today)
	if stockedOutOn != nil {
		end = UTCDate(*stockedOutOn)
	}
	return int(end.Sub(UTCDate(enteredOn)).Hours() / 24)
}

// StockState is in stock while stocked-out on is empty.
func StockState(stockedOutOn *time.Time) string {
	if stockedOutOn == nil {
		return "in_stock"
	}
	return "stocked_out"
}

// IsAging reports aging stock: in stock, and inventory age greater than 90 days.
func IsAging(stockedOutOn *time.Time, inventoryAge int) bool {
	return stockedOutOn == nil && inventoryAge > constant.AgingStockDays
}
