package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/m/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func main() {
	// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// db, err := sql.Open("postgres", psqlconn)
	// CheckErr(err)

	// defer db.Close()

	// insertStmt := `insert into "employee"("Name","EmpId") values('Tohit',14)`
	// _, e := db.Exec(insertStmt)

	// CheckErr(e)

	r := mux.NewRouter()
	r.HandleFunc("/list", handlers.GetList).Methods("GET")
	r.HandleFunc("/add", handlers.AddTask).Methods("POST")
	r.HandleFunc("/addUser", handlers.AddUser).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	// handler = r.Handler(r)
	log.Println("Listening on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("SDFSDF")
		panic(err)
	}
}
