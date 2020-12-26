package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func number2kanji(numStr string) (kanji string, err error) {
	var num int64
	num, err = strconv.ParseInt(numStr, 10, 64)
	if err != nil { return }
	if num > 9999999999999999 {
		err = errors.New("number2kanji: parsing \"" + strconv.FormatInt(num, 10) + "\": value out of range")
	}

	kanji = "kanji: " + strconv.FormatInt(num, 10)
	return
}

func kanji2number(kanji string) (numStr string) {
	numStr = "1234"
	return
}

func handleNumber2kanji(w http.ResponseWriter, r *http.Request) {
	log.Println("handleNumber2kanji called")
	var kanji string

	vars := mux.Vars(r)
	kanji, err := number2kanji(vars["number"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := fmt.Fprint(w, kanji); err != nil {
		log.Fatal(err)
	}
}

func handleKanji2number(w http.ResponseWriter, r *http.Request) {
	log.Println("handleKanji2number called")
	vars := mux.Vars(r)

	if _, err := fmt.Fprint(w, vars["kanji"]); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/number2kanji/{number}", handleNumber2kanji).Methods("GET")
	router.HandleFunc("/v1/kanji2number/{kanji}", handleKanji2number).Methods("GET")

	log.Println("Listen Server ...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
