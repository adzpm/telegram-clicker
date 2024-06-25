package math

import "testing"

func TestCalculateCoinsPerClick(t *testing.T) {
	testCases := map[string]struct {
		startCoins      uint64
		level           uint64
		coinsMultiplier float64

		expected uint64
	}{
		"level 0 / start 2": {2, 0, 2, 0},
		"level 1 / start 2": {2, 1, 2, 2},
		"level 2 / start 2": {2, 2, 2, 4},
		"level 3 / start 2": {2, 3, 2, 8},
		"level 4 / start 2": {2, 4, 2, 12},
		"level 5 / start 2": {2, 5, 2, 16},
		"level 6 / start 2": {2, 6, 2, 20},
		"level 7 / start 2": {2, 7, 2, 24},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := CalculateCoinsPerClick(tc.startCoins, tc.level, tc.coinsMultiplier)
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

		expected uint64
	}{
		"level 0 / start 100": {100, 0, 2, 100},
		"level 1 / start 100": {100, 1, 2, 200},
		"level 2 / start 100": {100, 2, 2, 400},
		"level 3 / start 100": {100, 3, 2, 600},
		"level 4 / start 100": {100, 4, 2, 800},
		"level 5 / start 100": {100, 5, 2, 1000},
		"level 6 / start 100": {100, 6, 2, 1200},
		"level 7 / start 100": {100, 7, 2, 1400},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := CalculateUpgradePrice(tc.startPrice, tc.level, tc.priceMultiplier)
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}
