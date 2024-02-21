package types

type Person struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
	Role string `json:"role"`
	// Where the person is
	Stage string `json:"stage"`
}
