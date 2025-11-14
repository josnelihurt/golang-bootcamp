package wire

// User representa un usuario genérico para ejemplos y pruebas.
type User struct {
	ID    int    `json:"id" toon:"id"`
	Name  string `json:"name" toon:"name"`
	Email string `json:"email" toon:"email"`
	Role  string `json:"role" toon:"role"`
}

// Organization modela un conjunto de usuarios bajo un mismo equipo.
type Organization struct {
	Name        string `json:"name" toon:"name"`
	Description string `json:"description" toon:"description"`
	Users       []User `json:"users" toon:"users"`
}
