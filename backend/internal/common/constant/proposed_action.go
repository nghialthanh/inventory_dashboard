package constant

// Proposed action values stored on a vehicle.
const (
	ProposedActionPriceReductionPlanned = "price_reduction_planned"
	ProposedActionWholesalePlanned      = "wholesale_planned"
	ProposedActionHold                  = "hold"
)

// KnownProposedActions is the closed set of proposed-action intentions.
var KnownProposedActions = map[string]struct{}{
	ProposedActionPriceReductionPlanned: {},
	ProposedActionWholesalePlanned:      {},
	ProposedActionHold:                  {},
}
