package util

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	errMustBePositive = errors.New("must be >= 0")
	errInvalidUnit    = errors.New("invalid time unit")
)

type Duration time.Duration

const Day = 24 * Duration(time.Hour)

func ParseDuration(s string) (Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if s == "0" {
		return 0, nil
	}

	i := strings.IndexFunc(s, func(r rune) bool { return !(r >= '0' && r <= '9') && r != '-' })
	prefix := s
	suffix := ""
	if i >= 0 {
		prefix = s[:i]
		suffix = s[i:]
	}

	prefixInt, err := strconv.Atoi(prefix)
	if err != nil {
		return 0, err
	}
	if prefixInt < 0 {
		return 0, errMustBePositive
	}
	n := Duration(prefixInt)

	switch suffix {
	case "", "d":
		return n * Day, nil
	case "w":
		return n * 7 * Day, nil
	default:
		return 0, errInvalidUnit
	}
}

func (d Duration) Ago() *time.Time {
	if d <= 0 {
		return nil
	}
	t := time.Now().Add(-time.Duration(d))
	return &t
}
