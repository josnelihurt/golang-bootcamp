# Contribuir a toonwire

¡Gracias por tu interés en colaborar! Este proyecto busca mantenerse simple y fácil de integrar. Para mantener la calidad del código te pedimos seguir estas pautas.

## Requisitos previos

- Go 1.22 o superior.
- `gofmt` ejecutado sobre todos los archivos modificados.
- `go test ./...` debe pasar sin errores ni cambios no comprometidos.

## Flujo de trabajo

1. Crea un fork y una rama descriptiva (`feature/toon-length-markers`).
2. Implementa los cambios siguiendo el estilo idiomático de Go.
3. Añade pruebas unitarias y ejemplos cuando aplique.
4. Ejecuta `go test ./...` antes de abrir el pull request.
5. Describe claramente la motivación del cambio y cualquier limitación conocida.

## Estilo y documentación

- Exporta únicamente los símbolos necesarios para el API público.
- Documenta todo elemento exportado con comentarios `//` válidos para `godoc`.
- Prefiere errores envueltos con contexto (`fmt.Errorf("contexto: %w", err)`).
- Evita dependencias innecesarias; prioriza la compatibilidad con la librería base de TOON.

## Código de conducta

Este proyecto se rige por el [Contributor Covenant](https://www.contributor-covenant.org/). Al participar aceptas actuar con respeto y profesionalismo.

¡Gracias por ayudar a que toonwire mejore!
