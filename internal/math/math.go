package math

const (
	OneInvestorCoins   = 1000000 // one investor for 1M coins
	InvestorMultiplier = 0.02    // 2% for each investor
)

func CalculateCoinsPerClick(startCoins, level uint64, coinsMultiplier, investorsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := float64(startCoins)
	for i := uint64(1); i < level; i++ {
		coinsPerClick *= coinsMultiplier
	}

	return uint64(coinsPerClick * investorsMultiplier)
}

// CalculateUpgradePrice calculates the price of the upgrade to the next level.
func CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level < 1 {
		return startPrice
	}

	upgradePrice := float64(startPrice)
	for i := uint64(1); i <= level; i++ {
		upgradePrice *= priceMultiplier
	}

	return uint64(upgradePrice)
}

func CalculateInvestorsCount(earnedCoins uint64) uint64 {
	return earnedCoins / OneInvestorCoins
}

func CalculateInvestorsMultiplier(investors uint64) float64 {
	return 1 + float64(investors)*InvestorMultiplier
}
