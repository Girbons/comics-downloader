package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIssuesRange(t *testing.T) {
	tt := []struct {
		input         string
		expectedStart float64
		expectedEnd   float64
		hasError      bool
	}{
		{"1-1", 1, 1, false},
		{"1-5", 1, 5, false},
		{"3-9", 3, 9, false},
		{"3.1-9.5", 3.1, 9.5, false},
		{"3.-9.5", 3, 9.5, false},
		{"12-123", 12, 123, false},
		{"0-0", 0, 0, true},
		{"0-1", 0, 1, false},
		{"1-0", 0, 0, true},
		{"2-1", 0, 0, true},
		{"1", 0, 0, true},
		{"wrong range", 0, 0, true},
		// Volume.issue format tests
		{"4.78-4.99", 4.78, 4.99, false},     // Nightwing V4 #078-099
		{"4.078-4.099", 4.078, 4.099, false}, // Same with leading zeros
		{"1.10-1.20", 1.10, 1.20, false},     // Volume 1, issues 10-20
		{"2.01-2.05", 2.01, 2.05, false},     // Volume 2, issues 01-05
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			start, end, err := ParseIssuesRange(tc.input)
			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)

			if tc.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
