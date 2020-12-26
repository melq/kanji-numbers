package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func returnHello(w http.ResponseWriter, r *http.Request) {
	log.Println("returnHello called")
	if n, err := fmt.Fprint(w, "hello world!"); err != nil {
		log.Fatal(n, err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/hello", returnHello).Methods("GET")

	log.Println("Listen Server ...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
