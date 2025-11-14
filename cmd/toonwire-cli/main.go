package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cursor-ai/toonwire/pkg/wire"
)

func main() {
	var (
		from   = flag.String("from", "json", "formato de entrada: json o toon")
		to     = flag.String("to", "toon", "formato de salida: json o toon")
		llm    = flag.Bool("llm", false, "aplicar ajustes recomendados para prompts LLM")
		sample = flag.Bool("sample", false, "usar el ejemplo Organization embebido en lugar de STDIN")
	)
	flag.Parse()

	transformer := wire.NewTransformer()
	if *llm {
		transformer.Toon.LengthMarkers = true
	}

	input, err := readInput(*from, *sample, transformer)
	if err != nil {
		log.Fatalf("error leyendo entrada: %v", err)
	}

	var output []byte

	switch {
	case *from == "json" && *to == "toon":
		if *llm {
			var doc map[string]any
			if err := transformer.UnmarshalFromJSON(input, &doc); err != nil {
				log.Fatalf("json inválido: %v", err)
			}
			text, err := transformer.MarshalForLLM(doc)
			if err != nil {
				log.Fatalf("marshal LLM: %v", err)
			}
			fmt.Println(text)
			return
		}
		output, err = wire.JSONToTOON[map[string]any](transformer, input)
	case *from == "toon" && *to == "json":
		output, err = wire.TOONToJSON[map[string]any](transformer, input)
	case *from == *to:
		output = input
	default:
		log.Fatalf("conversión no soportada: %s -> %s", *from, *to)
	}
	if err != nil {
		log.Fatalf("error durante la conversión: %v", err)
	}
	os.Stdout.Write(output)
	if len(output) == 0 || output[len(output)-1] != '\n' {
		fmt.Println()
	}
}

func readInput(from string, useSample bool, t *wire.Transformer) ([]byte, error) {
	if !useSample {
		return io.ReadAll(os.Stdin)
	}
	org := wire.Organization{
		Name:        "Sample Org",
		Description: "Librería TOON/JSON demo",
		Users: []wire.User{
			{ID: 1, Name: "Ada Lovelace", Email: "ada@example.com", Role: "admin"},
			{ID: 2, Name: "Alan Turing", Email: "alan@example.com", Role: "researcher"},
		},
	}
	format := strings.ToLower(from)
	switch format {
	case "json":
		return t.Marshal(org, wire.FormatJSON)
	case "toon":
		return t.Marshal(org, wire.FormatTOON)
	default:
		return nil, fmt.Errorf("formato de muestra desconocido: %s", from)
	}
}
