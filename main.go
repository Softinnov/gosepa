package main

import "encoding/xml"
import "fmt"
import "log"

// Document is the SEPA format for the document containing all transfers
type Document struct {
	XMLName               xml.Name      `xml:"Document"`
	XMLNs                 string        `xml:"xmlns,attr"`
	XMLxsi                string        `xml:"xmlns:xsi,attr"`
	GroupheaderMsgID      string        `xml:"CstmrCdtTrfInitn>GrpHdr>MsgId"`
	GroupheaderCreateDate string        `xml:"CstmrCdtTrfInitn>GrpHdr>CreDtTm"`
	GroupheaderTransacNb  int           `xml:"CstmrCdtTrfInitn>GrpHdr>NbOfTxs"`
	GroupheaderCtrlSum    float32       `xml:"CstmrCdtTrfInitn>GrpHdr>CtrlSum"`
	GroupheaderEmiterName string        `xml:"CstmrCdtTrfInitn>GrpHdr>InitgPty>Nm"`
	PaymentInfoID         string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtInfId"`
	PaymentInfoMethod     string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtMtd"`
	PaymentInfoTransacNb  int           `xml:"CstmrCdtTrfInitn>PmtInf>NbOfTxs"`
	PaymentInfoCtrlSum    float32       `xml:"CstmrCdtTrfInitn>PmtInf>CtrlSum"`
	PaymentTypeInfo       string        `xml:"CstmrCdtTrfInitn>PmtInf>PmtTpInf>SvcLvl>Cd"`
	PaymentExecDate       string        `xml:"CstmrCdtTrfInitn>PmtInf>ReqdExctnDt"`
	PaymentEmiterName     string        `xml:"CstmrCdtTrfInitn>PmtInf>Dbtr>Nm"`
	PaymentEmiterIBAN     string        `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAcct>Id>IBAN"`
	PaymentEmiterBIC      string        `xml:"CstmrCdtTrfInitn>PmtInf>DbtrAgt>FinInstnId>BIC"`
	PaymentCharge         string        `xml:"CstmrCdtTrfInitn>PmtInf>ChrgBr"`
	PaymentTransactions   []Transaction `xml:"CstmrCdtTrfInitn>PmtInf>CdtTrfTxInf"`
}

// Transaction is the transfer SEPA format
type Transaction struct {
	TransacID           string  `xml:"PmtId>InstrId"`
	TransacIDe2e        string  `xml:"PmtId>EndToEndId"`
	TransacAmount       TAmount `xml:"Amt>InstdAmt"`
	TransacCreditorName string  `xml:"Cdtr>Nm"`
	TransacCreditorIBAN string  `xml:"CdtrAcct>Id>IBAN"`
	TransacRegulatory   string  `xml:"RgltryRptg>Dtls>Cd"`
	TransacMotif        string  `xml:"RmtInf>Ustrd"`
}

// TAmount is the transaction amount with its currency
type TAmount struct {
	Amount   float32 `xml:",chardata"`
	Currency string  `xml:"Ccy,attr"`
}

func main() {
	// Initialize doc example
	doc := &Document{
		XMLNs:                 "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03",
		XMLxsi:                "http://www.w3.org/2001/XMLSchema-instance",
		GroupheaderMsgID:      "VIR201705",
		GroupheaderCreateDate: "2017-05-01T12:00:00", // format : AAAA-MM-JJTHH:HH:SS
		GroupheaderTransacNb:  3,
		GroupheaderCtrlSum:    187654.32,
		GroupheaderEmiterName: "Franz Holzapfel GMBH",
		PaymentInfoID:         "VIR201705",
		PaymentInfoMethod:     "TRF",        // always TRF
		PaymentInfoTransacNb:  3,            // same as GroupheaderTransacNb
		PaymentInfoCtrlSum:    187654.32,    // same as GroupheaderCtrlSum
		PaymentTypeInfo:       "SEPA",       // always SEPA
		PaymentExecDate:       "2017-05-03", // format : AAAA-MM-JJ
		PaymentEmiterName:     "Franz Holzapfel GMBH",
		PaymentEmiterIBAN:     "AT611904300234573201",
		PaymentEmiterBIC:      "BKAUATWW",
		PaymentCharge:         "SLEV", // always SLEV
	}

	doc.PaymentTransactions = append(doc.PaymentTransactions, Transaction{
		TransacID:           "F201705",
		TransacIDe2e:        "F201705",
		TransacAmount:       TAmount{Amount: 70000, Currency: "EUR"},
		TransacCreditorName: "DEF Electronics",
		TransacCreditorIBAN: "GB29NWBK60161331926819",
		TransacRegulatory:   "150", // always 150
		TransacMotif:        "F201705",
	})

	str, err := xml.MarshalIndent(doc, "", "  ")
	//str, err := xml.Marshal(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", str)
}
