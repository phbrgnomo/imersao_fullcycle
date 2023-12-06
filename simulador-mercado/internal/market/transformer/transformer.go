package transformer

import (
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/dto"
	"github.com/phbrgnomo/imersao_fullcycle/simulador-mercado/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	// Create a new Asset based on the input data
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)

	// Create a new Investor based on the input data
	investor := entity.NewInvestor(input.InvestorID)

	// Create a new Order based on the input data
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)
	
	// If there are existing shares, create an InvestorAssetPosition and add it to the Investor
	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}
	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutput {

	// Create a new OrderOutput based on the Order
	output := &dto.OrderOutput{
		OrderID:    order.ID,
		InvestorID: order.Investor.ID,
		AssetID:    order.Asset.ID,
		OrderType:  order.OrderType,
		Status:     order.Status,
		Partial:    order.PendingShares,
		Shares:     order.Shares,
	}

	// Create TransactionOutput objects for each Transaction in the Order
	var transactionsOutput []*dto.TransactionOutput
	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: t.ID,
			BuyerID:       t.BuyingOrder.Investor.ID,
			SellerID:      t.SellingOrder.Investor.ID,
			AssetID:       t.SellingOrder.Asset.ID,
			Price:         t.Price,
			Shares:        t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	// Add the TransactionOutput objects to the OrderOutput
	output.TransactionsOutput = transactionsOutput
	return output
}