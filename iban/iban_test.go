package iban

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewGerman(t *testing.T) {
	bankID := "10090044"
	accountID := "532013018"
	var result IBAN

	result, err := NewGerman(bankID, accountID)

	if err != nil {
		t.Logf("Expected no error, got %T:%v", err, err)
		t.Fail()
	}

	expectedResult := "DE10100900440532013018"

	if string(result) != expectedResult {
		t.Logf("Expected result to equal %q, got %q", expectedResult, result)
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		testCases := []struct {
			country     string
			countryCode string
			bban        string
			result      IBAN
		}{
			{"Albania", "AL", "202111090000000005012075", "AL06202111090000000005012075"},
			{"Andorra", "AD", "00060004451247870930", "AD1000060004451247870930"},
			{"Azerbaijan", "AZ", "UBAZ04003214540060AZN001", "AZ04UBAZ04003214540060AZN001"},
			{"Bahrain", "BH", "CITI00001077181611", "BH02CITI00001077181611"},
			{"Belgium", "BE", "096920886089", "BE45096920886089"},
			{"Bosnia and Herzegovina", "BA", "1011606058553319", "BA391011606058553319"},
			{"Brasilia", "BR", "00000000010670000117668C1", "BR0200000000010670000117668C1"},
			{"British Virgin Islands", "VG", "NOSC0000000005002993", "VG48NOSC0000000005002993"},
			{"Bulgaria", "BG", "RZBB91551002755190", "BG02RZBB91551002755190"},
			{"Costa Rica", "CR", "015202220005614288", "CR79015202220005614288"},
			{"Denmark", "DK", "20005036459478", "DK0220005036459478"},
			{"Germany", "DE", "100500000024290661", "DE02100500000024290661"},
			{"Dominican Republic", "DO", "BCBH00000000011003290022", "DO22BCBH00000000011003290022"},
			{"El Salvador", "SV", "ACAT00000000000000123123", "SV43ACAT00000000000000123123"},
			{"Estonia", "EE", "1700017000459042", "EE021700017000459042"},
			{"Faeroe Islands", "FO", "91810001441878", "FO1291810001441878"},
			{"Finland", "FI", "10403500314392", "FI0210403500314392"},
			{"France", "FR", "20041000016219433J02076", "FR0220041000016219433J02076"},
			{"Georgia", "GE", "TB7523045063700002", "GE02TB7523045063700002"},
			{"Gibraltar", "GI", "BARC020452163087000", "GI04BARC020452163087000"},
			{"Greece", "GR", "01102160000021661309175", "GR0201102160000021661309175"},
			{"Greenland", "GL", "64710001504964", "GL2664710001504964"},
			{"Great Britain", "GB", "CITI18500811417983", "GB11CITI18500811417983"},
			{"Guatemala", "GT", "CITI01010000000004146026", "GT24CITI01010000000004146026"},
			{"Iraq", "IQ", "CBIQ861800101010500", "IQ20CBIQ861800101010500"},
			{"Ireland", "IE", "BOFI90008413113207", "IE02BOFI90008413113207"},
			{"Iceland", "IS", "0116381002305610911109", "IS040116381002305610911109"},
			{"Israel", "IL", "0108380000002149431", "IL020108380000002149431"},
			{"Italy", "IT", "K0310412701000000820420", "IT43K0310412701000000820420"},
			{"Jordan", "JO", "SCBL1260000000018525836101", "JO02SCBL1260000000018525836101"},
			{"Kazakhstan", "KZ", "319C010005569698", "KZ04319C010005569698"},
			{"Qatar", "QA", "QNBA000000000060565452001", "QA03QNBA000000000060565452001"},
			{"Kosovo", "XK", "1301001002074155", "XK051301001002074155"},
			{"Croatia", "HR", "23400093216312031", "HR0223400093216312031"},
			{"Kuwait", "KW", "NBOK0000000000001000614589", "KW02NBOK0000000000001000614589"},
			{"Latvia", "LV", "HABA0551007820897", "LV02HABA0551007820897"},
			{"Lebanon", "LB", "001400000302300023018319", "LB02001400000302300023018319"},
			{"Liechtenstein", "LI", "08800000022875748", "LI0308800000022875748"},
			{"Lithuania", "LT", "7300010134441147", "LT027300010134441147"},
			{"Luxembourg", "LU", "0019175546294000", "LU020019175546294000"},
			{"Malta", "MT", "VALL22013000000040013752732", "MT02VALL22013000000040013752732"},
			{"Mauritania", "MR", "00012000010000009880016", "MR1300012000010000009880016"},
			{"Mauritius", "MU", "MCBL0901000001879025000USD", "MU03MCBL0901000001879025000USD"},
			{"Macedonia", "MK", "200000625758632", "MK07200000625758632"},
			{"Moldova", "MD", "MO2224ASV41884097100", "MD14MO2224ASV41884097100"},
			{"Monaco", "MC", "12739000710075018000P14", "MC2412739000710075018000P14"},
			{"Montenegro", "ME", "505120000000466170", "ME25505120000000466170"},
			{"Netherlands", "NL", "ABNA0457180536", "NL02ABNA0457180536"},
			{"Norway", "NO", "39916835985", "NO0239916835985"},
			{"Austria", "AT", "1100000622888600", "AT021100000622888600"},
			{"Pakistan", "PK", "SCBL0000001925518401", "PK02SCBL0000001925518401"},
			{"Palestine", "PS", "ARAB000000009040781605610", "PS06ARAB000000009040781605610"},
			{"Poland", "PL", "103000190109780401676562", "PL02103000190109780401676562"},
			{"Portugal", "PT", "003600409911001102673", "PT50003600409911001102673"},
			{"Romania", "RO", "BRDE445SV75163474450", "RO02BRDE445SV75163474450"},
			{"Saint Lucia", "LC", "HEMM000100010012001200023015", "LC55HEMM000100010012001200023015"},
			{"San Marino", "SM", "U0854009803000030174419", "SM07U0854009803000030174419"},
			{"São Tomé and Príncipe", "ST", "000200000289355710148", "ST23000200000289355710148"},
			{"Saudi Arabia", "SA", "20000002480647579940", "SA0220000002480647579940"},
			{"Sweden", "SE", "30000000030301099952", "SE0230000000030301099952"},
			{"Swiss", "CH", "0020720710117540C", "CH020020720710117540C"},
			{"Serbia", "RS", "105008054113238018", "RS35105008054113238018"},
			{"Seychelles", "SC", "NOVH00000021002035257028SCR", "SC74NOVH00000021002035257028SCR"},
			{"Slovak Republic", "SK", "02000000003679748552", "SK0202000000003679748552"},
			{"Slowenia", "SI", "011006000005649", "SI56011006000005649"},
			{"Spain", "ES", "21000555370200853027", "ES1321000555370200853027"},
			{"Timor-Leste", "TL", "0030000000025923744", "TL380030000000025923744"},
			{"Turkey", "TR", "0001000201529153355002", "TR020001000201529153355002"},
			{"Czech Republic", "CZ", "01000000199216760237", "CZ0201000000199216760237"},
			{"Tunisia", "TN", "01026067111999766058", "TN5901026067111999766058"},
			{"Ukraine", "UA", "3052990004149497803982794", "UA123052990004149497803982794"},
			{"Hungary", "HU", "116000060000000064247067", "HU02116000060000000064247067"},
			{"United Arab Emirates", "AE", "0090004001079346500", "AE020090004001079346500"},
			{"Belarus", "BY", "AKBB10100000002966000000", "BY86AKBB10100000002966000000"},
			{"Cyprus", "CY", "002001950000357009822416", "CY02002001950000357009822416"},
		}

		for _, tt := range testCases {
			t.Run(tt.country, func(t *testing.T) {
				result, err := New(tt.countryCode, tt.bban)
				if err != nil {
					t.Logf("Expected no error, got %T:%v", err, err)
					t.Fail()
				}

				if result != tt.result {
					t.Logf("Expected result to equal %q, got %q", tt.result, result)
					t.Fail()
				}
			})
		}
	})
	t.Run("uncommon input", func(t *testing.T) {
		testCases := []struct {
			desc        string
			countryCode string
			bban        string
			result      IBAN
		}{
			{
				desc:        "lowercase country code",
				countryCode: "de",
				bban:        "100500000024290661",
				result:      "DE02100500000024290661",
			},
			{
				desc:        "mixed case countryCode",
				countryCode: "Gb",
				bban:        "CITI18500811417983",
				result:      "GB11CITI18500811417983",
			},
			{
				desc:        "lowercase bban",
				countryCode: "GB",
				bban:        "CITI18500811417983",
				result:      "GB11CITI18500811417983",
			},
			{
				desc:        "mixed case bban",
				countryCode: "GB",
				bban:        "CiTi18500811417983",
				result:      "GB11CITI18500811417983",
			},
			{
				desc:        "mixed case countryCode and bban",
				countryCode: "gB",
				bban:        "cItI18500811417983",
				result:      "GB11CITI18500811417983",
			},
		}

		for _, tt := range testCases {
			t.Run(tt.desc, func(t *testing.T) {
				result, err := New(tt.countryCode, tt.bban)
				if err != nil {
					t.Logf("Expected no error, got %T:%v", err, err)
					t.Fail()
				}

				if result != tt.result {
					t.Logf("Expected result to equal %q, got %q", tt.result, result)
					t.Fail()
				}
			})
		}
	})
	t.Run("errors", func(t *testing.T) {
		testCases := []struct {
			desc        string
			countryCode string
			bban        string
		}{
			{
				desc:        "empty country code",
				countryCode: "",
				bban:        "1234567890123456789",
			},
			{
				desc:        "empty country code",
				countryCode: "ABCDE",
				bban:        "1234567890123456789",
			},
			{
				desc:        "empty country code",
				countryCode: "A",
				bban:        "1234567890123456789",
			},
			{
				desc:        "too long BBAN",
				countryCode: "AT",
				bban:        "1234567890123456789012345678901234567890",
			},
		}

		for _, tt := range testCases {
			_, err := New(tt.countryCode, tt.bban)
			if err == nil {
				t.Logf("Expected error, got nil")
				t.Fail()
			}
		}
	})
}

func TestIsValid(t *testing.T) {
	t.Run("valid german IBAN", func(t *testing.T) {
		iban, _ := NewGerman("10090044", "0532013018")

		ok := IsValid(iban)
		if !ok {
			t.Logf("Expected iban to be valid")
			t.Fail()
		}
	})
	t.Run("valid british IBAN", func(t *testing.T) {
		iban, _ := New("GB", "CITI18500811417983")

		ok := IsValid(iban)
		if !ok {
			t.Logf("Expected iban to be valid")
			t.Fail()
		}
	})
	t.Run("valid uncommon IBAN", func(t *testing.T) {
		iban, _ := New("GB", "CITI18500811417983")
		iban = IBAN(fmt.Sprintf(
			"%s%s%s",
			strings.ToLower(iban.CountryCode()),
			iban.ProofNumber(),
			strings.ToLower(iban.BBAN()),
		))

		ok := IsValid(iban)
		if !ok {
			t.Logf("Expected iban to be valid")
			t.Fail()
		}
	})
	t.Run("invalid german IBAN", func(t *testing.T) {
		iban := IBAN("DE9910090044053201301812345678901234567890")

		ok := IsValid(iban)
		if ok {
			t.Logf("Expected iban to be invalid")
			t.Fail()
		}
	})
	t.Run("invalid german IBAN", func(t *testing.T) {
		iban := IBAN("DE99100900440532013018")

		ok := IsValid(iban)
		if ok {
			t.Logf("Expected iban to be invalid")
			t.Fail()
		}
	})
	t.Run("invalid IBAN", func(t *testing.T) {
		iban, _ := New("GB", "CITI18500811417983")
		iban = IBAN(fmt.Sprintf(
			"%s%s%s",
			"12",
			iban.ProofNumber(),
			strings.ToLower(iban.BBAN()),
		))

		ok := IsValid(iban)
		if ok {
			t.Logf("Expected iban to be invalid")
			t.Fail()
		}
	})
}

func TestIbanBBAN(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	bban := iban.BBAN()

	expectedBban := "100900440532013018"

	if bban != expectedBban {
		t.Logf("Expected bankId to equal %q, got %q\n", expectedBban, bban)
		t.Fail()
	}
}

func TestIbanBankId(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	bankID := iban.BankID()

	expectedBankID := "10090044"

	if bankID != expectedBankID {
		t.Logf("Expected bankId to equal %q, got %q\n", expectedBankID, bankID)
		t.Fail()
	}
}

func TestIbanAccountId(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	accountID := iban.AccountID()

	expectedAccountID := "532013018"

	if accountID != expectedAccountID {
		t.Logf("Expected accountId to equal %q, got %q\n", expectedAccountID, accountID)
		t.Fail()
	}
}

func TestIbanCountry(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	country := iban.CountryCode()

	expectedCountry := "DE"

	if country != expectedCountry {
		t.Logf("Expected country to equal %q, got %q\n", expectedCountry, country)
		t.Fail()
	}
}

func TestIbanProofNumber(t *testing.T) {
	iban := IBAN("DE10100900440532013018")

	proofNumber := iban.ProofNumber()

	expectedProofNumber := "10"

	if proofNumber != expectedProofNumber {
		t.Logf("Expected proofNumber to equal %q, got %q\n", expectedProofNumber, proofNumber)
		t.Fail()
	}
}

func TestIbanString(t *testing.T) {
	iban := IBAN("GB11CITI18500811417983")

	s := iban.String()

	expected := "GB11CITI18500811417983"

	if expected != s {
		t.Logf("Expected iban.String() to equal %q, got %q", expected, s)
		t.Fail()
	}
}

func TestPrint(t *testing.T) {
	iban := IBAN("GB11CITI18500811417983")

	printed := Print(iban)

	expected := "GB11 CITI 1850 0811 4179 83"

	if expected != printed {
		t.Logf("Expected printed iban to equal %q, got %q", expected, printed)
		t.Fail()
	}
}
