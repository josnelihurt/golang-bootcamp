package wire

// User representa un usuario ejemplo usado en documentación y tests.
type User struct {
	ID       int            `json:"id" toon:"id"`
	Name     string         `json:"name" toon:"name"`
	Email    string         `json:"email" toon:"email"`
	Roles    []string       `json:"roles" toon:"roles"`
	Metadata map[string]any `json:"metadata,omitempty" toon:"metadata,omitempty"`
}

// Organization es un ejemplo de entidad con estructuras anidadas.
type Organization struct {
	Name  string `json:"name" toon:"name"`
	Users []User `json:"users" toon:"users"`
}
