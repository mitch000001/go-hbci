package client

import (
	"fmt"
	"time"

	"github.com/mitch000001/go-hbci/bankinfo"
	"github.com/mitch000001/go-hbci/dialog"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/internal"
	"github.com/mitch000001/go-hbci/message"
	"github.com/mitch000001/go-hbci/segment"
	"github.com/mitch000001/go-hbci/swift"
	"github.com/mitch000001/go-hbci/transport"
)

// Config defines the basic configuration needed for a Client to work.
type Config struct {
	BankID             string `json:"bank_id"`
	AccountID          string `json:"account_id"`
	PIN                string `json:"pin"`
	URL                string `json:"url"`
	HBCIVersion        int    `json:"hbci_version"`
	Transport          transport.Transport
	ProductName        string `json:"product_name"`
	ProductVersion     string `json:"product_version"`
	EnableDebugLogging bool   `json:"enable_debug_logging"`
}

func (c Config) hbciVersion() (segment.HBCIVersion, error) {
	version, ok := segment.SupportedHBCIVersions[c.HBCIVersion]
	if !ok {
		return version, fmt.Errorf("unsupported HBCI version. Supported versions are %v", domain.SupportedHBCIVersions)
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
	internal.SetDebugMode(config.EnableDebugLogging)
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
			return nil, fmt.Errorf("unsupported HBCI version. Supported versions are %v", domain.SupportedHBCIVersions)
		}
		hbciVersion = version
	}
	dcfg := dialog.Config{
		BankID:         bankID,
		HBCIURL:        url,
		UserID:         config.AccountID,
		HBCIVersion:    hbciVersion,
		ProductName:    config.ProductName,
		ProductVersion: config.ProductVersion,
		Transport:      config.Transport,
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
			return fmt.Errorf("error while fetching accounts: %v", err)
		}
	}
	return nil
}

// Accounts return the basic account information for the provided client config.
func (c *Client) Accounts() ([]domain.AccountInformation, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	err := c.pinTanDialog.SyncUserParameterData()
	if err != nil {
		return nil, fmt.Errorf("error getting accounts")
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
	requestBuilder := func() (segment.AccountTransactionRequest, error) {
		builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
		return builder.AccountTransactionRequest(account, allAccounts)
	}
	bookedSwiftTransactions, err := c.accountTransactions(requestBuilder, timeframe, continuationReference)
	if err != nil {
		return nil, fmt.Errorf("error executing HBCI request: %w", err)
	}
	unmarshaler := swift.NewMT940MessagesUnmarshaler()
	tx, err := unmarshaler.UnmarshalMT940(bookedSwiftTransactions.Data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling SWIFT transactions: %w", err)
	}
	return tx, nil
}

// SepaAccountTransactions return all transactions for the provided timeframe.
// If allAccouts is true, it will fetch all transactions associated with the
// provided account. For the initial request no continuationReference is
// needed, as this method will be called recursivly if the server sends one.
func (c *Client) SepaAccountTransactions(account domain.InternationalAccountConnection, timeframe domain.Timeframe, allAccounts bool, continuationReference string) ([]domain.AccountTransaction, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	requestBuilder := func() (segment.AccountTransactionRequest, error) {
		builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
		return builder.SepaAccountTransactionRequest(account, allAccounts)
	}
	bookedSwiftTransactions, err := c.accountTransactions(requestBuilder, timeframe, continuationReference)
	if err != nil {
		return nil, fmt.Errorf("error executing HBCI request: %w", err)
	}
	unmarshaler := swift.NewMT940MessagesUnmarshaler()
	tx, err := unmarshaler.UnmarshalMT940(bookedSwiftTransactions.Data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling SWIFT transactions: %w", err)
	}
	return tx, nil
}

func (c *Client) accountTransactions(requestBuilder func() (segment.AccountTransactionRequest, error), timeframe domain.Timeframe, continuationReference string) (*swift.MT940Messages, error) {
	accountTransactionRequest, err := requestBuilder()
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}
	accountTransactionRequest.SetTransactionRange(timeframe)
	if continuationReference != "" {
		accountTransactionRequest.SetContinuationReference(continuationReference)
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(
		message.NewHBCIMessage(c.hbciVersion, c.hbciVersion.TanProcess4Request(segment.IdentificationID), accountTransactionRequest),
	)
	if err != nil {
		return nil, fmt.Errorf("error sending hbci request: %w", err)
	}
	var bookedSwiftTransactions []*swift.MT940Messages
	accountTransactionResponses := decryptedMessage.FindSegments("HIKAZ")
	for _, unmarshaledSegment := range accountTransactionResponses {
		seg, ok := unmarshaledSegment.(segment.AccountTransactionResponse)
		if !ok {
			return nil, fmt.Errorf("malformed segment found with ID `HIKAZ`")
		}
		bookedSwiftTransactions = append(bookedSwiftTransactions, seg.BookedSwiftTransactions())
	}
	var newContinuationReference string
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.Code == element.AcknowledgementAdditionalInformation {
			newContinuationReference = ack.Params[0]
			break
		}
	}
	tx := swift.MergeMT940Messages(bookedSwiftTransactions...)
	if newContinuationReference == "" {
		return tx, nil
	}
	msg, err := c.accountTransactions(requestBuilder, timeframe, newContinuationReference)
	if err != nil {
		return nil, err
	}
	return swift.MergeMT940Messages(msg, tx), err
}

// AccountInformation will print all information attached to the provided
// account. If allAccounts is true it will fetch also the information
// associated with the account.
func (c *Client) AccountInformation(account domain.AccountConnection, allAccounts bool) error {
	if err := c.init(); err != nil {
		return err
	}
	accountInformationRequest := segment.NewAccountInformationRequestSegmentV1(account, allAccounts)
	decryptedMessage, err := c.pinTanDialog.SendMessage(
		message.NewHBCIMessage(c.hbciVersion, c.hbciVersion.TanProcess4Request(segment.IdentificationID), accountInformationRequest),
	)
	if err != nil {
		return err
	}
	accountInfoResponse := decryptedMessage.FindMarshaledSegment("HIKIF")
	if accountInfoResponse != nil {
		fmt.Printf("Account Info: %s\n", accountInfoResponse)
		return nil
	}
	return fmt.Errorf("malformed response: expected HIKIF segment")
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
		message.NewHBCIMessage(
			c.hbciVersion,
			c.hbciVersion.TanProcess4Request(segment.IdentificationID),
			accountBalanceRequest,
		),
	)
	if err != nil {
		return nil, err
	}
	var balances []domain.AccountBalance
	balanceResponses := decryptedMessage.FindSegments(segment.AccountBalanceResponseID)
	for _, unmarshaledSegment := range balanceResponses {
		seg, ok := unmarshaledSegment.(segment.AccountBalanceResponse)
		if !ok {
			return nil, fmt.Errorf("malformed segment found with ID %q", segment.AccountBalanceResponseID)
		}
		balances = append(balances, seg.AccountBalance())
	}
	if len(balanceResponses) == 0 {
		return nil, fmt.Errorf("malformed response: expected HISAL segment")
	}

	return balances, nil
}

// AccountBalances retrieves the balance for the provided account.
// If allAccounts is true it will fetch also the balances for all accounts
// associated with the account.
func (c *Client) SepaAccountBalances(account domain.InternationalAccountConnection, allAccounts bool, continuationReference string) ([]domain.SepaAccountBalance, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	builder := segment.NewBuilder(c.pinTanDialog.SupportedSegments())
	accountBalanceRequest, err := builder.SepaAccountBalanceRequest(account, allAccounts)
	if err != nil {
		return nil, err
	}
	if continuationReference != "" {
		accountBalanceRequest.SetContinuationMark(continuationReference)
	}
	decryptedMessage, err := c.pinTanDialog.SendMessage(
		message.NewHBCIMessage(
			c.hbciVersion,
			c.hbciVersion.TanProcess4Request(segment.IdentificationID),
			accountBalanceRequest,
		),
	)
	if err != nil {
		return nil, err
	}
	var balances []domain.SepaAccountBalance
	balanceResponses := decryptedMessage.FindSegments(segment.AccountBalanceResponseID)
	for _, unmarshaledSegment := range balanceResponses {
		seg, ok := unmarshaledSegment.(segment.AccountBalanceResponse)
		if !ok {
			return nil, fmt.Errorf("malformed segment found with ID %q", segment.AccountBalanceResponseID)
		}
		sepaBalances, err := seg.SepaAccountBalance()
		if err != nil {
			return nil, fmt.Errorf("could not get sepa balances: %w", err)
		}
		balances = append(balances, sepaBalances)
	}
	if len(balanceResponses) == 0 {
		return nil, fmt.Errorf("malformed response: expected HISAL segment")
	}
	var newContinuationReference string
	acknowledgements := decryptedMessage.Acknowledgements()
	for _, ack := range acknowledgements {
		if ack.Code == element.AcknowledgementAdditionalInformation {
			newContinuationReference = ack.Params[0]
			break
		}
	}
	if newContinuationReference == "" {
		return balances, nil
	}
	nextBal, err := c.SepaAccountBalances(account, allAccounts, newContinuationReference)
	if err != nil {
		return nil, err
	}
	balances = append(balances, nextBal...)
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
	bankMessage, err := c.pinTanDialog.SendMessage(
		message.NewHBCIMessage(c.hbciVersion, c.hbciVersion.TanProcess4Request(segment.IdentificationID), statusRequest),
	)
	if err != nil {
		return nil, err
	}
	var statusAcknowledgements []domain.StatusAcknowledgement
	statusResponses := bankMessage.FindSegments("HIPRO")
	for _, seg := range statusResponses {
		statusResponse := seg.(segment.StatusProtocolResponse)
		statusAcknowledgements = append(statusAcknowledgements, statusResponse.Status())
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

// BankParameterData return the bankparameter data from the institute.
func (c *AnonymousClient) BankParameterData() (*dialog.BankParameterData, error) {
	bpd, err := c.pinTanDialog.GetAnonymousBankParameterData()
	if err != nil {
		return nil, fmt.Errorf("error getting anpnymous bank parameter data: %w", err)
	}
	return bpd, nil
}
