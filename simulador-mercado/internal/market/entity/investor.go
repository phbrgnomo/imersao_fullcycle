package entity

// Structure - composite data type that groups together variables
// Investor represents an investor on the market.
type Investor struct {
	ID			string
	Name        string
	AssetPosition []*InvestorAssetPosition
}

// Constructor Function
// NewInvestor creates a new instance of the Investor entity.
func NewInvestor(id string) *Investor {
	return &Investor{
		ID: id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

// Methods associated with Investor struct
// AddAssetPosition adds a new asset position to the investor's portfolio.
func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition){
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

// UpdateAssetPosition updates an existing asset position or adds a new one.
func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) {
	// Retrieve the existing asset position
	assetPosition := i.GetAssetPosition(assetID)

	// If the asset position exists, update the shares. If the asset position doesn't exist, create a new one
	if assetPosition == nil  {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qtdShares))
	} else {
		assetPosition.Shares += qtdShares
	}

}

// GetAssetPosition retrieves the asset position based on the assetID.
func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}

	return nil
}


// InvestorAssetPosition represents an asset position on the market.
type InvestorAssetPosition struct {
	AssetID		string
	Shares		int
}

// Constructor Function
// NewInvestorAssetPosition creates a new instance of the InvestorAssetPosition entity.
func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares: shares,
	}
}