## Contribuir a toonwire

Gracias por considerar contribuir. Este proyecto busca ofrecer una base sólida para manipular TOON/JSON en Go.

### Flujo sugerido

1. Abre un issue describiendo el problema o la propuesta.
2. Crea una rama desde `main` o la rama activa del feature.
3. Añade tests: ningún cambio funcional se acepta sin cobertura.
4. Ejecuta `go test ./...` y, si aplica, `go test -run Example`.
5. Envía el PR enlazando el issue y describe impactos o migraciones.

### Estilo y convenciones

- Usa `gofmt` y `golangci-lint` (opcional) antes de subir cambios.
- Prefiere APIs pequeñas y opciones funcionales sobre configuraciones globales.
- Documenta tipos y funciones exportadas con comentarios `//`.
- Los nuevos formatos deben integrarse al enum `Format` y a `Transformer`.

### Seguridad y calidad

- Evita dependencias con licencias incompatibles.
- No introduzcas logging global; usa errores envueltos (`fmt.Errorf("context: %w", err)`).
- Incluye benchmarks cuando el cambio afecte rendimiento.

Gracias por ayudar a que `toonwire` sea útil para la comunidad Go + LLM.
