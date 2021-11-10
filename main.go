package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/skskuzan/rest_api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Student App API
// @version 1.0
// @description API Server for Students

// host: localhost:5000
// @BasePath /

// securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var db []Student
var database *sql.DB

type Student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	Age  string `json:"age"`
}

func getDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Current Database:\n")
	rows, err := database.Query("select * from people")

	if err != nil {
		fmt.Println(err)
	}
	p1 := make([]Student, 0)
	for rows.Next() {
		p2 := Student{}
		err := rows.Scan(&p2.Id, &p2.Name, &p2.Mail, &p2.Age)
		if err != nil {
			fmt.Println(err)
			continue
		}
		p1 = append(p1, p2)
	}
	json.NewEncoder(w).Encode(p1)

}

func getSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rows, err := database.Query("SELECT * FROM people WHERE id=" + params["id"])
	if err != nil {
		fmt.Println(err)
	}

	rows.Next()
	p := Student{}
	err = rows.Scan(&p.Id, &p.Name, &p.Mail, &p.Age)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "Here is that Student:\n")
	json.NewEncoder(w).Encode(p)
}

func newSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var std Student
	json.NewDecoder(r.Body).Decode(&std)
	_, err := database.Exec("INSERT INTO people (name, mail,age) VALUES ('" + std.Name + "','" + std.Mail + "','" + std.Age + "')")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, "New Student created:\n")
	json.NewEncoder(w).Encode(std)

}
func changeSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var std Student
	json.NewDecoder(r.Body).Decode(&std)

	params := mux.Vars(r)
	_, err := database.Exec("UPDATE people SET name='" + std.Name + "', mail='" + std.Mail + "', age='" + std.Age + "' WHERE id=" + params["id"])

	if err != nil {
		fmt.Println(err)
	}

	std.Id = params["id"]
	fmt.Fprint(w, "Student Changed:\n")
	json.NewEncoder(w).Encode(std)

}

func deleteSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := database.Exec("DELETE FROM people WHERE id=" + params["id"])

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "Student "+params["id"]+" Deleted!\n")

}

func main() {
	var err error

	fmt.Println("Starting Router...")
	mysql.DeregisterReaderHandler("hsas")
	database, err = sql.Open("mysql", "veles-connect:kL7SMBVgEllwkOma@tcp(176.114.14.32:3306)/veles")
	if err != nil {
		log.Println(err)
	}

	db = append(db, Student{Id: "1", Name: "Oleksandr Kuzan"})

	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	router.HandleFunc("/db", getDB).Methods("GET")
	router.HandleFunc("/db/{id}", getSTD).Methods("GET")
	router.HandleFunc("/db", newSTD).Methods("POST")
	router.HandleFunc("/db/{id}", changeSTD).Methods("PUT")
	router.HandleFunc("/del/{id}", deleteSTD).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
