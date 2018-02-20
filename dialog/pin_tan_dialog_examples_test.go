package dialog_test

import (
	"github.com/mitch000001/go-hbci/domain"
)
import "github.com/mitch000001/go-hbci/dialog"

func ExamplePinTanDialog() {
	cfg := dialog.Config{
		HBCIURL: "https://bank.de/hbci",
		UserID:  "100000000",
		BankID: domain.BankID{
			CountryCode: 280,
			ID:          "1000000",
		},
	}
	d := dialog.NewPinTanDialog(cfg)
	d.SetPin("12345")
}
