package shipping

import (
	"testing"
	"math"
)

func almostEqual(a, b float64) bool {
	// For floating point comparisons in tests
	return math.Abs(a-b) < 0.0001
}

// Comprehensive test for CalculateShippingFee (v2) specification
func TestCalculateShippingFee_V2(t *testing.T) {
	testCases := []struct {
		name        string
		weight      float64
		zone        string
		insured     bool
		expectedFee float64
		expectError bool
	}{
		// --- Weight Invalid: Too Small and Too Large
		{"Weight = 0 (invalid)", 0, "Domestic", false, 0, true},   // Lower boundary, invalid
		{"Weight < 0 (invalid)", -5, "Express", true, 0, true},
		{"Weight > 50 (invalid)", 51, "Domestic", false, 0, true}, // Upper boundary, invalid
		{"Weight = 100 (invalid)", 100, "International", true, 0, true},

		// --- Zone Invalid
		{"Zone invalid (Local)", 5, "Local", false, 0, true},
		{"Zone invalid (empty)", 10, "", false, 0, true},
		{"Zone invalid (case)", 10, "domestic", true, 0, true}, // Should be case-sensitive

		// --- Standard Package: 0 < w <= 10
		{"Standard, Domestic, not insured", 1, "Domestic", false, 5.0, false},
		{"Standard, Domestic, insured", 1, "Domestic", true, 5.0 * 1.015, false},
		{"Standard, International, not insured", 5, "International", false, 20.0, false},
		{"Standard, International, insured", 10, "International", true, 20.0 * 1.015, false},
		{"Standard, Express, not insured", 10, "Express", false, 30.0, false},
		{"Standard, Express, insured", 10, "Express", true, 30.0 * 1.015, false},

		// --- Standard/Heavy Boundary
		{"Boundary: weight just below heavy", 10, "Domestic", false, 5.0, false},
		{"Boundary: weight just above standard", 10.01, "Domestic", false, 5.0 + 7.5, false},

		// --- Heavy Package: 10 < w <= 50
		{"Heavy, Domestic, not insured", 20, "Domestic", false, 5.0 + 7.5, false},
		{"Heavy, Domestic, insured", 20, "Domestic", true, (5.0 + 7.5) * 1.015, false},
		{"Heavy, International, not insured", 15, "International", false, 20.0 + 7.5, false},
		{"Heavy, International, insured", 49.99, "International", true, (20.0 + 7.5) * 1.015, false},
		{"Heavy, Express, not insured", 11, "Express", false, 30.0 + 7.5, false},
		{"Heavy, Express, insured", 50, "Express", true, (30.0 + 7.5) * 1.015, false},

		// --- Upper Boundary Cases
		{"Boundary: weight = 50 (max valid)", 50, "Domestic", false, 5.0 + 7.5, false},
		{"Boundary: weight just above max", 50.01, "Domestic", false, 0, true},

		// --- Lower Boundary for Standard
		{"Boundary: weight just above 0", 0.01, "Domestic", false, 5.0, false},
		// insured at the lowest possible valid weight
		{"Boundary: weight just above 0, insured", 0.01, "Domestic", true, 5.0 * 1.015, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFee(tc.weight, tc.zone, tc.insured)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil (fee: %v)", fee)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error, got: %v", err)
				}
				if !almostEqual(fee, tc.expectedFee) {
					t.Errorf("Expected fee %.4f, got %.4f", tc.expectedFee, fee)
				}
			}
		})
	}
}
