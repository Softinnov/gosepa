package main

import (
	"fmt"
	"github.com/softinnov/gosepa/sepa"
	"log"
)

func main() {

	doc := &sepa.Document{}
	err := doc.InitDoc("MSGID", "2017-06-07T14:39:33", "2017-06-09", "Emiter Name", "FR1420041010050500013M02606", "BKAUATWW")
	if err != nil {
		log.Fatal("can't create sepa document : ", err)
	}

	err = doc.AddTransaction("F201705", 70000, "EUR", "DEF Electronics", "GB29NWBK60161331926819")
	if err != nil {
		log.Fatal("can't add transaction in the sepa document : ", err)
	}

	res, err := doc.PrettySerialize()
	if err != nil {
		log.Fatal("can't get the xml doc : ", err)
	}

	fmt.Println(string(res))
}
