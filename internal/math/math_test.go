package math

import (
	"testing"

	config "github.com/adzpm/telegram-clicker/internal/config"
)

func TestCalculateCoinsPerClick(t *testing.T) {
	testCases := map[string]struct {
		startCoins      uint64
		level           uint64
		coinsMultiplier float64
		investors       uint64
		expected        uint64
	}{
		"level 0 / start 1":                {1, 0, 2.1, 0, 0},
		"level 1 / start 1":                {1, 1, 2.1, 0, 1},
		"level 2 / start 1":                {1, 2, 2.1, 0, 2},
		"level 3 / start 1":                {1, 3, 2.1, 0, 4},
		"level 4 / start 1":                {1, 4, 2.1, 0, 9},
		"level 5 / start 1":                {1, 5, 2.1, 0, 19},
		"level 6 / start 1":                {1, 6, 2.1, 0, 40},
		"level 7 / start 1":                {1, 7, 2.1, 0, 85},
		"level 8 / start 1":                {1, 8, 2.1, 0, 180},
		"level 0 / start 1 / investors 50": {1, 0, 2.1, 50, 0},
		"level 1 / start 1 / investors 50": {1, 1, 2.1, 50, 2},
		"level 2 / start 1 / investors 50": {1, 2, 2.1, 50, 4},
		"level 3 / start 1 / investors 50": {1, 3, 2.1, 50, 8},
		"level 4 / start 1 / investors 50": {1, 4, 2.1, 50, 18},
		"level 5 / start 1 / investors 50": {1, 5, 2.1, 50, 38},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var (
				mth = New(&config.GameVariables{
					EarnedCoinsForInvestor: 5000000,
					PercentsForInvestor:    0.02,
				})
				investorsMultiplier = mth.CalculateInvestorsMultiplier(tc.investors)
				result              = mth.CalculateGeometricCoinsPerClick(
					tc.startCoins,
					tc.level,
					tc.coinsMultiplier,
					investorsMultiplier,
				)
			)

			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestCalculateCoinsPerClickV2(t *testing.T) {
	testCases := map[string]struct {
		startCoins uint64
		level      uint64
		investors  uint64
		expected   uint64
	}{
		"case 1": {1, 0, 0, 0},
		"case 2": {1, 1, 0, 1},
		"case 3": {1, 2, 0, 2},
		"case 4": {1, 3, 0, 3},
		"case 5": {1, 4, 0, 4},
		"case 6": {1, 5, 0, 5},

		"case 1-1": {1, 0, 50, 0},
		"case 1-2": {1, 1, 50, 2},
		"case 1-3": {1, 2, 50, 4},
		"case 1-4": {1, 3, 50, 6},
		"case 1-5": {1, 4, 50, 8},
		"case 1-6": {1, 5, 50, 10},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var (
				mth = New(&config.GameVariables{
					EarnedCoinsForInvestor: 5000000,
					PercentsForInvestor:    0.02,
				})
				investorsMultiplier = mth.CalculateInvestorsMultiplier(tc.investors)
				result              = mth.CalculateAlgebraCoinsPerClick(
					tc.startCoins,
					tc.level,
					investorsMultiplier,
				)
			)

			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestCalculateUpgradePrice(t *testing.T) {
	testCases := map[string]struct {
		startPrice      uint64
		level           uint64
		priceMultiplier float64
		expected        uint64
	}{
		"level 0 / start 100": {100, 0, 1.5, 100},
		"level 1 / start 100": {100, 1, 1.5, 150},
		"level 2 / start 100": {100, 2, 1.5, 225},
		"level 3 / start 100": {100, 3, 1.5, 337},
		"level 4 / start 100": {100, 4, 1.5, 506},
		"level 5 / start 100": {100, 5, 1.5, 759},
		"level 6 / start 100": {100, 6, 1.5, 1139},
		"level 7 / start 100": {100, 7, 1.5, 1708},
		"level 8 / start 100": {100, 8, 1.5, 2562},
		"case 1":              {10, 0, 3, 10},
		"case 2":              {10, 1, 3, 30},
		"case 3":              {10, 2, 3, 90},
		"case 4":              {10, 3, 3, 270},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var (
				mth = New(&config.GameVariables{
					EarnedCoinsForInvestor: 5000000,
					PercentsForInvestor:    0.02,
				})
				result = mth.CalculateUpgradePrice(tc.startPrice, tc.level, tc.priceMultiplier)
			)

			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}
