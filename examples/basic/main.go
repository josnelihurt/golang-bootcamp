package main

import (
	"encoding/json"
	"fmt"

	"github.com/llm-tools/toonwire/pkg/wire"
)

func main() {
	transformer := wire.NewTransformer()
	transformer.Toon.LengthMarkers = true

	org := wire.Organization{
		Name: "Acme",
		Users: []wire.User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Roles: []string{"viewer"}},
		},
	}

	toonPayload, err := transformer.MarshalToTOON(org)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TOON:\n%s\n\n", toonPayload)

	jsonPayload, err := transformer.MarshalToJSON(org)
	if err != nil {
		panic(err)
	}

	// Calcular tamaños aproximados.
	fmt.Printf("JSON bytes: %d\n", len(jsonPayload))
	fmt.Printf("TOON bytes: %d\n", len(toonPayload))

	// Convertir de nuevo a estructura para verificación.
	var round wire.Organization
	if err := transformer.Unmarshal(toonPayload, wire.FormatTOON, &round); err != nil {
		panic(err)
	}

	pretty, _ := json.MarshalIndent(round, "", "  ")
	fmt.Printf("Roundtrip:\n%s\n", pretty)
}
