package gosepa

import (
	"encoding/xml"
	"errors"
	"math/big"
	"time"
)

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

// InitDoc fixes every constants in the document + emiter informations
func (doc *Document) InitDoc(msgID string, creationDate string, executionDate string, emiterName string, emiterIBAN string, emiterBIC string) error {
	_, err := time.Parse("2006-01-02T15:04:05", creationDate) // format : AAAA-MM-JJTHH:HH:SS
	if err != nil {
		return err
	}
	_, err = time.Parse("2006-01-02", executionDate) // format : AAAA-MM-JJ
	if err != nil {
		return err
	}
	if !IsValid(emiterIBAN) {
		return errors.New("Invalid IBAN")
	}
	doc.XMLNs = "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"
	doc.XMLxsi = "http://www.w3.org/2001/XMLSchema-instance"
	doc.GroupheaderMsgID = msgID
	doc.PaymentInfoID = msgID
	doc.GroupheaderCreateDate = creationDate
	doc.PaymentExecDate = executionDate
	doc.GroupheaderEmiterName = emiterName
	doc.PaymentEmiterName = emiterName
	doc.PaymentEmiterIBAN = emiterIBAN
	doc.PaymentEmiterBIC = emiterBIC
	doc.PaymentInfoMethod = "TRF" // always TRF
	doc.PaymentTypeInfo = "SEPA"  // always SEPA
	doc.PaymentCharge = "SLEV"    // always SLEV
	return nil
}

// AddTransaction adds a transfer transaction and adjust the transaction number and the sum control
func (doc *Document) AddTransaction(id string, amount float32, currency string, creditorName string, creditorIBAN string) error {
	if !IsValid(creditorIBAN) {
		return errors.New("Invalid IBAN")
	}
	if amount != float32(int(amount*100))/100 {
		return errors.New("Amount 2 decimals only")
	}
	doc.PaymentTransactions = append(doc.PaymentTransactions, Transaction{
		TransacID:           id,
		TransacIDe2e:        id,
		TransacMotif:        id,
		TransacAmount:       TAmount{Amount: amount, Currency: currency},
		TransacCreditorName: creditorName,
		TransacCreditorIBAN: creditorIBAN,
		TransacRegulatory:   "150", // always 150
	})
	doc.GroupheaderTransacNb++
	doc.PaymentInfoTransacNb++
	doc.GroupheaderCtrlSum += amount
	doc.PaymentInfoCtrlSum += amount
	return nil
}

// Serialize returns the xml document in byte stream
func (doc *Document) Serialize() ([]byte, error) {
	return xml.Marshal(doc)
}

// PrettySerialize returns the indented xml document in byte stream
func (doc *Document) PrettySerialize() ([]byte, error) {
	return xml.MarshalIndent(doc, "", "  ")
}

// IsValid IBAN
func IsValid(iban string) bool {
	i := new(big.Int)
	t := big.NewInt(10)
	if len(iban) < 4 || len(iban) > 34 {
		return false
	}
	for _, v := range iban[4:] + iban[:4] {
		switch {
		case v >= 'A' && v <= 'Z':
			ch := v - 'A' + 10
			i.Add(i.Mul(i, t), big.NewInt(int64(ch/10)))
			i.Add(i.Mul(i, t), big.NewInt(int64(ch%10)))
		case v >= '0' && v <= '9':
			i.Add(i.Mul(i, t), big.NewInt(int64(v-'0')))
		case v == ' ':
		default:
			return false
		}
	}
	return i.Mod(i, big.NewInt(97)).Int64() == 1
}
