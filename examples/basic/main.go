package main

import (
	"fmt"

	"github.com/cursor-ai/toonwire/pkg/wire"
)

func main() {
	transformer := wire.NewTransformer()
	org := wire.Organization{
		Name:        "PromptLab",
		Description: "Equipo de prompt engineering",
		Users: []wire.User{
			{ID: 1, Name: "Luna", Email: "luna@example.com", Role: "lead"},
			{ID: 2, Name: "Kai", Email: "kai@example.com", Role: "research"},
		},
	}

	toon, err := transformer.MarshalToTOON(org)
	if err != nil {
		panic(err)
	}
	fmt.Println("TOON:\n", string(toon))

	jsonPayload, err := transformer.MarshalToJSON(org)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:\n", string(jsonPayload))
}
