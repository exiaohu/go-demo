package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		expected  int
		expectErr bool
	}{{
		name:      "Add positive numbers",
		a:         10,
		b:         5,
		expected:  15,
		expectErr: false,
	}, {
		name:      "Add negative numbers",
		a:         -10,
		b:         5,
		expected:  -5,
		expectErr: false,
	}, {
		name:      "Add zero",
		a:         0,
		b:         0,
		expected:  0,
		expectErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Add(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		expected  int
		expectErr bool
	}{{
		name:      "Subtract positive numbers",
		a:         10,
		b:         5,
		expected:  5,
		expectErr: false,
	}, {
		name:      "Subtract negative numbers",
		a:         -10,
		b:         5,
		expected:  -15,
		expectErr: false,
	}, {
		name:      "Subtract zero",
		a:         0,
		b:         0,
		expected:  0,
		expectErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Subtract(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		expected  int
		expectErr bool
	}{{
		name:      "Multiply positive numbers",
		a:         10,
		b:         5,
		expected:  50,
		expectErr: false,
	}, {
		name:      "Multiply negative numbers",
		a:         -10,
		b:         5,
		expected:  -50,
		expectErr: false,
	}, {
		name:      "Multiply by zero",
		a:         10,
		b:         0,
		expected:  0,
		expectErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Multiply(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a         int
		b         int
		expected  int
		expectErr bool
	}{{
		name:      "Divide positive numbers",
		a:         10,
		b:         5,
		expected:  2,
		expectErr: false,
	}, {
		name:      "Divide negative numbers",
		a:         -10,
		b:         5,
		expected:  -2,
		expectErr: false,
	}, {
		name:      "Divide by zero",
		a:         10,
		b:         0,
		expected:  0,
		expectErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
