package entity

import (
	"container/heap"
	"sync"
)

// Book respresents an order book for trading
type Book struct {
	Order				[]*Order 		// List of orders in the book
	Transactions		[]*Transaction	// List of transaction in the book
	OrdersChan			chan *Order 	// Input - Orders sent by Kafka
	OrdersChanOut		chan *Order		// Output - Orders processed and sent to Kafka
	Wg					*sync.WaitGroup // Synchronization primitive. Coordinate the orders execution
}

// NewBook creates a new instance of the Book
func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         	[]*Order{},
		Transactions:   []*Transaction{},
		OrdersChan:     orderChan,
		OrdersChanOut:  orderChanOut,
		Wg:             wg,
	}
}

// Trade processess incoming orders and execute trades
func (b *Book)  Trade() {
	buyOrders := NewOrderQueue() // Create buy orders queue
	sellOrders := NewOrderQueue() // Create sell orders queue

	heap.Init(buyOrders) // Initialize the heap sorting algo on the buy queue
	heap.Init(sellOrders) // Initialize the heap sorting algo on the sell queue
	
	// Loop checking if there is any new order on the Orders Channel
	for order := range b.OrdersChan {
		
		// Check order type
		if order.OrderType == "BUY" {
			
			// Add order to the buy queue
			buyOrders.Push(order)

			// Check if there is an order on the sell queue with a price that matches the buy order
			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				
				// Remove the sell orderd from the queue
				sellOrder := sellOrders.Orders.Pop().(*Order)

				// Check if the sell order matched still have pending shares to sell
				if sellOrder.PendingShares > 0 {
					
					// Creates a new transaction with the sell order information
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					
					// Add the transaction to the book's transaction list
					b.AddTransaction(transaction, b.Wg)
					
					// Append the transaction to the sell order's transactions list
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					
					// Append the transaction to the buy order's transactions list
					order.Transactions = append(order.Transactions, transaction)
					
					// Send the order information exit channel to be sent to Kafka
					b.OrdersChanOut <- sellOrder
					b.OrdersChanOut <- order
					
					// If there is still shares remaining on the sell order, sent the order back to the queue
					if sellOrder.PendingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			// Add sell order to the sell queue
			sellOrders.Push(order)

			// Check if there is a buy order on the buy queue with a price that matches the sell order
			for buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				// Remove the buy order from the queue
				buyOrder := buyOrders.Orders.Pop().(*Order)
		
				// Check if the buy order matched still has pending shares to buy
				if buyOrder.PendingShares > 0 {
					// Creates a new transaction with the buy order information
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
		
					// Add the transaction to the book's transaction list
					b.AddTransaction(transaction, b.Wg)
		
					// Append the transaction to the buy order's transactions list
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
		
					// Append the transaction to the sell order's transactions list
					order.Transactions = append(order.Transactions, transaction)
		
					// Send the order information exit channel to be sent to Kafka
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- order
		
					// If there is still shares remaining on the buy order, send the order back to the queue
					if buyOrder.PendingShares > 0 {
						buyOrders.Push(buyOrder)
					}
		}

	}
}

// AddTransaction updates the book's state based on a completed transaction.
func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	// Extract pending shares from both the selling and buying orders
	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	// Determine the minimum shares to process in the transaction
	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	// Update asset positions and pending shares for the selling and buying orders
	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.SellingOrder.PendingShares -= minShares
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.BuyingOrder.PendingShares -= minShares
	
	// Calculate and set the total amount for the transaction
	transaction.Total = float64(transaction.Shares) * transaction.BuyingOrder.Price
	// transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	// Check if the orders are fully executed and update its status
	// transaction.CloseBuyTransaction()
	// transaction.CloseSellTransaction()
	if transaction.BuyingOrder.PendingShares == 0{
		transaction.BuyingOrder.Status = "CLOSED"
	}
	if transaction.SellingOrder.PendingShares == 0{
		transaction.SellingOrder.Status = "CLOSED"
	}

}