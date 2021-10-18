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
		{"0-1", 0, 0, true},
		{"1-0", 0, 0, true},
		{"2-1", 0, 0, true},
		{"1", 0, 0, true},
		{"wrong range", 0, 0, true},
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
