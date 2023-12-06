package entity

import (
	"container/heap"
	"sync"
)

// Book represents a trading book that manages buy and sell orders and executes transactions
type Book struct {
	Order         []*Order			// List of all orders in the book
	Transactions  []*Transaction	// List of completed transactions
	OrdersChan    chan *Order 		// Input channel for receiving new orders
	OrdersChanOut chan *Order		// Output channel for sending processed ordes
	Wg            *sync.WaitGroup	// WaitGroup for managing concurrent processing
}

// NewBook creates a new Book instance with the given order channels and WaitGroup
func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

// Trade procesess incoming orders and executes trades based on matching logic
func (b *Book) Trade() {

	// Create queues for buy and sell orders for each asses
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)
	// buyOrders := NewOrderQueue()
	// sellOrders := NewOrderQueue()

	// heap.Init(buyOrders)
	// heap.Init(sellOrders)

	for order := range b.OrdersChan {
		asset := order.Asset.ID

		// Initialize order queues for the asset if not already created
		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}
		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		// Process BUY and SELL ordes separately
		if order.OrderType == "BUY" {
			
			// Add the buy order to the queue
			buyOrders[asset].Push(order)

			// Check for potential matches with texisting SELL orders
			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price {
				sellOrder := sellOrders[asset].Pop().(*Order)

				// Check if the matched SELL order have pending shares
				if sellOrder.PendingShares > 0 {
					// Create a transaction and update orders and book
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.Wg)

					// Update orders and send them to the output channel
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- sellOrder
					b.OrdersChanOut <- order

					// If there are pending shares on the SELL order, add it back to the queue
					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			// Process SELL order. Add the order to the sell queue
			sellOrders[asset].Push(order)

			// Check for potential matches with existing BUY orders
			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)

				// Check if the matched BUY order have pending shares
				if buyOrder.PendingShares > 0 {
					// Create a transaction and update orders and book
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)

					// Update orders and send them to the output channel
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- order

					// If there are pending shares on the BUY order, add it back to the queue
					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

// AddTransaction adds a completed transactio to the book, updating investor asset position and order status
func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	// Pull the current pending shares
	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	// Determine the minimum shares that can be transacted
	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	// Update investor asset positions and order status
	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	// Add the transaction to the book
	b.Transactions = append(b.Transactions, transaction)
}