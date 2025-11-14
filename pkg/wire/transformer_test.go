package wire

import (
	"errors"
	"strings"
	"testing"
)

func TestTransformerRoundTripTOON(t *testing.T) {
	transformer := NewTransformer(WithToonOptions(ToonOptions{
		LengthMarkers: true,
		Indent:        "",
		SortKeys:      true,
		EscapeHTML:    false,
	}))

	original := Organization{
		Name:        "Playground",
		Description: "Org para pruebas",
		Users: []User{
			{ID: 1, Name: "Tess", Email: "tess@example.com", Role: "admin"},
			{ID: 2, Name: "Jules", Email: "jules@example.com", Role: "analyst"},
		},
	}

	data, err := transformer.MarshalToTOON(original)
	if err != nil {
		t.Fatalf("marshal toon: %v", err)
	}
	if !strings.Contains(string(data), "[#2]") {
		t.Fatalf("esperaba marcador de longitud, obtuvo: %s", string(data))
	}

	var copy Organization
	if err := transformer.UnmarshalFromTOON(data, &copy); err != nil {
		t.Fatalf("unmarshal toon: %v", err)
	}
	if copy.Name != original.Name || len(copy.Users) != len(original.Users) {
		t.Fatalf("roundtrip falló: %+v", copy)
	}
}

func TestJSONToTOONHelper(t *testing.T) {
	transformer := NewTransformer()
	jsonPayload := []byte(`{"id":10,"name":"Helper"}`)

	data, err := JSONToTOON[User](transformer, jsonPayload)
	if err != nil {
		t.Fatalf("json->toon: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("json->toon retornó vacío")
	}
}

func TestTOONToJSONHelper(t *testing.T) {
	transformer := NewTransformer()
	original := User{ID: 42, Name: "TOON", Email: "toon@example.com", Role: "bot"}

	toonData, err := transformer.MarshalToTOON(original)
	if err != nil {
		t.Fatalf("marshal toon: %v", err)
	}

	jsonData, err := TOONToJSON[User](transformer, toonData)
	if err != nil {
		t.Fatalf("toon->json: %v", err)
	}
	if !strings.Contains(string(jsonData), `"id":42`) {
		t.Fatalf("json inesperado: %s", string(jsonData))
	}
}

type customMarshaler struct{}

func (customMarshaler) MarshalTOON() ([]byte, error) {
	return []byte("custom-toon"), nil
}

type customUnmarshaler struct {
	value string
}

func (c *customUnmarshaler) UnmarshalTOON(data []byte) error {
	if len(data) == 0 {
		return errors.New("sin datos")
	}
	c.value = string(data)
	return nil
}

func TestCustomInterfaces(t *testing.T) {
	transformer := NewTransformer()
	data, err := transformer.MarshalToTOON(customMarshaler{})
	if err != nil {
		t.Fatalf("marshal custom: %v", err)
	}
	if string(data) != "custom-toon" {
		t.Fatalf("marshal custom inesperado: %s", data)
	}

	target := &customUnmarshaler{}
	if err := transformer.UnmarshalFromTOON([]byte("payload"), target); err != nil {
		t.Fatalf("unmarshal custom: %v", err)
	}
	if target.value != "payload" {
		t.Fatalf("valor custom inesperado: %s", target.value)
	}
}

func TestMarshalForLLM(t *testing.T) {
	transformer := NewTransformer()
	org := Organization{
		Name: "LLM",
		Users: []User{
			{ID: 1, Name: "Bot", Email: "bot@example.com", Role: "ai"},
		},
	}
	text, err := transformer.MarshalForLLM(org)
	if err != nil {
		t.Fatalf("marshal LLM: %v", err)
	}
	if !strings.Contains(text, "[#1]") {
		t.Fatalf("MarshalForLLM debería incluir marcadores, obtuvo: %s", text)
	}
}

func TestMarshalErrors(t *testing.T) {
	transformer := NewTransformer()
	if _, err := transformer.Marshal(nil, FormatJSON); !errors.Is(err, ErrNilValue) {
		t.Fatalf("esperaba ErrNilValue, obtuvo: %v", err)
	}
	if err := transformer.Unmarshal(nil, FormatJSON, nil); !errors.Is(err, ErrNilTarget) {
		t.Fatalf("esperaba ErrNilTarget, obtuvo: %v", err)
	}
	if _, err := transformer.Marshal(User{}, Format(99)); !errors.Is(err, ErrUnsupportedFormat) {
		t.Fatalf("esperaba ErrUnsupportedFormat, obtuvo: %v", err)
	}
}
