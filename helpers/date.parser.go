package helpers

import (
	"fmt"
	"strings"
	"time"
)

func ParseDate(input string) (time.Time, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	t := time.Now()

	switch input {
	case "today":
		return t, nil
	case "yesterday":
		return t.AddDate(0, 0, -1), nil
	case "tomorrow":
		return t.AddDate(0, 0, 1), nil
	}

	allowedFormats := []string{
		"2006-01-02",
		"20060102",
		"02-01-2006",
		"02012006",
	}

	var err error
	for _, f := range allowedFormats {
		t, err = time.Parse(f, input)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("Invalid Date")
}
