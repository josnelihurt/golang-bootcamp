# toonwire

Una librería de Go para serializar estructuras a **TOON (Token-Oriented Object Notation)** y JSON con un enfoque especial en **prompts para LLMs** y cargas útiles optimizadas en tokens.

## Motivación

- Reducir el costo en tokens frente a JSON tradicional en intercambios con modelos de lenguaje.
- Ofrecer un API idiomática en Go, con soporte para `Marshal`/`Unmarshal` simétricos.
- Permitir a los equipos integrar fácilmente TOON en pipelines existentes de prompts y datos estructurados.

## Características

- `Transformer` configurable con opciones globales de TOON (`LengthMarkers`, `Indent`).
- Interfaces `TOONMarshaler` y `TOONUnmarshaler` para controlar la serialización personalizada.
- Helpers genéricos para convertir `JSON ↔ TOON` mediante un tipo intermedio fuerte.
- Helper específico `MarshalForLLM` que activa ajustes sugeridos para prompts repetitivos.
- CLI de referencia para convertir entre formatos desde la terminal.

## Instalación

```bash
go get github.com/llm-tools/toonwire
```

## Uso básico

```go
package main

import (
    "fmt"

    "github.com/llm-tools/toonwire/pkg/wire"
)

type User struct {
    ID    int      `json:"id" toon:"id"`
    Name  string   `json:"name" toon:"name"`
    Email string   `json:"email" toon:"email"`
    Roles []string `json:"roles" toon:"roles"`
}

func main() {
    transformer := wire.NewTransformer()
    transformer.Toon.LengthMarkers = true // sugerido para LLMs

    user := User{ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}}

    toonPayload, err := transformer.MarshalToTOON(user)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(toonPayload))
}
```

Salida aproximada:

```json
{"id":1,"name":"Alice","email":"alice@example.com","roles":["#len=2","admin","editor"]}
```

## Conversión JSON ↔ TOON

```go
transformer := wire.NewTransformer()
transformer.Toon.LengthMarkers = true

jsonData := []byte(`{"id":1,"name":"Alice","email":"alice@example.com","roles":["admin","editor"]}`)
toonData, err := wire.JSONToTOON[User](transformer, jsonData)
if err != nil {
    // manejar error
}

backToJSON, err := wire.TOONToJSON[User](transformer, toonData)
if err != nil {
    // manejar error
}
```

## Prompts para LLMs

```go
org := wire.Organization{
    Name: "Acme",
    Users: []wire.User{
        {ID: 1, Name: "Alice", Email: "alice@example.com", Roles: []string{"admin", "editor"}},
        {ID: 2, Name: "Bob", Email: "bob@example.com", Roles: []string{"viewer"}},
    },
}

prompt, err := transformer.MarshalForLLM(org)
if err != nil {
    // manejar error
}
fmt.Println(prompt)
```

`MarshalForLLM` activa marcadores de longitud e indentación legible para repetir grupos homogéneos (tablas, listas, etc.).

## CLI de ejemplo

```bash
# De JSON a TOON (stdin → stdout)
go run ./cmd/toonwire -from json -to toon < input.json > output.toon

# De TOON a JSON prettificado
go run ./cmd/toonwire -from toon -pretty -output out.json -input input.toon
```

## Limitaciones y advertencias

- TOON se basa en JSON; los marcadores de longitud son meta-información que los modelos interpretan según el contexto.
- La reducción real de tokens depende del tokenizer específico de cada modelo.
- Para cargas `streaming`, puede ser preferible JSON tradicional por compatibilidad con APIs existentes.

## Estado del proyecto

- Estado: **Experimental**.
- Versionado: **Semantic Versioning** (cuando se alcance v1.0.0). Antes de esa versión pueden darse cambios incompatibles.

## Contribuciones

Consulta [CONTRIBUTING.md](CONTRIBUTING.md) para las pautas de desarrollo y estilo.

## Licencia

Distribuido bajo la licencia MIT. Consulta [LICENSE](LICENSE).
