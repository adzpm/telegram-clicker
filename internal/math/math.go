package math

func CalculateCoinsPerClick(startCoins, level uint64, coinsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	// example:
	//
	// currentCoinsPerClick = 1, coinsMultiplier = 2.5, level = 1
	// 1 * 2.5 * 1 = 2 (rounded from 2.5)
	//
	// currentCoinsPerClick = 1, coinsMultiplier = 2.5, level = 2
	// 1 * 2.5 * 2 = 5
	//
	// currentCoinsPerClick = 1, coinsMultiplier = 2.5, level = 3
	// 1 * 2.5 * 3 = 7 (rounded from 7.5)
	//
	// currentCoinsPerClick = 1, coinsMultiplier = 2.5, level = 4
	// 1 * 2.5 * 4 = 10

	return uint64(float64(startCoins) * coinsMultiplier * float64(level))
}

func CalculateCoinsPerMinute(startCoins, level uint64, coinsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	return CalculateCoinsPerClick(startCoins, level, coinsMultiplier)
}

func CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level == 0 {
		return startPrice
	}

	// example:
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 1
	// 1 * 2.5 * 1 = 2 (rounded from 2.5)
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 2
	// 1 * 2.5 * 2 = 5
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 3
	// 1 * 2.5 * 3 = 7 (rounded from 7.5)
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 4
	// 1 * 2.5 * 4 = 10
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 5
	// 1 * 2.5 * 5 = 12 (rounded from 12.5)
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 6
	// 1 * 2.5 * 6 = 15
	//
	// startPrice = 1, priceMultiplier = 2.5, level = 7
	// 1 * 2.5 * 7 = 17 (rounded from 17.5)

	return uint64(float64(startPrice) * priceMultiplier * float64(level))
}
