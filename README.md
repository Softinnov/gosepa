# goSEPA

goSEPA is a sepa xml file generation example written in Go compatible with pain.002.001.03 schema (Customer Credit Transfer Initiation V03).

## run the tests

`
go test -v
`

## check validity

You can use any xsd validation tool. I use xmllint from libxml :

### install xmllint
`
sudo apt install libxml2-utils
`

### validate
`
xmllint --noout --schema pain.001.001.03.xsd test.xml
`

### ressources

* [sepa xsd](https://www.iso20022.org/message_archive.page)
* [go xml](https://golang.org/pkg/encoding/xml/)
