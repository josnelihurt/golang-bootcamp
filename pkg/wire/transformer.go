package wire

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/cursor-ai/toonwire/internal/tooncodec"
)

// ErrUnsupportedFormat indica que se solicitó un formato desconocido.
var ErrUnsupportedFormat = errors.New("wire: formato no soportado")

// ErrNilValue indica que se intentó serializar un valor nulo.
var ErrNilValue = errors.New("wire: valor nulo")

// ErrNilTarget indica que la operación de Unmarshal recibió un destino nulo.
var ErrNilTarget = errors.New("wire: destino nulo")

// Transformer centraliza las opciones de serialización.
type Transformer struct {
	Toon           ToonOptions
	jsonIndent     string
	jsonEscapeHTML bool
}

// Option representa una opción funcional para NewTransformer.
type Option func(*Transformer)

// NewTransformer crea una nueva instancia con opciones seguras por defecto.
func NewTransformer(opts ...Option) *Transformer {
	t := &Transformer{
		Toon:           DefaultToonOptions(),
		jsonIndent:     "",
		jsonEscapeHTML: true,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// WithToonOptions reemplaza la configuración por defecto de TOON.
func WithToonOptions(opts ToonOptions) Option {
	return func(t *Transformer) {
		t.Toon = opts
	}
}

// WithJSONIndent establece la indentación para Marshal JSON.
func WithJSONIndent(indent string) Option {
	return func(t *Transformer) {
		t.jsonIndent = indent
	}
}

// WithJSONEscapeHTML controla si se escapan caracteres HTML en JSON.
func WithJSONEscapeHTML(enable bool) Option {
	return func(t *Transformer) {
		t.jsonEscapeHTML = enable
	}
}

// Marshal serializa un valor en el formato indicado.
func (t *Transformer) Marshal(v any, f Format) ([]byte, error) {
	switch f {
	case FormatJSON:
		return t.MarshalToJSON(v)
	case FormatTOON:
		return t.MarshalToTOON(v)
	default:
		return nil, ErrUnsupportedFormat
	}
}

// Unmarshal deserializa datos en el formato indicado dentro del destino.
func (t *Transformer) Unmarshal(data []byte, f Format, v any) error {
	switch f {
	case FormatJSON:
		return t.UnmarshalFromJSON(data, v)
	case FormatTOON:
		return t.UnmarshalFromTOON(data, v)
	default:
		return ErrUnsupportedFormat
	}
}

// MarshalToTOON serializa un valor a TOON respetando las interfaces personalizadas.
func (t *Transformer) MarshalToTOON(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrNilValue
	}
	if marshaler, ok := asTOONMarshaler(v); ok {
		return marshaler.MarshalTOON()
	}
	return tooncodec.Marshal(v, t.toInternalOptions())
}

// MarshalToJSON serializa un valor a JSON siguiendo las opciones del Transformer.
func (t *Transformer) MarshalToJSON(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrNilValue
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(t.jsonEscapeHTML)
	if t.jsonIndent != "" {
		enc.SetIndent("", t.jsonIndent)
	}
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// UnmarshalFromTOON deserializa datos TOON respetando interfaces personalizadas.
func (t *Transformer) UnmarshalFromTOON(data []byte, v any) error {
	if v == nil {
		return ErrNilTarget
	}
	if unmarshaler, ok := asTOONUnmarshaler(v); ok {
		if err := unmarshaler.UnmarshalTOON(data); err != nil {
			return fmt.Errorf("custom toon unmarshal: %w", err)
		}
		return nil
	}
	if err := tooncodec.Unmarshal(data, v, t.toInternalOptions()); err != nil {
		return fmt.Errorf("unmarshal toon: %w", err)
	}
	return nil
}

// UnmarshalFromJSON deserializa datos JSON usando encoding/json.
func (t *Transformer) UnmarshalFromJSON(data []byte, v any) error {
	if v == nil {
		return ErrNilTarget
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}

// MarshalForLLM genera una representación TOON recomendada para prompts LLM.
func (t *Transformer) MarshalForLLM(v any) (string, error) {
	opts := t.Toon
	opts.LengthMarkers = true
	if opts.Indent == "" {
		opts.Indent = ""
	}
	if v == nil {
		return "", ErrNilValue
	}
	if marshaler, ok := asTOONMarshaler(v); ok {
		data, err := marshaler.MarshalTOON()
		return string(data), err
	}
	data, err := tooncodec.Marshal(v, toInternalToonOptions(opts))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *Transformer) toInternalOptions() tooncodec.Options {
	return toInternalToonOptions(t.Toon)
}

func toInternalToonOptions(opts ToonOptions) tooncodec.Options {
	return tooncodec.Options{
		LengthMarkers: opts.LengthMarkers,
		Indent:        opts.Indent,
		SortKeys:      opts.SortKeys,
		EscapeHTML:    opts.EscapeHTML,
	}
}

func asTOONMarshaler(v any) (TOONMarshaler, bool) {
	if v == nil {
		return nil, false
	}
	if marshaler, ok := v.(TOONMarshaler); ok {
		return marshaler, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer && !rv.IsNil() {
		if marshaler, ok := rv.Interface().(TOONMarshaler); ok {
			return marshaler, true
		}
	}
	return nil, false
}

func asTOONUnmarshaler(v any) (TOONUnmarshaler, bool) {
	if v == nil {
		return nil, false
	}
	if unmarshaler, ok := v.(TOONUnmarshaler); ok {
		return unmarshaler, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer && !rv.IsNil() {
		if unmarshaler, ok := rv.Interface().(TOONUnmarshaler); ok {
			return unmarshaler, true
		}
	}
	return nil, false
}
