package dto

// TradeInput represents the input data for a trade order
type TradeInput struct {
	OrderID       string  `json:"order_id"`
	InvestorID    string  `json:"investor_id"`
	AssetID       string  `json:"asset_id"`
	CurrentShares int     `json:"current_shares"`
	Shares        int     `json:"shares"`
	Price         float64 `json:"price"`
	OrderType     string  `json:"order_type"`
}

// OrderOutput represents the output data for a trade order
type OrderOutput struct {
	OrderID            string               `json:"order_id"`
	InvestorID         string               `json:"investor_id"`
	AssetID            string               `json:"asset_id"`
	OrderType          string               `json:"order_type"`
	Status             string               `json:"status"`
	Partial            int                  `json:"partial"`
	Shares             int                  `json:"shares"`
	TransactionsOutput []*TransactionOutput `json:"transactions"`
}

// TransactionOutput represents the output data for a transaction
type TransactionOutput struct {
	TransactionID string  `json:"transaction_id"`
	BuyerID       string  `json:"buyer_id"`
	SellerID      string  `json:"seller_id"`
	AssetID       string  `json:"asset_id"`
	Price         float64 `json:"price"`
	Shares        int     `json:"shares"`
}