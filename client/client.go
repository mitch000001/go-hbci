package client

import (
	"fmt"
	"strings"

	"github.com/mitch000001/go-hbci/dialog"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

type Config struct {
	BankID      string `json:"bank_id"`
	AccountID   string `json:"account_id"`
	PIN         string `json:"pin"`
	URL         string `json:"url"`
	HBCIVersion int    `json:"hbci_version"`
}

func (c Config) hbciVersion() (segment.HBCIVersion, error) {
	version, ok := segment.SupportedHBCIVersions[c.HBCIVersion]
	if !ok {
		return version, fmt.Errorf("Unsupported HBCI version. Supported versions are %s", domain.SupportedHBCIVersions)
	}
	return version, nil
}

func New(config Config) (*Client, error) {
	bankId := domain.BankId{
		CountryCode: 280,
		ID:          config.BankID,
	}
	var d *dialog.PinTanDialog
	hbciVersion, err := config.hbciVersion()
	if err != nil {
		return nil, err
	}
	d = dialog.NewPinTanDialog(bankId, config.URL, config.AccountID, hbciVersion)
	d.SetPin(config.PIN)
	client := &Client{
		config:       config,
		hbciVersion:  hbciVersion,
		pinTanDialog: d,
	}
	return client, nil
}

type Client struct {
	config       Config
	hbciVersion  segment.HBCIVersion
	pinTanDialog *dialog.PinTanDialog
}

func (c *Client) Accounts() ([]domain.AccountInformation, error) {
	if c.pinTanDialog.UserParameterDataVersion() == 0 {
		_, err := c.pinTanDialog.SyncClientSystemID()
		if err != nil {
			return nil, fmt.Errorf("Error while fetching accounts: %v", err)
		}
	}
	return c.pinTanDialog.Accounts, nil
}

func (c *Client) AccountTransactions(account domain.AccountConnection, timeframe domain.Timeframe, allAccounts bool, aufsetzpunkt string) ([]domain.AccountTransaction, error) {
	accountTransactionRequest := c.hbciVersion.AccountTransactionRequest(account, allAccounts)
	accountTransactionRequest.SetTransactionRange(timeframe)
	if aufsetzpunkt != "" {
		accountTransactionRequest.SetAufsetzpunkt(aufsetzpunkt)
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountTransactionRequest))
	if err != nil {
		return nil, err
	}
	var accountTransactions []domain.AccountTransaction
	accountTransactionResponses := decryptedMessage.FindSegments("HIKAZ")
	if accountTransactionResponses != nil {
		type response struct {
			transactions []domain.AccountTransaction
			err          error
		}
		resFn := func(tr []domain.AccountTransaction, err error) response {
			return response{tr, err}
		}
		responses := make(chan response, len(accountTransactionResponses))
		for _, marshaledSegment := range accountTransactionResponses {
			segment := &segment.AccountTransactionResponseSegment{}
			err = segment.UnmarshalHBCI(marshaledSegment)
			if err != nil {
				return nil, err
			}
			accountTransactions = append(accountTransactions, segment.Transactions()...)
			if segment != nil {
				go func() {
					responses <- resFn(c.AccountTransactions(account, timeframe, allAccounts, aufsetzpunkt))
				}()
			} else {
				responses <- resFn([]domain.AccountTransaction{}, nil)
			}
		}
		var errs []string
		for {
			if len(responses) == 0 {
				break
			}
			res := <-responses
			accountTransactions = append(accountTransactions, res.transactions...)
			if res.err != nil {
				errs = append(errs, fmt.Sprintf("%T:%v", res.err, res.err))
			}
		}
		if len(errs) != 0 {
			return nil, fmt.Errorf("Got errors: %s", strings.Join(errs, "\t"))
		}
	} else {
		return nil, fmt.Errorf("Malformed response: expected HIKAZ segment")
	}

	return accountTransactions, nil
}

func (c *Client) AccountInformation(account domain.AccountConnection, allAccounts bool) error {
	accountInformationRequest := segment.NewAccountInformationRequestSegmentV1(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountInformationRequest))
	if err != nil {
		return err
	}
	accountInfoResponse := decryptedMessage.FindSegment("HIKIF")
	if accountInfoResponse != nil {
		fmt.Printf("Account Info: %s\n", accountInfoResponse)
		return nil
	} else {
		return fmt.Errorf("Malformed response: expected HIKIF segment")
	}
	return nil
}

func (c *Client) AccountBalances(account domain.AccountConnection, allAccounts bool) ([]domain.AccountBalance, error) {
	accountBalanceRequest := c.hbciVersion.AccountBalanceRequest(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountBalanceRequest))
	if err != nil {
		return nil, err
	}
	var balances []domain.AccountBalance
	balanceResponses := decryptedMessage.FindSegments("HISAL")
	if balanceResponses != nil {
		for _, marshaledSegment := range balanceResponses {
			balanceSegment := &segment.AccountBalanceResponseSegment{}
			err = balanceSegment.UnmarshalHBCI(marshaledSegment)
			if err != nil {
				return nil, fmt.Errorf("Error while parsing account balance: %v", err)
			}
			balances = append(balances, balanceSegment.AccountBalance())
		}
	} else {
		return nil, fmt.Errorf("Malformed response: expected HISAL segment")
	}

	return balances, nil
}
