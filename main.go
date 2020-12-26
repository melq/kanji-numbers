package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func singleTrans(input rune) (output rune) {
	switch input {
	case '0': output = '零'
	case '1': output = '壱'
	case '2': output = '弐'
	case '3': output = '参'
	case '4': output = '四'
	case '5': output = '五'
	case '6': output = '六'
	case '7': output = '七'
	case '8': output = '八'
	case '9': output = '九'
	}
	return
}

func number2kanji(numStr string) (kanji string, err error) {
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil { return }
	if len(numStr) > 16 {
		err = errors.New("number2kanji: parsing \"" + numStr + "\": value out of range")
	}

	numStr = fmt.Sprintf("%016d", num)

	var runes []rune
	littleUnits := []rune{'千', '_', '拾', '百'}
	indexLittleUnits := 0
	bigUnits := []rune{'_', '万', '億', '兆'}
	indexBigUnits := len(bigUnits) - 1
	zeroFlag := true
	for _, c := range numStr {
		if c != '0' {
			zeroFlag = false
			if indexLittleUnits == 1 {
				runes = append(runes, singleTrans(c))
			} else {
				runes = append(runes, singleTrans(c), littleUnits[indexLittleUnits])
			}
		}

		if indexLittleUnits == 1 {
			if indexBigUnits != 0 && !zeroFlag { runes = append(runes, bigUnits[indexBigUnits]) }
			zeroFlag = true
			indexBigUnits--
		}

		if indexLittleUnits == 0 {
			indexLittleUnits = 3
		} else { indexLittleUnits-- }
	}

	kanji = string(runes)
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
