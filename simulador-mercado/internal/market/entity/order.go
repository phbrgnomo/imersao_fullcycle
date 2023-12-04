package entity

// Order represents a buy or sell order in the market.
type Order struct {
	ID 				string
	Investor		*Investor
	Asset       	*Asset
	Shares			int
	PendingShares	int
	Price           float64
	OrderType		string
	Status          string
	Transactions	[]*Transaction
}

// Constructor Function
// NewOrder creates a new instance of the Order entity.
func NewOrder(orderID string, investor *Investor, asset *Asset, shares int, price float64, orderType) {
	return &Order{
		ID: 			orderID,
		Investor: 		investor,
		Asset: 			asset,
		Shares: 		shares,
		PendingShares: 	shares,
		Price: 			price,
		OrderType: 		orderType,
		Status:         "OPEN",
		Transaction: 	[]*Transaction{},
	}
}