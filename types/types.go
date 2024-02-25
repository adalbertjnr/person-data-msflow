package types

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  uint8  `json:"age"`
	Role string `json:"role"`
	// Where the person is
	Stage string `json:"stage"`
}
