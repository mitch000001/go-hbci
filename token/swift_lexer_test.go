package token

import (
	"reflect"
	"testing"
)

const githubIssue13TestData = "\r\n" +
	":20:MT940-1803060458\r\n" +
	":21:NONREF\r\n" +
	":25:20040000/12345678EUR\r\n" +
	":28C:0/13\r\n" +
	":60M:C170201EUR1234,56\r\n" +
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
	":62M:C170203EUR5378,36\r\n-" +
	"\r\n" +
	":20:MT940-1803060458\r\n" +
	":21:NONREF\r\n" +
	":25:20040000/12345678EUR\r\n" +
	":28C:0/13\r\n" +
	":60M:C170201EUR1234,56\r\n" +
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
	":62M:C170203EUR5378,36\r\n-" +
	"\r\n" +
	":20:MT940-1804300355\r\n" +
	":21:NONREF\r\n" +
	":25:20012345/112233445EUR\r\n" +
	":28C:0/2\r\n" +
	":60M:C170203EUR1234,56\r\n" +
	":61:1702030203DR1,00NMSCNONREF//POS 3416383187\r\n" +
	":86:005?20LASTSCHRIFT/BELAST.?21SEPADDDD00009999999-01 DRI?22VENOW 9\r\n" +
	"3?2350790427 IHRE BILLPAYZAHLUN?24G 0055 3?555 555 555 555 WWW.BI\r\n" +
	"LLPAY.D?26E?27END-TO-END-REF.:?28SEPADDDD00009999999-01?29CORE /\r\n" +
	" MANDATSREF.:?32BILLPAY GMBH?12345556667-1-12345678?61GLÄUBIGER-ID\r\n" +
	":?62DE19ZZZ00000999999?63Ref. I9999099B999999/280\r\n" +
	":61:9909999203DR1000,00NMSCNONREF//POS 9999999999\r\n" +
	":86:820?20ÜBERTRAG/ÜBERWEISUNG?21CBAEURXABCDEFG?22END-TO-END-REF.:?23\r\n" +
	"NICHT ANGEGEBEN?24Ref. H09999999999999/2?30ABCDEF22XXX?31EE99999\r\n" +
	"0999009999999?32COINBASE UK, LTD.\r\n" +
	":62M:C170206EUR1234,56\r\n-"

func Test_SwiftLexer(t *testing.T) {
	t.Run("github issue 13", func(t *testing.T) {
		lexer := NewSwiftLexer("testlexer", []byte(githubIssue13TestData))

		var tokens []Token

		for lexer.HasNext() {
			tokens = append(tokens, lexer.Next())
		}

		var messageSeparatorCount int
		var tags []string
		for _, tk := range tokens {
			if tk.Type() == SWIFT_MESSAGE_SEPARATOR {
				messageSeparatorCount++
			}
			if tk.Type() == SWIFT_TAG {
				tags = append(tags, string(tk.Value()))
			}
		}

		if messageSeparatorCount != 3 {
			t.Logf("Expected %d message separator, got %d", 3, messageSeparatorCount)
			t.Fail()
		}

		expectedTagCount := 30
		if len(tags) != expectedTagCount {
			t.Logf("Expected %d tags, got %d", expectedTagCount, len(tags))
			t.Fail()
		}
		expectedTags := []string{
			":20:",
			":21:",
			":25:",
			":28C:",
			":60M:",
			":61:",
			":86:",
			":61:",
			":86:",
			":62M:",
			":20:",
			":21:",
			":25:",
			":28C:",
			":60M:",
			":61:",
			":86:",
			":61:",
			":86:",
			":62M:",
			":20:",
			":21:",
			":25:",
			":28C:",
			":60M:",
			":61:",
			":86:",
			":61:",
			":86:",
			":62M:",
		}

		if !reflect.DeepEqual(tags, expectedTags) {
			t.Logf("Expected tags to equal\n%v\n\tgot\n%v\n", expectedTags, tags)
			t.Fail()
		}
	})
}
