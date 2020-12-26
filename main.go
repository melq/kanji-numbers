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
func number2kanji(w http.ResponseWriter, r *http.Request) {
	log.Println("number2kanji called")
	vars := mux.Vars(r)
	if n, err := fmt.Fprint(w, vars["number"]); err != nil {
		log.Fatal(n, err)
	}
}

func kanji2number(w http.ResponseWriter, r *http.Request) {
	log.Println("kanji2number called")
	vars := mux.Vars(r)
	if n, err := fmt.Fprint(w, vars["kanji"]); err != nil {
		log.Fatal(n, err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/hello", returnHello).Methods("GET")
	router.HandleFunc("/v1/number2kanji/{number}", number2kanji).Methods("GET")
	router.HandleFunc("/v1/kanji2number/{kanji}", kanji2number).Methods("GET")

	log.Println("Listen Server ...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
