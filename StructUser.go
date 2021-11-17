package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	Age  string `json:"age"`
}

// getDB godoc
// @Summary getDB
// @Description Get All Students
// @Accept json
// @Produce json
// @Success 200 {object} Student
// @Router /db [get]
