package domain

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

// AccountTransactions represents a printable version of a slice of
// account transactions
type AccountTransactions struct {
	BookedTransactions   []BookedAccountTransactions
	UnbookedTransactions []UnbookedAccountTransactions
}

func (at AccountTransactions) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	for _, a := range at.BookedTransactions {
		buf.WriteString(a.String())
	}
	for _, a := range at.UnbookedTransactions {
		buf.WriteString(a.String())
	}
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 24, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

// AccountTransaction represents one transaction entry for a given account
type AccountTransaction struct {
	Account              AccountConnection
	Amount               Amount
	ValutaDate           time.Time
	BookingDate          time.Time
	BookingText          string
	BankID               string
	AccountID            string
	Name                 string
	Purpose              string
	Purpose2             string
	TransactionID        int
	AccountBalanceBefore Balance
	AccountBalanceAfter  Balance
}

func (a AccountTransaction) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BookingDate\tAmount\tBankID\tAccountID\tPurpose")
	buf.WriteString("\n")
	buf.WriteString(a.BookingDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
	buf.WriteString("\t")
	buf.WriteString(a.BankID)
	buf.WriteString("\t")
	buf.WriteString(a.AccountID)
	buf.WriteString("\t")
	buf.WriteString(a.Purpose)
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 0, 8, 0, '\t', 0)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

// BookedAccountTransactions represents a printable version of
// account transactions
type BookedAccountTransactions struct {
	Account              AccountConnection
	AccountBalanceBefore Balance
	AccountBalanceAfter  Balance
	Transactions         []Transaction
}

func (at BookedAccountTransactions) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	fmt.Fprintf(
		&buf, "Booked transactions for account %s/%s\n",
		at.Account.BankID, at.Account.AccountID,
	)
	fmt.Fprintf(
		&buf, "Balance at %s: %.2f %s\n",
		at.AccountBalanceBefore.TransmissionDate.Format("2006-01-02"),
		at.AccountBalanceBefore.Amount.Amount,
		at.AccountBalanceBefore.Amount.Currency,
	)
	buf.WriteString("BookingDate\tBooking Text\tAmount\tBankID\tAccountID\tName\tPurpose")
	buf.WriteString("\n")
	for _, a := range at.Transactions {
		buf.WriteString(a.BookingDate.Format("2006-01-02"))
		buf.WriteString("\t")
		buf.WriteString(a.BookingText)
		buf.WriteString("\t")
		buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
		buf.WriteString("\t")
		buf.WriteString(a.BankID)
		buf.WriteString("\t")
		buf.WriteString(a.AccountID)
		buf.WriteString("\t")
		buf.WriteString(a.Name)
		buf.WriteString("\t")
		buf.WriteString(a.Purpose)
		buf.WriteString("\n")
	}
	fmt.Fprintf(
		&buf, "Balance at %s: %.2f %s\n",
		at.AccountBalanceAfter.TransmissionDate.Format("2006-01-02"),
		at.AccountBalanceAfter.Amount.Amount,
		at.AccountBalanceAfter.Amount.Currency,
	)
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 24, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

type UnbookedAccountTransactions struct {
	Account            AccountConnection
	CreationDate       time.Time
	DebitAmount        Amount
	CreditAmount       Amount
	DebitTransactions  int
	CreditTransactions int
	Transactions       []Transaction
}

func (u UnbookedAccountTransactions) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	fmt.Fprintf(
		&buf, "Unbooked transactions for account %s/%s\n",
		u.Account.BankID, u.Account.AccountID,
	)
	fmt.Fprintf(
		&buf, "Created at %s\n",
		u.CreationDate.Format("2006-01-02T"),
	)
	fmt.Fprintf(
		&buf, "Debit balance %.2f %s (%d transactions)\n",
		u.DebitAmount.Amount,
		u.DebitAmount.Currency,
		u.DebitTransactions,
	)
	fmt.Fprintf(
		&buf, "Credit balance %.2f %s (%d transactions)\n",
		u.CreditAmount.Amount,
		u.CreditAmount.Currency,
		u.CreditTransactions,
	)
	buf.WriteString("BookingDate\tBooking Text\tAmount\tBankID\tAccountID\tName\tPurpose")
	buf.WriteString("\n")
	for _, a := range u.Transactions {
		buf.WriteString(a.BookingDate.Format("2006-01-02"))
		buf.WriteString("\t")
		buf.WriteString(a.BookingText)
		buf.WriteString("\t")
		buf.WriteString(fmt.Sprintf("%.2f %s", a.Amount.Amount, a.Amount.Currency))
		buf.WriteString("\t")
		buf.WriteString(a.BankID)
		buf.WriteString("\t")
		buf.WriteString(a.AccountID)
		buf.WriteString("\t")
		buf.WriteString(a.Name)
		buf.WriteString("\t")
		buf.WriteString(a.Purpose)
		buf.WriteString("\n")
	}
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 24, 1, 0, ' ', tabwriter.TabIndent)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}

type Transaction struct {
	Amount        Amount
	ValutaDate    time.Time
	BookingDate   time.Time
	CreationDate  time.Time
	BookingText   string
	BankID        string
	AccountID     string
	Name          string
	Purpose       string
	Purpose2      string
	TransactionID int
}

func (t Transaction) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	buf.WriteString("BookingDate\tAmount\tBankID\tAccountID\tPurpose")
	buf.WriteString("\n")
	buf.WriteString(t.BookingDate.Format("2006-01-02"))
	buf.WriteString("\t")
	buf.WriteString(fmt.Sprintf("%.2f %s", t.Amount.Amount, t.Amount.Currency))
	buf.WriteString("\t")
	buf.WriteString(t.BankID)
	buf.WriteString("\t")
	buf.WriteString(t.AccountID)
	buf.WriteString("\t")
	buf.WriteString(t.Purpose)
	var out bytes.Buffer
	tabw := tabwriter.NewWriter(&out, 0, 8, 0, '\t', 0)
	fmt.Fprint(tabw, buf.String())
	tabw.Flush()
	return out.String()
}
