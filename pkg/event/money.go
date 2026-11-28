package event

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money is an amount of US dollars in whole micro-dollars (1 USD = 1_000_000).
// It is never a float. In JSON it is a string with exactly six decimals, e.g. "12.345600".
type Money int64

const microPerDollar = 1_000_000

// ParseMoney parses a non-negative decimal dollar amount with at most six decimals,
// such as "12", "0.5" or "12.345600". Signs, exponents and spaces are rejected.
func ParseMoney(s string) (Money, error) {
	if s == "" {
		return 0, errors.New("empty amount")
	}
	whole, frac, hasDot := strings.Cut(s, ".")
	if whole == "" || (hasDot && frac == "") {
		return 0, fmt.Errorf("invalid amount %q", s)
	}
	if len(frac) > 6 {
		return 0, fmt.Errorf("amount %q has more than 6 decimals", s)
	}
	if !allDigits(whole) || !allDigits(frac) {
		return 0, fmt.Errorf("invalid amount %q", s)
	}

	w, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || w > math.MaxInt64/microPerDollar-1 {
		return 0, fmt.Errorf("amount %q is too large", s)
	}
	f := int64(0)
	if frac != "" {
		f, _ = strconv.ParseInt(frac+strings.Repeat("0", 6-len(frac)), 10, 64)
	}
	return Money(w*microPerDollar + f), nil
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// String formats the amount with six decimals, e.g. "12.345600".
func (m Money) String() string {
	sign := ""
	v := int64(m)
	if v < 0 {
		sign = "-"
		v = -v
	}
	return fmt.Sprintf("%s%d.%06d", sign, v/microPerDollar, v%microPerDollar)
}

// MarshalJSON writes the amount as a JSON string with six decimals.
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.String() + `"`), nil
}

// UnmarshalJSON accepts only a JSON string, never a number, so no float ever
// touches a money value.
func (m *Money) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil || len(data) == 0 || data[0] != '"' {
		return fmt.Errorf("money must be a JSON string like \"1.250000\", got %s", data)
	}
	v, err := ParseMoney(s)
	if err != nil {
		return err
	}
	*m = v
	return nil
}
