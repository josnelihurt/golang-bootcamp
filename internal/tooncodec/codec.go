package tooncodec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Options define cómo renderizar estructuras en formato TOON.
type Options struct {
	// LengthMarkers inserta un marcador "#len=N" como primer elemento en los arreglos.
	LengthMarkers bool
	// Indent aplica indentación cuando se serializa. Si está vacío, el resultado es compacto.
	Indent string
}

// Codec implementa una codificación TOON simplificada basada en JSON.
type Codec struct{}

// New crea una instancia lista para usar.
func New() *Codec {
	return &Codec{}
}

// Marshal serializa un valor a una representación TOON basada en JSON.
// La implementación usa JSON como representación subyacente para simplificar
// la interoperabilidad con bibliotecas reales de TOON.
func (c *Codec) Marshal(v any, opts Options) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("tooncodec: no se puede serializar un valor nil")
	}

	// Convertimos el valor a una estructura intermedia.
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("tooncodec: json marshal: %w", err)
	}

	var intermediate any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&intermediate); err != nil {
		return nil, fmt.Errorf("tooncodec: json decode intermedio: %w", err)
	}

	if opts.LengthMarkers {
		intermediate = applyLengthMarkers(intermediate)
	}

	if opts.Indent != "" {
		return json.MarshalIndent(intermediate, "", opts.Indent)
	}

	return json.Marshal(intermediate)
}

// Unmarshal decodifica datos TOON y los proyecta en el valor destino.
func (c *Codec) Unmarshal(data []byte, opts Options, v any) error {
	if v == nil {
		return fmt.Errorf("tooncodec: el destino no puede ser nil")
	}

	var intermediate any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&intermediate); err != nil {
		return fmt.Errorf("tooncodec: json decode desde TOON: %w", err)
	}

	intermediate = stripLengthMarkers(intermediate)

	cleaned, err := json.Marshal(intermediate)
	if err != nil {
		return fmt.Errorf("tooncodec: reconstrucción json: %w", err)
	}

	if err := json.Unmarshal(cleaned, v); err != nil {
		return fmt.Errorf("tooncodec: json unmarshal final: %w", err)
	}

	return nil
}

func applyLengthMarkers(v any) any {
	switch val := v.(type) {
	case []any:
		result := make([]any, 0, len(val)+1)
		result = append(result, fmt.Sprintf("#len=%d", len(val)))
		for _, item := range val {
			result = append(result, applyLengthMarkers(item))
		}
		return result
	case map[string]any:
		for k, item := range val {
			val[k] = applyLengthMarkers(item)
		}
		return val
	default:
		return v
	}
}

func stripLengthMarkers(v any) any {
	switch val := v.(type) {
	case []any:
		var start int
		if len(val) > 0 {
			if marker, ok := val[0].(string); ok && strings.HasPrefix(marker, "#len=") {
				start = 1
			}
		}
		for i := start; i < len(val); i++ {
			val[i] = stripLengthMarkers(val[i])
		}
		return val[start:]
	case map[string]any:
		for k, item := range val {
			val[k] = stripLengthMarkers(item)
		}
		return val
	default:
		return v
	}
}
