package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

/*数字一文字を漢数字に変換する関数*/
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

/*数字の文字列を漢数字の文字列に変換する関数*/
func number2kanji(numStr string) (kanji string, err error) {
	if numStr == "0" {
		kanji = string(singleTrans('0'))
		return
	}
	if len(numStr) > 16 {
		err = errors.New("number2kanji: parsing \"" + numStr + "\": value out of range")
	}

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil { return }
	numStr = fmt.Sprintf("%016d", num) //入力文字列を0で埋める

	var runes []rune //文字列を一旦格納するためにスライス
	littleUnits := [3]rune{'千', '百', '拾'}
	indexLittleUnits := 0
	bigUnits := [4]rune{'兆', '億', '万', '_'} //'_'はダミー(兆より大きい単位の実装を用意にするため)
	indexBigUnits := 0
	zeroFlag := true //0000万のようにならないようにするため利用するフラグ
	for i, c := range numStr {
		if c != '0' { //"零千"のようにしないため
			zeroFlag = false
			if indexLittleUnits == 3 { //千、百、拾がつかない桁なら
				runes = append(runes, singleTrans(c))
			} else {
				runes = append(runes, singleTrans(c), littleUnits[indexLittleUnits]) //千、百、拾をつける
			}
		}
		if indexLittleUnits == 3 {	  //千、百、拾の配列が一周したら
			indexLittleUnits = 0	  //千に戻す
		} else { indexLittleUnits++ } //それ以外は次の単位に進める

		if i % 4 == 3 {	//兆、億、万をつけるタイミングなら
			if indexBigUnits != 3 && !zeroFlag { runes = append(runes, bigUnits[indexBigUnits]) } //兆、億、万をつける
			zeroFlag = true
			indexBigUnits++
		}
	}

	kanji = string(runes) //できたスライスをstringへ変換
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
