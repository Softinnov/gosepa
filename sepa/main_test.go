package sepa

import (
	"strings"
	"testing"
)

func TestDecimalsNumber(t *testing.T) {
	suite := []struct {
		f float64
		n int
	}{
		{0, 0},
		{123.0, 0},
		{144.2, 1},
		{1.123456789, 9},
		{3.1415900000, 5},
		{-1250, 0},
		{-252123.123, 3},
	}
	for _, s := range suite {
		received := DecimalsNumber(s.f)
		expected := s.n
		if received != expected {
			t.Errorf("Expected %v received %v", expected, received)
		}
	}
}
func TestGenerateSepaXML(t *testing.T) {

	var err error

	// targetDoc is a verified valid sepa xml file
	var targetDoc = `<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pain.001.001.03" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><CstmrCdtTrfInitn><GrpHdr><MsgId>VIR201705</MsgId><CreDtTm>2017-05-01T12:00:00</CreDtTm><NbOfTxs>5</NbOfTxs><CtrlSum>170000</CtrlSum><InitgPty><Nm>Franz Holzapfel GMBH</Nm></InitgPty></GrpHdr><PmtInf><PmtInfId>VIR201705</PmtInfId><PmtMtd>TRF</PmtMtd><NbOfTxs>5</NbOfTxs><CtrlSum>170000</CtrlSum><PmtTpInf><SvcLvl><Cd>SEPA</Cd></SvcLvl></PmtTpInf><ReqdExctnDt>2017-05-03</ReqdExctnDt><Dbtr><Nm>Franz Holzapfel GMBH</Nm></Dbtr><DbtrAcct><Id><IBAN>AT611904300234573201</IBAN></Id></DbtrAcct><DbtrAgt><FinInstnId><BIC>BKAUATWW</BIC></FinInstnId></DbtrAgt><ChrgBr>SLEV</ChrgBr><CdtTrfTxInf><PmtId><InstrId>F201705</InstrId><EndToEndId>F201705</EndToEndId></PmtId><Amt><InstdAmt Ccy="EUR">70000</InstdAmt></Amt><Cdtr><Nm>DEF Electronics</Nm></Cdtr><CdtrAcct><Id><IBAN>GB29NWBK60161331926819</IBAN></Id></CdtrAcct><RgltryRptg><Dtls><Cd>150</Cd></Dtls></RgltryRptg><RmtInf><Ustrd>F201705</Ustrd></RmtInf></CdtTrfTxInf><CdtTrfTxInf><PmtId><InstrId>F201706</InstrId><EndToEndId>F201706</EndToEndId></PmtId><Amt><InstdAmt Ccy="EUR">10000</InstdAmt></Amt><Cdtr><Nm>D1F Electronics</Nm></Cdtr><CdtrAcct><Id><IBAN>AT611904300234573201</IBAN></Id></CdtrAcct><RgltryRptg><Dtls><Cd>150</Cd></Dtls></RgltryRptg><RmtInf><Ustrd>F201706</Ustrd></RmtInf></CdtTrfTxInf><CdtTrfTxInf><PmtId><InstrId>F201707</InstrId><EndToEndId>F201707</EndToEndId></PmtId><Amt><InstdAmt Ccy="EUR">20000</InstdAmt></Amt><Cdtr><Nm>D2F Electronics</Nm></Cdtr><CdtrAcct><Id><IBAN>BE62510007547061</IBAN></Id></CdtrAcct><RgltryRptg><Dtls><Cd>150</Cd></Dtls></RgltryRptg><RmtInf><Ustrd>F201707</Ustrd></RmtInf></CdtTrfTxInf><CdtTrfTxInf><PmtId><InstrId>F201708</InstrId><EndToEndId>F201708</EndToEndId></PmtId><Amt><InstdAmt Ccy="EUR">30000</InstdAmt></Amt><Cdtr><Nm>D3F Electronics</Nm></Cdtr><CdtrAcct><Id><IBAN>BG80BNBG96611020345678</IBAN></Id></CdtrAcct><RgltryRptg><Dtls><Cd>150</Cd></Dtls></RgltryRptg><RmtInf><Ustrd>F201708</Ustrd></RmtInf></CdtTrfTxInf><CdtTrfTxInf><PmtId><InstrId>F201709</InstrId><EndToEndId>F201709</EndToEndId></PmtId><Amt><InstdAmt Ccy="EUR">40000</InstdAmt></Amt><Cdtr><Nm>D4F Electronics</Nm></Cdtr><CdtrAcct><Id><IBAN>EE382200221020145685</IBAN></Id></CdtrAcct><RgltryRptg><Dtls><Cd>150</Cd></Dtls></RgltryRptg><RmtInf><Ustrd>F201709</Ustrd></RmtInf></CdtTrfTxInf></PmtInf></CstmrCdtTrfInitn></Document>`

	// our doc
	var sepaDoc = &Document{}

	// Bad format for creation date, expecting AAAA-MM-JJTHH:HH:SS
	err = sepaDoc.InitDoc("", "2017-05-01", "", "", "", "")
	if err == nil {
		t.Error("Expected InitDoc return an error for bad creation date format", "got", err)
	}

	// Bad format for execution date, expecting AAAA-MM-JJ
	err = sepaDoc.InitDoc("", "2017-05-01T22:45:03", "2017-05-01T22:45:03", "", "", "")
	if err == nil {
		t.Error("Expected InitDoc return an error for bad exection date format", "got", err)
	}

	// Bad IBAN
	err = sepaDoc.InitDoc("", "2017-05-01T22:45:03", "2017-05-01", "", "XX12345678901234567", "")
	if err == nil {
		t.Error("Expected InitDoc return an error for bad IBAN", "got", err)
	}

	// Good IBAN
	err = sepaDoc.InitDoc("", "2017-05-01T22:45:03", "2017-05-01", "", "FR1420041010050500013M02606", "")
	if err != nil {
		t.Error("Expected InitDoc return nil for good IBAN", "got", err)
	}

	// Initialize doc test
	err = sepaDoc.InitDoc("VIR201705", "2017-05-01T12:00:00", "2017-05-03", "Franz Holzapfel GMBH", "AT611904300234573201", "BKAUATWW")
	if err != nil {
		t.Error("Expected InitDoc return nil", "got", err)
	}

	// Add Transaction with incorrect IBAN
	err = sepaDoc.AddTransaction("XXX", 0, "XXX", "XXX", "ZZ382200221020145685")
	if err == nil {
		t.Error("Exepected AddTransaction return an error for bad IBAN", "got", err)
	}

	// Add Transaction with incorrect amount (>2 decimals)
	err = sepaDoc.AddTransaction("XXX", 1.234, "XXX", "XXX", "EE382200221020145685")
	if err == nil {
		t.Error("Exepected AddTransaction return an error for bad amount", "got", err)
	}

	// Transactions Test Array
	type testTransac struct {
		id         string
		amount     float64
		currency   string
		debtorName string
		debtorIban string
	}
	TTest := []testTransac{
		{"F201705", 70000, "EUR", "DEF Electronics", "GB29NWBK60161331926819"},
		{"F201706", 10000, "EUR", "D1F Electronics", "AT611904300234573201"},
		{"F201707", 20000, "EUR", "D2F Electronics", "BE62510007547061"},
		{"F201708", 30000, "EUR", "D3F Electronics", "BG80BNBG96611020345678"},
		{"F201709", 40000, "EUR", "D4F Electronics", "EE382200221020145685"},
	}

	// For each transaction, we check that the cumul amount and number of transactions remain correct in header and payment block
	var cumul = float64(0)

	var nb = 0
	for count, transac := range TTest {
		err = sepaDoc.AddTransaction(transac.id, transac.amount, transac.currency, transac.debtorName, transac.debtorIban)
		if err != nil {
			t.Error("Expected AddTransaction return nil", "got", err)
		}
		cumul += transac.amount
		nb = 1 + count
		if sepaDoc.GroupheaderCtrlSum != cumul {
			t.Error("Expected GroupheaderCtrlSum", cumul, "got", sepaDoc.GroupheaderCtrlSum)
		}
		if sepaDoc.PaymentInfoCtrlSum != cumul {
			t.Error("Expected PaymentInfoCtrlSum", cumul, "got", sepaDoc.PaymentInfoCtrlSum)
		}
		if sepaDoc.GroupheaderTransacNb != nb {
			t.Error("Expected GroupheaderTransacNb", nb, "got", sepaDoc.GroupheaderTransacNb)
		}
		if sepaDoc.PaymentInfoTransacNb != nb {
			t.Error("Expected PaymentInfoTransacNb", nb, "got", sepaDoc.PaymentInfoTransacNb)
		}
	}

	// Get the result
	str, err := sepaDoc.Serialize()
	if err != nil {
		t.Error("Expected xml in []byte, got ", err)
	}
	// Ultimate test : compare the all generated doc with the predefined doc
	res := strings.Compare(string(str), targetDoc)
	if res != 0 {
		t.Error("Expected", string(targetDoc), "got", string(str))
	}
}
