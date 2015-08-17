package client

import (
	"fmt"

	"github.com/mitch000001/go-hbci/dialog"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
)

const (
	Version220 = 220
)

var supportedVersions = []int{
	Version220,
}

type Config struct {
	BankID      string `json:"bank_id"`
	AccountID   string `json:"account_id"`
	PIN         string `json:"pin"`
	URL         string `json:"url"`
	HBCIVersion int    `json:"hbci_version"`
}

func New(config Config) (*Client, error) {
	bankId := domain.BankId{
		CountryCode: 280,
		ID:          config.BankID,
	}
	d := dialog.NewPinTanDialog(bankId, config.URL, config.AccountID)
	d.SetPin(config.PIN)
	client := &Client{
		config:       config,
		pinTanDialog: d,
	}
	return client, nil
}

type Client struct {
	config       Config
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

func (c *Client) AccountTransactions(account domain.AccountConnection, timeframe domain.Timeframe, allAccounts bool) ([]domain.AccountTransaction, error) {
	accountTransactionRequest := segment.NewAccountTransactionRequestSegment(account, allAccounts)
	accountTransactionRequest.SetTransactionRange(timeframe)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(accountTransactionRequest))
	if err != nil {
		return nil, err
	}
	var accountTransactions []domain.AccountTransaction
	accountTransactionResponses := decryptedMessage.FindSegments("HIKAZ")
	if accountTransactionResponses != nil {
		for _, marshaledSegment := range accountTransactionResponses {
			segment := &segment.AccountTransactionResponseSegment{}
			err = segment.UnmarshalHBCI(marshaledSegment)
			if err != nil {
				return nil, err
			}
			accountTransactions = append(accountTransactions, segment.Transactions()...)
		}
	} else {
		return nil, fmt.Errorf("Malformed response: expected HIKAZ segment")
	}

	return accountTransactions, nil
}

func (c *Client) AccountInformation(account domain.AccountConnection, allAccounts bool) error {
	accountInformationRequest := segment.NewAccountInformationRequestSegment(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(accountInformationRequest))
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
	accountBalanceRequest := segment.NewAccountBalanceRequestSegment(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(accountBalanceRequest))
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
