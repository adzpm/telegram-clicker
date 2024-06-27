package math

func CalculateCoinsPerClick(startCoins, level uint64, coinsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := startCoins
	for i := uint64(1); i < level; i++ {
		coinsPerClick = uint64(float64(coinsPerClick) * coinsMultiplier)
	}

	return coinsPerClick
}

func CalculateCoinsPerClickVariant2(startCoins, level, coinsPerLevel uint64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := startCoins
	for i := uint64(1); i < level; i++ {
		coinsPerClick += coinsPerLevel
	}

	return coinsPerClick
}

// CalculateUpgradePrice calculates the price of the upgrade to the next level.
func CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level < 1 {
		return startPrice
	}

	upgradePrice := startPrice * uint64(priceMultiplier)
	for i := uint64(1); i <= level; i++ {
		upgradePrice = uint64(float64(upgradePrice) * priceMultiplier)
	}

	return upgradePrice
}
