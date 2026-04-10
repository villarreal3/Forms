package main

import (
	"fmt"
	"time"

	"forms/handlers/helpers"
)

func main() {
	fmt.Println("=== Test ParseExpiry ===")
	fmt.Printf("Hora del servidor (UTC): %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	tests := []struct {
		label string
		input string
	}{
		{"Fecha futura (debería pasar)", "2026-12-31 23:59:00"},
		{"Fecha pasada (debería fallar)", "2020-01-01 00:00:00"},
		{"Formato con T", "2026-12-31T23:59:00"},
		{"Formato sin segundos", "2026-12-31 23:59"},
		{"Formato inválido", "31/12/2026"},
		{"Vacío", ""},
	}

	for _, tc := range tests {
		fmt.Printf("  [%s] input=%q\n", tc.label, tc.input)
		result, err := helpers.ParseExpiry(tc.input)
		if err != nil {
			fmt.Printf("    ERROR: %v\n\n", err)
		} else {
			fmt.Printf("    OK: %s\n\n", result.Format("2006-01-02 15:04:05 MST"))
		}
	}
}
