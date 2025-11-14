package wire

// ToonOptions controla ajustes específicos del formato TOON.
type ToonOptions struct {
	// LengthMarkers inserta anotaciones [#N] antes de cada arreglo para
	// ayudar a los LLMs a razonar sobre el tamaño de las colecciones.
	LengthMarkers bool
	// Indent controla la indentación cuando se habilita pretty-print.
	Indent string
	// SortKeys fuerza el orden alfabético en los objetos para mantener la
	// estabilidad de los diffs.
	SortKeys bool
	// EscapeHTML replica el comportamiento por defecto de encoding/json.
	EscapeHTML bool
}

// DefaultToonOptions devuelve los valores recomendados para la mayoría de escenarios.
func DefaultToonOptions() ToonOptions {
	return ToonOptions{
		LengthMarkers: false,
		Indent:        "",
		SortKeys:      true,
		EscapeHTML:    false,
	}
}
