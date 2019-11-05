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
	"github.com/mitch000001/go-hbci/transport"
)

// Config defines the basic configuration needed for a Client to work.
type Config struct {
	BankID      string `json:"bank_id"`
	AccountID   string `json:"account_id"`
	PIN         string `json:"pin"`
	URL         string `json:"url"`
	HBCIVersion int    `json:"hbci_version"`
	ProductName string `json:"product_name"`
	Transport   transport.Transport
}

func (c Config) hbciVersion() (segment.HBCIVersion, error) {
	version, ok := segment.SupportedHBCIVersions[c.HBCIVersion]
	if !ok {
		return version, fmt.Errorf("Unsupported HBCI version. Supported versions are %v", domain.SupportedHBCIVersions)
	}
	return version, nil
}

// New creates a new HBCI client. It returns an error if the provided
// HBCI-Version of the config is not supported or if there is no entry in the
// bank institute database for the provided BankID.
//
// If the provided Config does not provide a URL or a HBCI-Version it will be
// looked up in the bankinfo database.
func New(config Config) (*Client, error) {
	bankID := domain.BankID{
		CountryCode: 280,
		ID:          config.BankID,
	}
	bankInfo := bankinfo.FindByBankID(config.BankID)
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
			return nil, fmt.Errorf("Unsupported HBCI version. Supported versions are %v", domain.SupportedHBCIVersions)
		}
		hbciVersion = version
	}
	dcfg := dialog.Config{
		BankID:      bankID,
		HBCIURL:     url,
		UserID:      config.AccountID,
		HBCIVersion: hbciVersion,
		ProductName: config.ProductName,
		Transport:   config.Transport,
	}

	d := dialog.NewPinTanDialog(dcfg)
	d.SetPin(config.PIN)
	client := &Client{
		config:       config,
		hbciVersion:  hbciVersion,
		pinTanDialog: d,
	}
	return client, nil
}

// Client is the main entrypoint to perform high level HBCI requests.
//
// Its methods reflect possible actions and abstract the lower level dialog
// methods.
type Client struct {
	config       Config
	hbciVersion  segment.HBCIVersion
	pinTanDialog *dialog.PinTanDialog
}

func (c *Client) init() error {
	if c.pinTanDialog.BankParameterDataVersion() == 0 {
		_, err := c.pinTanDialog.SyncClientSystemID()
		if err != nil {
			return fmt.Errorf("Error while fetching accounts: %v", err)
		}
	}
	return nil
}

// Accounts return the basic account information for the provided client config.
func (c *Client) Accounts() ([]domain.AccountInformation, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	return c.pinTanDialog.Accounts, nil
}

// AccountTransactions return all transactions for the provided timeframe.
// If allAccouts is true, it will fetch all transactions associated with the
// proviced account. For the initial request no continuationReference is
// needed, as this method will be called recursivly if the server sends one.
func (c *Client) AccountTransactions(account domain.AccountConnection, timeframe domain.Timeframe, allAccounts bool, continuationReference string) ([]domain.AccountTransaction, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
	accountTransactionRequest, err := builder.AccountTransactionRequest(account, allAccounts)
	if err != nil {
		return nil, err
	}
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

// SepaAccountTransactions return all transactions for the provided timeframe.
// If allAccouts is true, it will fetch all transactions associated with the
// provided account. For the initial request no continuationReference is
// needed, as this method will be called recursivly if the server sends one.
func (c *Client) SepaAccountTransactions(account domain.InternationalAccountConnection, timeframe domain.Timeframe, allAccounts bool, continuationReference string) ([]domain.AccountTransaction, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
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

// AccountInformation will print all information attached to the provided
// account. If allAccounts is true it will fetch also the information
// associated with the account.
func (c *Client) AccountInformation(account domain.AccountConnection, allAccounts bool) error {
	if err := c.init(); err != nil {
		return err
	}
	accountInformationRequest := segment.NewAccountInformationRequestSegmentV1(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(message.NewHBCIMessage(c.hbciVersion, accountInformationRequest))
	if err != nil {
		return err
	}
	accountInfoResponse := decryptedMessage.FindMarshaledSegment("HIKIF")
	if accountInfoResponse != nil {
		fmt.Printf("Account Info: %s\n", accountInfoResponse)
		return nil
	}
	return fmt.Errorf("Malformed response: expected HIKIF segment")
}

// AccountBalances retrieves the balance for the provided account.
// If allAccounts is true it will fetch also the balances for all accounts
// associated with the account.
func (c *Client) AccountBalances(account domain.AccountConnection, allAccounts bool) ([]domain.AccountBalance, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
	accountBalanceRequest, err := builder.AccountBalanceRequest(account, allAccounts)
	if err != nil {
		return nil, err
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(
		message.NewHBCIMessage(c.hbciVersion, accountBalanceRequest),
	)
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

// Status returns information about open jobs to fetch from the institute.
// If a continuationReference is present, the status information attached to it
// will be fetched.
func (c *Client) Status(from, to time.Time, maxEntries int, continuationReference string) ([]domain.StatusAcknowledgement, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
	statusRequest, err := builder.StatusProtocolRequest(from, to, maxEntries, continuationReference)
	if err != nil {
		return nil, err
	}
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

// AnonymousClient wraps a Client and allows anonymous requests to bank
// institutes. Examples for those jobs are stock exchange news.
type AnonymousClient struct {
	*Client
}

// CommunicationAccess returns data used to make calls to a given institute.
// Not yet properly implemented, therefore only the raw data are returned.
func (a *AnonymousClient) CommunicationAccess(from, to domain.BankID, maxEntries int) ([]byte, error) {
	commRequest := segment.NewCommunicationAccessRequestSegment(from, to, maxEntries, "")
	decryptedMessage, err := a.pinTanDialog.SendAnonymousMessage(message.NewHBCIMessage(a.hbciVersion, commRequest))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%+#v", decryptedMessage)), nil
}
