package math

func CalculateCoinsPerClick(startCoins, level uint64, coinsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := float64(startCoins)
	for i := uint64(1); i < level; i++ {
		coinsPerClick *= coinsMultiplier
	}

	return uint64(coinsPerClick)
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
