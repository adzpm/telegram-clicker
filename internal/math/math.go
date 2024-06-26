package math

// CalculateCoinsPerClick calculates the amount of coins per click.
// for example, if startCoins = 1, level = 1, coinsMultiplier = 1.5
// then the result will be 1.5 on the first level.
// if level = 2, then the result will be 2.25.
// if level = 3, then the result will be 3.375.
// if level = 4, then the result will be 5.0625.
// if level = 5, then the result will be 7.59375.
// if level = 6, then the result will be 11.390625.
// if level = 7, then the result will be 17.0859375.
// and so on.
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

// CalculateUpgradePrice calculates the price of the upgrade to the next level.
func CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level == 0 {
		return startPrice
	}

	upgradePrice := startPrice * uint64(priceMultiplier)
	for i := uint64(1); i <= level; i++ {
		upgradePrice = uint64(float64(upgradePrice) * priceMultiplier)
	}

	return upgradePrice
}
