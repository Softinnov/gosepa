# gosepa

[![Go Report Card](https://goreportcard.com/badge/github.com/softinnov/gosepa)](https://goreportcard.com/report/github.com/softinnov/gosepa)

gosepa is a sepa xml file generator written in Go compatible with pain.002.001.03 schema (Customer Credit Transfer Initiation V03).

This generator uses shortcuts to simplify the norm implementation : for example, there is only one id used for several differents references (I never have to deal with multiple references)

## Install

```console
$ go get github.com/softinnov/gosepa/sepa
```

## Usage

```go
    package main

    import (
        "fmt"
        "github.com/softinnov/gosepa/sepa"
        "log"
    )

    func main() {

        doc := &sepa.Document{}
        err := doc.InitDoc("MSGID", "2017-06-07T14:39:33", "2017-06-09", "Emiter Name", "FR1420041010050500013M02606", "BKAUATWWP")
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
```

## Tests

Unit test the go way :

```console
$ go test -v
```

You can use any xsd validation tool. I use xmllint from libxml.

```console
$ sudo apt install libxml2-utils
```

You have to generate a file so xmllint can check it. From the sample in the 'tests' folder :

```console
$ go run usage.go > test.xml
$ xmllint --noout --schema pain.001.001.03.xsd test.xml
```

## Ressources

* [sepa xsd](https://www.iso20022.org/message_archive.page)
* [go xml](https://golang.org/pkg/encoding/xml/)
