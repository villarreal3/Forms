package helpers

import (
	"fmt"
	"time"
)

// DateLayouts son los formatos de fecha soportados
var DateLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
}

// ParseExpiry parsea una fecha de expiración y valida que sea futura
// respecto a la hora actual del servidor (UTC).
func ParseExpiry(input string) (time.Time, error) {
	parsed, err := parseDate(input)
	if err != nil {
		return time.Time{}, err
	}

	if !parsed.After(time.Now()) {
		return time.Time{}, fmt.Errorf("La fecha de expiración debe ser una fecha futura")
	}

	return parsed, nil
}

// ParseDateTime parsea una fecha sin validar que sea futura.
func ParseDateTime(input string) (time.Time, error) {
	return parseDate(input)
}

func parseDate(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, fmt.Errorf("fecha vacía")
	}

	loc := time.Now().Location()

	for _, layout := range DateLayouts {
		if t, err := time.ParseInLocation(layout, input, loc); err == nil {
			return t, nil
		}
	}

	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t.In(loc), nil
	}

	return time.Time{}, fmt.Errorf("formato inválido (%s)", input)
}
