package wire

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type customTOON struct {
	Value string `json:"value"`
}

func (c customTOON) MarshalTOON() ([]byte, error) {
	return []byte("{\"value\":\"marshaled\"}"), nil
}

func (c *customTOON) UnmarshalTOON(data []byte) error {
	c.Value = "from-toon"
	return nil
}

func TestTransformerMarshalToTOON(t *testing.T) {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	input := User{
		ID:    42,
		Name:  "Carol",
		Email: "carol@example.com",
		Roles: []string{"ops"},
	}

	got, err := transformer.MarshalToTOON(input)
	if err != nil {
		t.Fatalf("MarshalToTOON devolvió error: %v", err)
	}

	if !strings.Contains(string(got), "\"#len=1\"") {
		t.Fatalf("no se insertó marcador de longitud: %s", got)
	}

	var decoded User
	if err := transformer.Unmarshal(got, FormatTOON, &decoded); err != nil {
		t.Fatalf("roundtrip falló: %v", err)
	}

	if decoded.ID != input.ID || decoded.Email != input.Email || len(decoded.Roles) != 1 || decoded.Roles[0] != "ops" {
		t.Fatalf("resultado inesperado tras roundtrip: %+v", decoded)
	}
}

func TestTransformerUnmarshalFromTOON(t *testing.T) {
	transformer := NewTransformer()

	data := []byte(`{"id":42,"name":"Carol","email":"carol@example.com","roles":["#len=1","ops"]}`)

	var dst User
	if err := transformer.Unmarshal(data, FormatTOON, &dst); err != nil {
		t.Fatalf("Unmarshal devolvió error: %v", err)
	}

	if dst.ID != 42 || dst.Email != "carol@example.com" || len(dst.Roles) != 1 {
		t.Fatalf("valor inesperado: %+v", dst)
	}
}

func TestTransformerMarshalNested(t *testing.T) {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	org := Organization{
		Name:  "Acme",
		Users: []User{{ID: 1, Name: "Alice", Email: "alice@example.com"}, {ID: 2, Name: "Bob", Email: "bob@example.com"}},
	}

	out, err := transformer.MarshalToTOON(org)
	if err != nil {
		t.Fatalf("MarshalToTOON devolvió error: %v", err)
	}

	if !strings.Contains(string(out), "#len=2") {
		t.Fatalf("se esperaba marcador de longitud en la salida: %s", out)
	}
}

func TestJSONToTOON(t *testing.T) {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	payload := []byte(`{"id":5,"name":"Eve","email":"eve@example.com","roles":["analyst","reviewer"]}`)

	out, err := JSONToTOON[User](transformer, payload)
	if err != nil {
		t.Fatalf("JSONToTOON devolvió error: %v", err)
	}

	if !strings.Contains(string(out), "#len=2") {
		t.Fatalf("se esperaba marcador de longitud\noutput: %s", out)
	}
}

func TestTOONToJSON(t *testing.T) {
	transformer := NewTransformer()

	data := []byte(`{"id":7,"name":"Frank","email":"frank@example.com","roles":["#len=0"]}`)
	out, err := TOONToJSON[User](transformer, data)
	if err != nil {
		t.Fatalf("TOONToJSON devolvió error: %v", err)
	}

	if !strings.Contains(string(out), `"roles":[]`) {
		t.Fatalf("se esperaba slice vacío en JSON: %s", out)
	}
}

func TestCustomMarshaler(t *testing.T) {
	transformer := NewTransformer()

	got, err := transformer.MarshalToTOON(customTOON{})
	if err != nil {
		t.Fatalf("MarshalToTOON devolvió error: %v", err)
	}
	if string(got) != `{"value":"marshaled"}` {
		t.Fatalf("se esperaba marshal personalizado: %s", got)
	}

	var dst customTOON
	if err := transformer.Unmarshal([]byte(`{"value":"ignored"}`), FormatTOON, &dst); err != nil {
		t.Fatalf("Unmarshal devolvió error: %v", err)
	}
	if dst.Value != "from-toon" {
		t.Fatalf("unmarshal personalizado no ejecutado: %s", dst.Value)
	}
}

func TestUnmarshalNilDestination(t *testing.T) {
	transformer := NewTransformer()

	if err := transformer.Unmarshal([]byte("{}"), FormatTOON, nil); err == nil {
		t.Fatalf("se esperaba error al pasar destino nil")
	}
}

func BenchmarkMarshalTOON(b *testing.B) {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	sample := Organization{
		Name: "Acme",
		Users: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Roles: []string{"viewer"}},
			{ID: 3, Name: "Carol", Email: "carol@example.com", Roles: []string{"ops", "security"}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := transformer.MarshalToTOON(sample); err != nil {
			b.Fatalf("MarshalToTOON error: %v", err)
		}
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	sample := Organization{
		Name: "Acme",
		Users: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Roles: []string{"viewer"}},
			{ID: 3, Name: "Carol", Email: "carol@example.com", Roles: []string{"ops", "security"}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(sample); err != nil {
			b.Fatalf("json.Marshal error: %v", err)
		}
	}
}

func ExampleTransformer_MarshalToTOON() {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Roles: []string{"admin", "editor"},
	}

	data, _ := transformer.MarshalToTOON(user)
	fmt.Println(string(data))
	// Output:
	// {"email":"alice@example.com","id":1,"name":"Alice","roles":["#len=2","admin","editor"]}
}

func ExampleJSONToTOON() {
	transformer := NewTransformer()
	transformer.Toon.LengthMarkers = true

	payload := []byte(`{"id":1,"name":"Alice","email":"alice@example.com","roles":["admin","editor"]}`)

	toonData, _ := JSONToTOON[User](transformer, payload)
	fmt.Println(string(toonData))
	// Output:
	// {"email":"alice@example.com","id":1,"name":"Alice","roles":["#len=2","admin","editor"]}
}

func ExampleTransformer_MarshalForLLM() {
	transformer := NewTransformer()

	org := Organization{
		Name: "Acme",
		Users: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Roles: []string{"viewer"}},
		},
	}

	data, _ := transformer.MarshalForLLM(org)
	fmt.Println(data)
	// Output:
	// {
	//   "name": "Acme",
	//   "users": [
	//     "#len=2",
	//     {
	//       "email": "alice@example.com",
	//       "id": 1,
	//       "name": "Alice",
	//       "roles": [
	//         "#len=2",
	//         "admin",
	//         "editor"
	//       ]
	//     },
	//     {
	//       "email": "bob@example.com",
	//       "id": 2,
	//       "name": "Bob",
	//       "roles": [
	//         "#len=1",
	//         "viewer"
	//       ]
	//     }
	//   ]
	// }
}
