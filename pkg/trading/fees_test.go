package trading

import (
	"testing"
)

func TestCalculateFee(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		feeBPS   int
		expected float64
	}{
		{"10 bp on 1000", 1000, 10, 1.0},
		{"20 bp on 1000", 1000, 20, 2.0},
		{"10 bp on 50000", 50000, 10, 50.0},
		{"0 bp fee", 1000, 0, 0},
		{"100 bp on 100", 100, 100, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateFee(tt.amount, 0, 0, tt.feeBPS)
			if result != tt.expected {
				t.Errorf("CalculateFee(%f, 0, 0, %d) = %f, want %f", tt.amount, tt.feeBPS, result, tt.expected)
			}
		})
	}
}

func TestIsMakerOrder(t *testing.T) {
	if !IsMakerOrder("GTC") {
		t.Error("GTC should be maker")
	}
	if IsMakerOrder("IOC") {
		t.Error("IOC should not be maker")
	}
	if IsMakerOrder("FOK") {
		t.Error("FOK should not be maker")
	}
}
