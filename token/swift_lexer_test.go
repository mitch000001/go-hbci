package token

import (
	"testing"
)

const githubIssue13TestData = "\r\n" +
	":20:MT940-1803060458\r\n" +
	":21:NONREF\r\n" +
	":25:20040000/12345678EUR\r\n" +
	":28C:0/13\r\n:60M:C170201EUR1234,56\r\n" +
	":61:1702010201DR86,40NMSCNONREF//POS 8888888888\r\n" +
	":86:005?20LASTSCHRIFT/BELAST.?888888888888 8888888884REFERE?22NZ HVV \r\n" +
	"A?23BO?24END-TO-END-REF.:?888888888888 8888888884?26CORE / MANDAT\r\n" +
	"SREF.:?27VMH008888880001?28GL\xc4UBIGER-ID:?29DE88888888888888888?32H\r\n" +
	"AMBURGER HOCHBAHN AG?60Ref. IL888888G8888888/6716\r\n" +
	":61:1702010201DR24,00NMSCNONREF//POS 3409790600\r\n" +
	":86:005?20LASTSCHRIFT/BELAST.?21110865 BEITRAG MITGLIED 888?228888?23\r\n" +
	"END-TO-END-REF.:?8888888888888ZV888888Z?25CORE / MANDATSREF.:?26EV\r\n" +
	"-000008888?27GL\xc4UBIGER-ID:?88DE88ZZZ8888888888?29Ref. IL888888G2\r\n" +
	"145077/1543?32SOME TEST-VEREIN E.V.\r\n" +
	":62M:C170203EUR5378,36\r\n-"

func Test_SwiftLexer(t *testing.T) {
	lexer := NewSwiftLexer("testlexer", githubIssue13TestData)

	var tokens []Token

	for lexer.HasNext() {
		tokens = append(tokens, lexer.Next())
	}

	expectedLen := 22
	if len(tokens) != expectedLen {
		t.Logf("Expected %d tokens, got %d", expectedLen, len(tokens))
		t.Fail()
	}

	var messageSeparatorCount int
	for _, tk := range tokens {
		if tk.Type() == SWIFT_MESSAGE_SEPARATOR {
			messageSeparatorCount++
		}
	}

	if messageSeparatorCount != 1 {
		t.Logf("Expected one message separator, got %d", messageSeparatorCount)
		t.Fail()
	}
}
