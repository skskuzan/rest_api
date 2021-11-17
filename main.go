package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"rest_api/database"
	"rest_api/grpc/pkg/api"
	"rest_api/grpc/pkg/crud"
	"rest_api/vars"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/skskuzan/rest_api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
)

// @title Student App API
// @version 1.0
// @description API Server for Students

// host: localhost:5000
// @BasePath /

// securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func getDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Current Database:\n")
	rows, err := vars.Database.Query("select * from people")

	if err != nil {
		fmt.Println(err)
	}
	p1 := make([]vars.Student, 0)
	for rows.Next() {
		p2 := vars.Student{}
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
	rows, err := vars.Database.Query("SELECT * FROM people WHERE id=" + params["id"])
	if err != nil {
		fmt.Println(err)
	}

	rows.Next()
	p := vars.Student{}
	err = rows.Scan(&p.Id, &p.Name, &p.Mail, &p.Age)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "Here is that Student:\n")
	json.NewEncoder(w).Encode(p)
}

func newSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var std vars.Student
	json.NewDecoder(r.Body).Decode(&std)
	_, err := vars.Database.Exec("INSERT INTO people (name, mail,age) VALUES ('" + std.Name + "','" + std.Mail + "','" + std.Age + "')")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, "New Student created:\n")
	json.NewEncoder(w).Encode(std)

}
func changeSTD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var std vars.Student
	json.NewDecoder(r.Body).Decode(&std)

	params := mux.Vars(r)
	_, err := vars.Database.Exec("UPDATE people SET name='" + std.Name + "', mail='" + std.Mail + "', age='" + std.Age + "' WHERE id=" + params["id"])

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
	_, err := vars.Database.Exec("DELETE FROM people WHERE id=" + params["id"])

	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "Student "+params["id"]+" Deleted!\n")

}

func Parallel(s *grpc.Server) {
	l, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	fmt.Println("Starting Router...")
	mysql.DeregisterReaderHandler("hsas")
	vars.Database, err = sql.Open("mysql", "veles-connect:kL7SMBVgEllwkOma@tcp(176.114.14.32:3306)/veles")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()
	srv := &crud.GRPCServer{}
	api.RegisterCRUDServer(s, srv)

	database.SetupDBConn()

	vars.DB = append(vars.DB, vars.Student{Id: "1", Name: "Oleksandr Kuzan"})

	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	router.HandleFunc("/db", getDB).Methods("GET")
	router.HandleFunc("/db/{id}", getSTD).Methods("GET")
	router.HandleFunc("/db", newSTD).Methods("POST")
	router.HandleFunc("/db/{id}", changeSTD).Methods("PUT")
	router.HandleFunc("/del/{id}", deleteSTD).Methods("DELETE")

	go Parallel(s)
	log.Fatal(http.ListenAndServe(":5000", router))
}
