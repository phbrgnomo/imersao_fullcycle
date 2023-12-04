package entity

// Asset represents an asset in the market.
type Asset struct {
	ID            	string
	Name          	string
	MarketVolume	int
}

// Constructor Function
// NewAsset creates a new instance of the Asset entity.
func NewAsset(id string, name string, marketVolume int) *Asset {
	return &Asset{
		ID: id,
		Name: name,
		MarketVolume: marketVolume,
	}
}
