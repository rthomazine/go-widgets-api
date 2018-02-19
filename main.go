package main

import (
	"log"
	"net/http"
	"database/sql"
)

var db *sql.DB

func main() {

	db = openDB()
	defer db.Close()

	router := createRouter()
	log.Println("Ready to serve...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
