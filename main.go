//http通信を用いたAPI

package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*数字一文字を漢数字に変換する関数*/
func num2kanjiSingle(input rune) (output rune, err error) {
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
	default:
		err = errors.New("num2kanjiSingle: " + string(input) + ": cannot translate")
		return
	}
	return
}

/*漢数字一文字を数字に変換する関数*/
func kanji2numSingle(input rune) (output rune, err error) {
	switch input {
	case '零': output = '0'
	case '壱': output = '1'
	case '弐': output = '2'
	case '参': output = '3'
	case '四': output = '4'
	case '五': output = '5'
	case '六': output = '6'
	case '七': output = '7'
	case '八': output = '8'
	case '九': output = '9'
	default:
		err = errors.New("num2kanjiSingle: " + string(input) + ": cannot translate")
		return
	}
	return
}

/*数字の文字列を漢数字の文字列に変換する関数*/
func number2kanji(numStr string) (kanji string, err error) {
	if numStr == "0" {
		tmp, _ := num2kanjiSingle('0')
		kanji = string(tmp)
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
	var tmpRune rune
	for i, c := range numStr {
		if c != '0' { //"零千"のようにしないため
			zeroFlag = false
			tmpRune, err = num2kanjiSingle(c)
			if err != nil { return }
			if indexLittleUnits == 3 { //千、百、拾がつかない桁なら
				runes = append(runes, tmpRune)
			} else {
				runes = append(runes, tmpRune, littleUnits[indexLittleUnits]) //千、百、拾をつける
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

/*漢数字の文字列を数字の文字列に変換する関数*/
func kanji2number(kanji string) (numStr string, err error) {
	num := 0
	tmpLittleNum := 0; tmpBigNum := 0 //拾、百、千の計算、万、億、兆の計算を行うのに使う変数
	littleUnits := "拾百千"
	bigUnits := "万億兆"

	var tmpRune rune
	for _, c := range kanji {
		if strings.Contains(bigUnits, string(c)) { //万、億、兆のとき
			tmpBigNum += tmpLittleNum
			switch c {
			case '万': num += tmpBigNum * 10000
			case '億': num += tmpBigNum * 100000000
			case '兆': num += tmpBigNum * 1000000000000
			}
			tmpBigNum = 0
			tmpLittleNum = 0
		} else if strings.Contains(littleUnits, string(c)) { //拾、百、千のとき
			switch c {
			case '拾': tmpBigNum += tmpLittleNum * 10
			case '百': tmpBigNum += tmpLittleNum * 100
			case '千': tmpBigNum += tmpLittleNum * 1000
			}
			tmpLittleNum = 0
		} else {
			tmpRune, err = kanji2numSingle(c)
			if err != nil { return }
			tmpLittleNum = int(tmpRune - '0')
		}
	}
	num += tmpBigNum
	num += tmpLittleNum

	numStr = strconv.Itoa(num) //できた値を文字列に変換
	return
}

/* /number2kanji/{number} のエンドポイントにリクエストがきたときの処理*/
func handleNumber2kanji(w http.ResponseWriter, r *http.Request) {
	log.Println("handleNumber2kanji called")
	vars := mux.Vars(r)
	kanji, err := number2kanji(vars["number"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("translate: ", vars["number"], "to", kanji)
	if _, err := fmt.Fprint(w, kanji); err != nil {
		log.Fatal(err)
	}
}

/* /kanji2number/{kanji} のエンドポイントにリクエストがきたときの処理*/
func handleKanji2number(w http.ResponseWriter, r *http.Request) {
	log.Println("handleKanji2number called")
	vars := mux.Vars(r)
	numStr, err := kanji2number(vars["kanji"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("translate: ", vars["kanji"], "to", numStr)
	if _, err := fmt.Fprint(w, numStr); err != nil {
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
