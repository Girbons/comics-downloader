package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseIssuesRange the range of issues.
// Format [start-end] or [volume.issue-volume.issue].
// Examples: "1-5", "3.1-9.5", "4.78-4.99" (volume 4, issues 78-99)
func ParseIssuesRange(rng string) (float64, float64, error) {
	values := strings.Split(rng, "-")
	if len(values) != 2 {
		return 0, 0, errors.New("wrong range format")
	}

	startRange, err := parseVolumeIssue(values[0])
	if err != nil {
		return 0, 0, fmt.Errorf("wrong the start range value: %v", err)
	}

	endRange, err := parseVolumeIssue(values[1])
	if err != nil {
		return 0, 0, fmt.Errorf("wrong the end range value: %v", err)
	}

	if endRange == 0 {
		return 0, 0, errors.New("the end range value must not be zero")
	}

	if startRange > endRange {
		return 0, 0, errors.New("the start range value must be less or equal to the end range value")
	}

	return startRange, endRange, nil
}

// parseVolumeIssue parses a value that may be in format "volume.issue" or just a simple number.
// For backwards compatibility and simplicity, it just parses as a regular float.
// Examples: "4.78" -> 4.78, "4.099" -> 4.099, "3.1" -> 3.1, "5" -> 5.0
func parseVolumeIssue(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
