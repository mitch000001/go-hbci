package dialog_test

import "github.com/mitch000001/go-hbci/domain"
import "github.com/mitch000001/go-hbci/dialog"

func ExamplePinTanDialog() {
	url := "https://bank.de/hbci"
	userId := "100000000"
	blz := "1000000"
	bankId := domain.BankId{
		CountryCode: 280,
		ID:          blz,
	}
	d := dialog.NewPinTanDialog(bankId, url, userId)
	d.SetPin("12345")
}
