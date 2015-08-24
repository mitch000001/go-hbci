package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

type AccountBalance struct {
	Account          AccountConnection
	ProductName      string
	Currency         string
	BookedBalance    Balance
	EarmarkedBalance *Balance
	CreditLimit      *Amount
	AvailableAmount  *Amount
	UsedAmount       *Amount
	BookingDate      *time.Time
	DueDate          *time.Time
}

func (a AccountBalance) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BankID\tAccountID\t\tBookingDate\tAmount")
	buf.WriteString("\n")
	buf.WriteString(a.Account.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.Account.AccountID)
	buf.WriteString("\t")
	buf.WriteString("\t")
	buf.WriteString(a.BookedBalance.TransmissionDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", a.BookedBalance.Amount.Amount, a.BookedBalance.Amount.Currency))
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 0, 4, 0, '\t', 0)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
