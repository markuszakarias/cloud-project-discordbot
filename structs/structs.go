package structs

// Struct todo_struct - struct for a todo object
type Todo_struct struct {
	Id       int    `json:"Id"`
	Title    string `json:"Title"`
	Category string `json:"Category"`
	State    string `json:"State"`
}
