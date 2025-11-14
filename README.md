# toonwire

Librería Go profesional para transformar estructuras entre JSON y TOON (Token-Oriented Object Notation) con énfasis en prompts para LLMs.

## Motivación

- JSON es ubicuo, pero su verbosidad incrementa el conteo de tokens en modelos de lenguaje.
- TOON prioriza estabilidad y brevedad (orden determinista, sin espacios innecesarios, marcadores opcionales de longitud).
- `Transformer` unifica opciones y permite cambiar de backend TOON sin romper la API.

## Instalación

```bash
go get github.com/cursor-ai/toonwire
```

## Uso rápido

```go
package main

import (
	"fmt"

	"github.com/cursor-ai/toonwire/pkg/wire"
)

type User struct {
	ID   int    `json:"id" toon:"id"`
	Name string `json:"name" toon:"name"`
	Role string `json:"role" toon:"role"`
}

func main() {
	transformer := wire.NewTransformer()
	user := User{ID: 7, Name: "Ada", Role: "research"}

	toonBytes, _ := transformer.MarshalToTOON(user)
	fmt.Println(string(toonBytes))

	jsonBytes, _ := transformer.MarshalToJSON(user)
	fmt.Println(string(jsonBytes))
}
```

### Helpers JSON ↔ TOON

```go
t := wire.NewTransformer()
toonData, _ := wire.JSONToTOON[wire.Organization](t, []byte(`{"name":"PromptOps"}`))
jsonData, _ := wire.TOONToJSON[wire.Organization](t, toonData)
```

### Payloads pensados para LLMs

```go
prompt, _ := t.MarshalForLLM(org)
// => {"name":"PromptOps","users":[#2][{"id":1,"name":"Ada"}, ...]}
```

## Ejemplos compilables

- `examples/basic`: programa completo con JSON y TOON.
- `cmd/toonwire-cli`: CLI para convertir flujos (`json↔toon`, modo LLM, muestra embebida).

Ejecuta `go test ./pkg/wire -run Example` para ver los ejemplos de la documentación.

## Limitaciones y advertencias

- TOON sigue siendo textual; la reducción de tokens real depende del tokenizer del modelo.
- JSON continúa siendo mejor cuando necesitas compatibilidad universal o validadores estrictos.
- Los marcadores `[#{n}]` son opcionales. Si otro parser los desconoce, desactívalos.
- El códec interno es de referencia; reemplázalo con `github.com/toon-format/toon-go` o equivalente cuando esté disponible.

## Estado y versionado

- Estado: **Alpha**. La API puede recibir ajustes menores, pero `Transformer` y los helpers se consideran estables.
- Versionado: **Semantic Versioning**. Parches sólo corrigen bugs; versiones menores agregan funcionalidad retrocompatible.

## Roadmap corto

- Agregar soporte opcional para YAML/CBOR mediante el mismo `Format`.
- Profundizar en métricas de conteo de tokens por modelo.
- Añadir perfiles de serialización para streams y partial updates.

## Licencia

MIT. Ver `LICENSE`.
