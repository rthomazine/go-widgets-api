package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
)

func createRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", authenticationHandler).Methods("POST")
	router.Handle("/api", negroni.New(
		negroni.HandlerFunc(validateTokenHandler),
		negroni.Wrap(http.HandlerFunc(apiHandler)),
	))
	log.Println("Initialized router")

	return router
}
