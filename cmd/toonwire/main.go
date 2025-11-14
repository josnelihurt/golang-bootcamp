package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/llm-tools/toonwire/pkg/wire"
)

func main() {
	var fromFlag string
	var toFlag string
	var inputPath string
	var outputPath string
	var pretty bool

	flag.StringVar(&fromFlag, "from", "json", "Formato de entrada: json | toon")
	flag.StringVar(&toFlag, "to", "", "Formato de salida: json | toon (por defecto el complemento)")
	flag.StringVar(&inputPath, "input", "", "Ruta del archivo de entrada (STDIN si se omite)")
	flag.StringVar(&outputPath, "output", "", "Ruta del archivo de salida (STDOUT si se omite)")
	flag.BoolVar(&pretty, "pretty", false, "Emite JSON con indentación legible")
	flag.Parse()

	fromFormat, err := parseFormat(fromFlag)
	if err != nil {
		exitErr(err)
	}

	if toFlag == "" {
		if fromFormat == wire.FormatJSON {
			toFlag = "toon"
		} else {
			toFlag = "json"
		}
	}

	toFormat, err := parseFormat(toFlag)
	if err != nil {
		exitErr(err)
	}

	inputData, err := readAll(inputPath)
	if err != nil {
		exitErr(err)
	}

	transformer := wire.NewTransformer()

	var output []byte

	switch {
	case fromFormat == wire.FormatJSON && toFormat == wire.FormatTOON:
		output, err = jsonToToon(transformer, inputData)
	case fromFormat == wire.FormatTOON && toFormat == wire.FormatJSON:
		output, err = toonToJSON(transformer, inputData, pretty)
	default:
		// Si ambos formatos son iguales, devolvemos los datos tal cual.
		output = inputData
	}

	if err != nil {
		exitErr(err)
	}

	if err := writeAll(outputPath, output); err != nil {
		exitErr(err)
	}
}

func parseFormat(value string) (wire.Format, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "json":
		return wire.FormatJSON, nil
	case "toon":
		return wire.FormatTOON, nil
	default:
		return 0, fmt.Errorf("formato desconocido: %s", value)
	}
}

func readAll(path string) ([]byte, error) {
	if path == "" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writeAll(path string, data []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(data)
		if err == nil && len(data) > 0 && data[len(data)-1] != '\n' {
			_, err = os.Stdout.Write([]byte("\n"))
		}
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "toonwire: %v\n", err)
	os.Exit(1)
}

func jsonToToon(t *wire.Transformer, data []byte) ([]byte, error) {
	var intermediate any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&intermediate); err != nil {
		return nil, fmt.Errorf("error leyendo JSON: %w", err)
	}
	return t.MarshalToTOON(intermediate)
}

func toonToJSON(t *wire.Transformer, data []byte, pretty bool) ([]byte, error) {
	var intermediate any
	if err := t.Unmarshal(data, wire.FormatTOON, &intermediate); err != nil {
		return nil, fmt.Errorf("error leyendo TOON: %w", err)
	}
	if pretty {
		return json.MarshalIndent(intermediate, "", "  ")
	}
	return json.Marshal(intermediate)
}
