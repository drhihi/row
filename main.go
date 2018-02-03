package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("the reminder of words"))
}

func main() {

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", indexHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
