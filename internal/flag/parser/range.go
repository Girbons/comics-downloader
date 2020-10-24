package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseIssuesRange the range of issues.
// Format [start-end].
func ParseIssuesRange(rng string) (int, int, error) {
	values := strings.Split(rng, "-")
	if len(values) != 2 {
		return 0, 0, errors.New("wrong range format")
	}

	startRange, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, fmt.Errorf("wrong the start range value: %w", err)
	}

	if startRange == 0 {
		return 0, 0, errors.New("the start range value must not be zero")
	}

	endRange, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, fmt.Errorf("wrong the end range value: %w", err)
	}

	if endRange == 0 {
		return 0, 0, errors.New("the end range value must not be zero")
	}

	if startRange > endRange {
		return 0, 0, errors.New("the start range value must be less or equal to the end range value")
	}

	return startRange, endRange, nil
}
