package client

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/bankinfo"
	"github.com/mitch000001/go-hbci/dialog"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
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
	bankInfo := bankinfo.FindByBankId(config.BankID)
	var (
		url         string
		hbciVersion segment.HBCIVersion
	)
	if config.URL != "" {
		url = config.URL
	} else {
		url = bankInfo.URL
	}
	if config.HBCIVersion > 0 {
		version, err := config.hbciVersion()
		if err != nil {
			return nil, err
		}
		hbciVersion = version
	} else {
		version, ok := segment.SupportedHBCIVersions[bankInfo.HbciVersion()]
		if !ok {
			return nil, fmt.Errorf("Unsupported HBCI version. Supported versions are %s", domain.SupportedHBCIVersions)
		}
		hbciVersion = version
	}
	d := dialog.NewPinTanDialog(bankId, url, config.AccountID, hbciVersion)
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

func (c *Client) AccountTransactions(account domain.AccountConnection, timeframe domain.Timeframe, allAccounts bool, continuationReference string) ([]domain.AccountTransaction, error) {
	accountTransactionRequest := c.hbciVersion.AccountTransactionRequest(account, allAccounts)
	accountTransactionRequest.SetTransactionRange(timeframe)
	if continuationReference != "" {
		accountTransactionRequest.SetContinuationReference(continuationReference)
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountTransactionRequest))
	if err != nil {
		return nil, err
	}
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.Code == element.AcknowledgementAdditionalInformation {
			fmt.Printf("Additional information: %+v\n", ack)
		}
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
		for _, unmarshaledSegment := range accountTransactionResponses {
			seg := unmarshaledSegment.(segment.AccountTransactionResponse)
			accountTransactions = append(accountTransactions, seg.Transactions()...)
			if seg != nil {
				go func() {
					responses <- resFn(c.AccountTransactions(account, timeframe, allAccounts, continuationReference))
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

func (c *Client) SepaAccountTransactions(account domain.InternationalAccountConnection, timeframe domain.Timeframe, allAccounts bool, continuationReference string) ([]domain.AccountTransaction, error) {
	accountTransactionRequest := c.hbciVersion.SepaAccountTransactionRequest(account, allAccounts)
	accountTransactionRequest.SetTransactionRange(timeframe)
	if continuationReference != "" {
		accountTransactionRequest.SetContinuationReference(continuationReference)
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountTransactionRequest))
	if err != nil {
		return nil, err
	}
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.Code == element.AcknowledgementAdditionalInformation {
			fmt.Printf("Additional information: %+v\n", ack)
		}
	}
	var accountTransactions []domain.AccountTransaction
	accountTransactionResponses := decryptedMessage.FindSegments("HIKAZ")
	if accountTransactionResponses != nil {
		for _, unmarshaledSegment := range accountTransactionResponses {
			seg := unmarshaledSegment.(segment.AccountTransactionResponse)
			accountTransactions = append(accountTransactions, seg.Transactions()...)
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
	accountInfoResponse := decryptedMessage.FindMarshaledSegment("HIKIF")
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
	balanceResponses := decryptedMessage.FindMarshaledSegments("HISAL")
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

func (c *Client) Status(from, to time.Time, maxEntries int, continuationReference string) ([]domain.StatusAcknowledgement, error) {
	statusRequest := c.hbciVersion.StatusProtocolRequest(from, to, maxEntries, continuationReference)
	bankMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, statusRequest))
	if err != nil {
		return nil, err
	}
	var statusAcknowledgements []domain.StatusAcknowledgement
	statusResponses := bankMessage.FindSegments("HIPRO")
	if statusResponses != nil {
		for _, seg := range statusResponses {
			statusResponse := seg.(segment.StatusProtocolResponse)
			statusAcknowledgements = append(statusAcknowledgements, statusResponse.Status())
		}
	}
	return statusAcknowledgements, nil
}

type AnonymousClient struct {
	*Client
}

func (a *AnonymousClient) CommunicationAccess(from, to domain.BankId, maxEntries int) ([]byte, error) {
	commRequest := segment.NewCommunicationAccessRequestSegment(from, to, maxEntries, "")
	decryptedMessage, err := a.pinTanDialog.SendAnonymousMessage(message.NewHBCIMessage(a.hbciVersion, commRequest))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%+#v", decryptedMessage)), nil
}
