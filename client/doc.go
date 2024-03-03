/*
Package client provides a high level API for HBCI-Requests

The main types of this package are the Config and the Client itself. The

Config provides general information about the account to use. It should be
sufficient to provide a config with BankID (i.e. 'Bankleitzahl, BLZ),
AccountID (i.e. Kontonummer) and the PIN.  The fields URL and HBCIVersion
are optional fields for users with deeper knowledge about the bank institute
and its HBCI endpoints. If one of these is not provided it will be looked up
from the bankinfo package.

Client provides a convenient way of issuing certain requests to the HBCI
server. All low level APIs are queried from the Client and it returns only
types from the domain package.

# Example to get account transactions

	c := newClient()

	testAccount := domain.AccountConnection{
		AccountID: "1000000000", CountryCode: 280, BankID: "100000000"
	}

	timeframe := domain.Timeframe{
		StartDate: domain.NewShortDate(time.Now().AddDate(0, 0, -90)),
	}

	transactions, err := c.AccountTransactions(testAccount, timeframe, false, "")
	if err != nil {
		panic(err)
	}

	for _, transaction := range transactions {
		fmt.Printf("%+v\n", transaction)
	}

Example to get account balances

	c := newClient()

	testAccount := domain.AccountConnection{
		AccountID: "1000000000", CountryCode: 280, BankID: "100000000"
	}

	balances, err := c.AccountBalances(testAccount, true)
	if err != nil {
		panic(err)
	}

	for _, balance := range balances {
		fmt.Printf("%+v\n", balance)
	}
*/
package client
