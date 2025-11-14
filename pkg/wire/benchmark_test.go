package wire

import (
	"testing"
)

func BenchmarkMarshalJSON(b *testing.B) {
	transformer := NewTransformer()
	org := benchmarkOrg()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := transformer.MarshalToJSON(org); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalTOON(b *testing.B) {
	transformer := NewTransformer(WithToonOptions(ToonOptions{
		LengthMarkers: true,
		SortKeys:      true,
	}))
	org := benchmarkOrg()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := transformer.MarshalToTOON(org); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkOrg() Organization {
	users := make([]User, 0, 10)
	for i := 0; i < 10; i++ {
		users = append(users, User{
			ID:    i + 1,
			Name:  "Researcher",
			Email: "researcher@example.com",
			Role:  "analyst",
		})
	}
	return Organization{
		Name:        "BenchOrg",
		Description: "Org utilizada en benchmarks; el conteo real de tokens depende del tokenizer del LLM.",
		Users:       users,
	}
}
