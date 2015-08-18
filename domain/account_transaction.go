package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

type AccountTransaction struct {
	Account              AccountConnection
	Amount               Amount
	ValutaDate           time.Time
	BookingDate          time.Time
	BankID               string
	AccountID            string
	Purpose              string
	Purpose2             string
	AccountBalanceBefore Balance
	AccountBalanceAfter  Balance
}

func (a AccountTransaction) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BookingDate\tAmount\tBankID\tAccountID\tPurpose")
	buf.WriteString("\n")
	buf.WriteString(a.ValutaDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
	buf.WriteString("\t")
	buf.WriteString(a.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.AccountID)
	buf.WriteString("\t")
	buf.WriteString(a.Purpose)
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 0, 4, 0, '\t', 0)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
