package math

import "testing"

func TestCalculateCoinsPerClick(t *testing.T) {
	testCases := map[string]struct {
		startCoins      uint64
		level           uint64
		coinsMultiplier float64

		expected uint64
	}{
		"level 0 / start 1": {1, 0, 2.1, 0},
		"level 1 / start 1": {1, 1, 2.1, 1},
		"level 2 / start 1": {1, 2, 2.1, 2},
		"level 3 / start 1": {1, 3, 2.1, 4},
		"level 4 / start 1": {1, 4, 2.1, 8},
		"level 5 / start 1": {1, 5, 2.1, 16},
		"level 6 / start 1": {1, 6, 2.1, 33},
		"level 7 / start 1": {1, 7, 2.1, 69},
		"level 8 / start 1": {1, 8, 2.1, 144},
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
		"level 0 / start 100": {100, 0, 1.5, 100},
		"level 1 / start 100": {100, 1, 1.5, 150},
		"level 2 / start 100": {100, 2, 1.5, 225},
		"level 3 / start 100": {100, 3, 1.5, 337},
		"level 4 / start 100": {100, 4, 1.5, 505},
		"level 5 / start 100": {100, 5, 1.5, 757},
		"level 6 / start 100": {100, 6, 1.5, 1135},
		"level 7 / start 100": {100, 7, 1.5, 1702},
		"level 8 / start 100": {100, 8, 1.5, 2553},

		"case 1": {100, 1, 2.2, 440},
		"case 2": {100, 2, 2.2, 968},
		"case 3": {100, 3, 2.2, 2129},
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
