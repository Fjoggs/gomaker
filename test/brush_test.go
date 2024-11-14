package test

import (
	"testing"

	"gomaker/internal/brush"
)

func TestIsBrush(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"// brush 1", true},
		{
			"( 104 400 176 ) ( 112 400 192 ) ( 104 272 176 ) common/caulk 32 0 0 0.5 0.5 134217728 0 0",
			false,
		},
		{"// Entity 0", false},
		{
			"( 96 80 192 ) ( 240 80 128 ) ( 240 80 192 ) testmap/test_texture 461.2879333496 22.0878295898 -26.5999984741 0.2808699906 0.280872494 134217728 0 0",
			false,
		},
		{"// Brush 1", true},
		{"// entity 1", false},
	}
	for _, test := range tests {
		value := brush.IsBrush(test.input)
		if value != test.expected {
			t.Errorf("Expected %v got %s for %v", value, test.input, test)
		}
	}
}
