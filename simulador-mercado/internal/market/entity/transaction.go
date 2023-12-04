package entity

import (
	"time"
	"github.com/google/uuid"
)

// Transaction represents a transaction in the market.
type Transaction struct {
	ID             	string
	SellingOrder	*Order
	BuyingOrder     *Order
	Shares			int
	Price           float64
	Total			float64
	DateTime		time.Time
}

// NewTransaction creates a new instance of the Transaction entity.
func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		ID: 			uuid.New().String(),
		SellingOrder: 	sellingOrder,
		BuyingOrder: 	buyingOrder,
		Shares: 		shares,
		Price: 			price,
		Total: 			total,
		DateTime: 		time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) float64 {
	t.Total = float64(t.Shares) * t.Price
}

func (t *Transaction) CloseBuyTransaction() *Order {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) CloseSellTransaction() *Order {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) AddBuyOrderPendingShares(shares int) {
	t.SellingOrder.PendingShares += shares
}

func (t *Transaction) AddSellOrderPendingShares(shares int) {
	t.BuyingOrder.PendingShares += shares
}