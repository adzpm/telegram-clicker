package math

import (
	config "github.com/adzpm/telegram-clicker/internal/config"
)

type (
	Math struct {
		config *config.GameVariables
	}
)

func New(cfg *config.GameVariables) *Math { return &Math{cfg} }

// CalculateGeometricCoinsPerClick calculates the number of coins per click based on the level.
func (m *Math) CalculateGeometricCoinsPerClick(startCoins, level uint64, coinsMultiplier, investorsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := float64(startCoins)
	for i := uint64(1); i < level; i++ {
		coinsPerClick *= coinsMultiplier
	}

	return uint64(coinsPerClick * investorsMultiplier)
}

// CalculateAlgebraCoinsPerClick calculates the number of coins per click based on the level.
func (m *Math) CalculateAlgebraCoinsPerClick(startCoins, level uint64, investorsMultiplier float64) uint64 {
	if level == 0 {
		return 0
	}

	coinsPerClick := startCoins
	for i := uint64(1); i < level; i++ {
		coinsPerClick += startCoins
	}

	return uint64(float64(coinsPerClick) * investorsMultiplier)
}

// CalculateUpgradePrice calculates the price of the upgrade to the next level.
func (m *Math) CalculateUpgradePrice(startPrice, level uint64, priceMultiplier float64) uint64 {
	if level < 1 {
		return startPrice
	}

	upgradePrice := float64(startPrice)
	for i := uint64(1); i <= level; i++ {
		upgradePrice *= priceMultiplier
	}

	return uint64(upgradePrice)
}

// CalculateInvestorsCount calculates the number of investors based on the earned coins.
func (m *Math) CalculateInvestorsCount(earnedCoins uint64) uint64 {
	return earnedCoins / m.config.EarnedCoinsForInvestor
}

// CalculateInvestorsMultiplier calculates the multiplier for the investors.
func (m *Math) CalculateInvestorsMultiplier(investors uint64) float64 {
	return 1 + float64(investors)*m.config.PercentsForInvestor
}

func (m *Math) GetGameVariables() *config.GameVariables {
	return m.config
}
