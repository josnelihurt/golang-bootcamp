package wire

import (
	"fmt"
)

func ExampleTransformer_MarshalToTOON() {
	transformer := NewTransformer()
	user := User{ID: 1, Name: "Ada", Email: "ada@example.com", Role: "admin"}
	data, _ := transformer.MarshalToTOON(user)
	fmt.Println(string(data))
	// Output:
	// {"email":"ada@example.com","id":1,"name":"Ada","role":"admin"}
}

func ExampleJSONToTOON() {
	transformer := NewTransformer()
	data, _ := JSONToTOON[User](transformer, []byte(`{"id":99,"name":"LLM"}`))
	fmt.Println(string(data))
	// Output:
	// {"email":"","id":99,"name":"LLM","role":""}
}

func ExampleTransformer_MarshalForLLM() {
	transformer := NewTransformer()
	text, _ := transformer.MarshalForLLM(Organization{
		Name:        "PromptOps",
		Description: "Pipeline de prompts",
		Users: []User{
			{ID: 7, Name: "Bot", Email: "bot@example.com", Role: "automation"},
		},
	})
	fmt.Println(text)
	// Output:
	// {"description":"Pipeline de prompts","name":"PromptOps","users":[#1][{"email":"bot@example.com","id":7,"name":"Bot","role":"automation"}]}
}
