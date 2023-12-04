package entity

// OrderQueue represents a priority queue (heap) for orders.
type OrderQueue struct {
	Orders []*Order
}

// Less returns true if the order at position i has a lower price than the order at position j.
// It implements the heap.Interface and is used for sorting orders based on their prices.
func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

// Swap swaps the positions of orders at indices i and j in the queue.
// It is used during heap operations to maintain the order of the priority queue.
func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

// Len returns the number of orders in the priority queue.
// It is used during heap operations to determine the size of the priority queue.
func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

// Push adds a new order to the priority queue.
// It is used during heap operations to insert new orders.
func (oq *OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order))
}

// Pop removes and returns the order with the highest priority (lowest price).
// It is used during heap operations to extract the top order from the priority queue.
func (oq *OrderQueue) Pop() interface{} {
	old := oq.Orders
	n := len(old)
	item := old[n-1]
	oq.Orders = old[0 : n-1]
	return item
}

// NewOrderQueue creates a new instance of the OrderQueue.
func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}