package math

func CalculateCoinsPerClick(startCoins, level uint64, coinsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	if level == 1 {
		return startCoins
	}

	// example:
	//
	// startCoins = 1, coinsMultiplier = 2, level = 2
	// then current coins perClick:
	// 1 * 2 * 2 = 4
	// startCoins = 1, coinsMultiplier = 2, level = 3
	// 1 * 2 * 3 = 6

	return uint64(float64(startCoins) * coinsMultiplier * float64(level))
}

func CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level == 0 {
		return startPrice
	}

	// example:
	//
	// startPrice = 10, priceMultiplier = 2.5, level = 1
	// then next level:
	// 10 * 2.5 * 2 = 50
	// startPrice = 10, priceMultiplier = 1.5, level = 2
	// 10 * 2.5 * 3 = 75

	return uint64(float64(startPrice) * priceMultiplier * float64(level+1))
}
