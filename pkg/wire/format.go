package wire

// Format representa un formato serializable soportado por la librería.
type Format int

// Lista de formatos soportados por Transformer.
const (
	FormatJSON Format = iota
	FormatTOON
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatTOON:
		return "toon"
	default:
		return "unknown"
	}
}
