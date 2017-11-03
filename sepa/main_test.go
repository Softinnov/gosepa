package sepa

import (
	"strings"
	"testing"
)

func TestCumul(t *testing.T) {
	var s = &Document{}
	s.InitDoc("", "2017-05-01T22:45:03", "2017-05-01", "", "FR1420041010050500013M02606", "")
	TTest := []float64{55, 140, 77, 105, 140, 76.3, 164.8, 62.3, 29.3, 125.3, 70, 78.22, 252.9, 35, 70, 173.6, 60.9, 63, 126, 215.6, 12.5, 35, 257.6, 75, 30, 72.5, 259.5, 302.62, 120.4, 35, 173.6, 104.54, 119, 22.5, 80.5, 135.8, 161.85, 1199.86, 32.5, 70, 140, 633.92, 159.6, 35, 196, 97.3, 90.3, 144.9, 258.7, 374.13, 27.5, 1575, 282.1, 56, 105, 57.4, 51.8, 56, 801.5, 66.99, 98.5, 212.8, 35, 109.9, 35, 269.5, 327.6, 224, 38.5, 35, 266, 256.2, 102.9, 201.6, 0.34, 35, 35, 341.6, 21, 217, 35.1, 19, 114, 25, 277.9, 70, 140, 21, 67.5, 41.3, 134.4, 143.36, 74, 21, 24, 27.07, 208.6, 43.75, 70, 58.8, 38.15, 61.5, 147, 378.8, 16.5, 52.5, 24.5, 60.2, 72.84, 175, 17.5, 70, 231.6, 161, 49, 70, 45.5, 291.2, 41.3, 35, 186.2, 154, 70, 35, 70, 35, 230, 119, 70, 20, 70, 175, 36.5, 217, 35, 52, 31.3, 109.2, 35, 24.5, 13.5, 63.5, 111.3, 60.2, 103, 203, 143.5, 35, 57.5, 35, 125.3, 175, 138.6, 153.82, 120.4, 62.5, 35.52, 63.5, 129.5, 70, 175, 224, 70, 126, 140, 35, 140, 25.5, 7.98, 70, 35, 65.2, 105, 77, 35, 98, 225.5, 38.5, 35, 158, 72.8, 147, 50, 210, 385, 28, 202.3, 128.8, 39.2, 117.6, 326, 30}
	cumul := float64(24443.66)
	//TTest := []float64{593.3, 164.8}
	//cumul := float64(758.1)
	for _, m := range TTest {
		s.AddTransaction("", m, "EUR", "", "GB29NWBK60161331926819")
	}
	if s.GroupheaderCtrlSum != cumul {
		t.Error("Expected GroupheaderCtrlSum", cumul, "got", s.GroupheaderCtrlSum)
	}
}
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
		t.Error("Expected InitDoc return an error for bad execution date format", "got", err)
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
