package vars

import "database/sql"

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	Age  string `json:"age"`
}

var DB []Student
var Database *sql.DB

type Student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	Age  string `json:"age"`
}
