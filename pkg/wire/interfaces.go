package wire

// TOONMarshaler permite que un tipo defina su propia serialización TOON.
type TOONMarshaler interface {
	MarshalTOON() ([]byte, error)
}

// TOONUnmarshaler permite que un tipo defina cómo deserializar datos TOON.
type TOONUnmarshaler interface {
	UnmarshalTOON([]byte) error
}
