package wire

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/llm-tools/toonwire/internal/tooncodec"
)

// Format representa un formato de serialización soportado.
type Format int

const (
	// FormatJSON representa la serialización estándar JSON.
	FormatJSON Format = iota
	// FormatTOON representa la serialización Token-Oriented Object Notation.
	FormatTOON
)

// ToonOptions describe los parámetros de renderizado para TOON.
type ToonOptions struct {
	// LengthMarkers añade marcadores de longitud al inicio de cada slice.
	LengthMarkers bool
	// Indent define el nivel de indentación para la salida. Cadena vacía => minificado.
	Indent string
}

// DefaultToonOptions devuelve la configuración base recomendada.
func DefaultToonOptions() ToonOptions {
	return ToonOptions{
		LengthMarkers: false,
		Indent:        "",
	}
}

// Transformer centraliza la conversión entre JSON y TOON.
type Transformer struct {
	// Toon especifica las opciones globales usadas al trabajar con TOON.
	Toon ToonOptions

	codec *tooncodec.Codec
}

// NewTransformer crea un Transformer con opciones por defecto.
func NewTransformer() *Transformer {
	return &Transformer{
		Toon:  DefaultToonOptions(),
		codec: tooncodec.New(),
	}
}

// TOONMarshaler permite personalizar la serialización TOON de un tipo.
type TOONMarshaler interface {
	MarshalTOON() ([]byte, error)
}

// TOONUnmarshaler permite personalizar la deserialización TOON de un tipo.
type TOONUnmarshaler interface {
	UnmarshalTOON([]byte) error
}

// Marshal serializa un valor en el formato indicado.
func (t *Transformer) Marshal(v any, f Format) ([]byte, error) {
	if t == nil {
		return nil, fmt.Errorf("wire: transformer nil")
	}
	switch f {
	case FormatJSON:
		return json.Marshal(v)
	case FormatTOON:
		if marshaler, ok := resolveTOONMarshaler(v); ok {
			return marshaler.MarshalTOON()
		}
		return t.codec.Marshal(v, toCodecOptions(t.Toon))
	default:
		return nil, fmt.Errorf("wire: formato %d no soportado", f)
	}
}

// Unmarshal deserializa datos en el formato indicado hacia el destino proporcionado.
func (t *Transformer) Unmarshal(data []byte, f Format, v any) error {
	if t == nil {
		return fmt.Errorf("wire: transformer nil")
	}
	if v == nil {
		return fmt.Errorf("wire: destino nil")
	}

	switch f {
	case FormatJSON:
		return json.Unmarshal(data, v)
	case FormatTOON:
		if unmarshaler, ok := resolveTOONUnmarshaler(v); ok {
			return unmarshaler.UnmarshalTOON(data)
		}
		return t.codec.Unmarshal(data, toCodecOptions(t.Toon), v)
	default:
		return fmt.Errorf("wire: formato %d no soportado", f)
	}
}

// MarshalToTOON es un helper especializado equivalente a Marshal(_, FormatTOON).
func (t *Transformer) MarshalToTOON(v any) ([]byte, error) {
	return t.Marshal(v, FormatTOON)
}

// MarshalToJSON es un helper especializado equivalente a Marshal(_, FormatJSON).
func (t *Transformer) MarshalToJSON(v any) ([]byte, error) {
	return t.Marshal(v, FormatJSON)
}

// MarshalForLLM serializa en TOON usando opciones recomendadas para prompts.
func (t *Transformer) MarshalForLLM(v any) (string, error) {
	if t == nil {
		return "", fmt.Errorf("wire: transformer nil")
	}
	if marshaler, ok := resolveTOONMarshaler(v); ok {
		data, err := marshaler.MarshalTOON()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	opts := t.Toon
	opts.LengthMarkers = true
	if opts.Indent == "" {
		opts.Indent = "  "
	}

	data, err := t.codec.Marshal(v, toCodecOptions(opts))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSONToTOON convierte un payload JSON a TOON a través de un tipo intermedio T.
func JSONToTOON[T any](t *Transformer, jsonData []byte) ([]byte, error) {
	if t == nil {
		t = NewTransformer()
	}
	var typed T
	if err := json.Unmarshal(jsonData, &typed); err != nil {
		return nil, fmt.Errorf("wire: error deserializando JSON intermedio: %w", err)
	}
	return t.MarshalToTOON(typed)
}

// TOONToJSON convierte un payload TOON a JSON usando T como forma intermedia.
func TOONToJSON[T any](t *Transformer, toonData []byte) ([]byte, error) {
	if t == nil {
		t = NewTransformer()
	}
	var typed T
	if err := t.Unmarshal(toonData, FormatTOON, &typed); err != nil {
		return nil, fmt.Errorf("wire: error deserializando TOON intermedio: %w", err)
	}
	return t.MarshalToJSON(typed)
}

func resolveTOONMarshaler(v any) (TOONMarshaler, bool) {
	if v == nil {
		return nil, false
	}
	if m, ok := v.(TOONMarshaler); ok {
		return m, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer && rv.CanAddr() {
		if m, ok := rv.Addr().Interface().(TOONMarshaler); ok {
			return m, true
		}
	}
	return nil, false
}

func resolveTOONUnmarshaler(v any) (TOONUnmarshaler, bool) {
	if v == nil {
		return nil, false
	}
	if um, ok := v.(TOONUnmarshaler); ok {
		return um, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer && !rv.IsNil() {
		if um, ok := rv.Interface().(TOONUnmarshaler); ok {
			return um, true
		}
	}
	return nil, false
}

func toCodecOptions(opts ToonOptions) tooncodec.Options {
	return tooncodec.Options{
		LengthMarkers: opts.LengthMarkers,
		Indent:        opts.Indent,
	}
}
