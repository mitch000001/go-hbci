package segment

import (
	"sort"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

var accountBalanceRequests = map[int]func(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest{
	5: NewAccountBalanceRequestV5,
	6: NewAccountBalanceRequestV6,
}

var sepaAccountBalanceRequests = map[int]func(account domain.InternationalAccountConnection, allAccounts bool) AccountBalanceRequest{
	7: NewAccountBalanceRequestV7,
	8: NewAccountBalanceRequestV8,
}

// AccountBalanceRequestBuilder returns the highest matching versioned segment
func AccountBalanceRequestBuilder(versions []int) (func(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		builder, ok := accountBalanceRequests[version]
		if ok {
			return builder, nil
		}
	}
	return nil, &unsupportedSegmentVersionError{segmentID: "HKSAL", versions: versions}
}

// SepaAccountBalanceRequestBuilder returns the highest matching versioned segment
func SepaAccountBalanceRequestBuilder(versions []int) (func(account domain.InternationalAccountConnection, allAccounts bool) AccountBalanceRequest, error) {
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))
	for _, version := range versions {
		builder, ok := sepaAccountBalanceRequests[version]
		if ok {
			return builder, nil
		}
	}
	return nil, fmt.Errorf("unsupported versions %v", versions)
}

type AccountBalanceRequest interface {
	ClientSegment
	SetContinuationMark(continuationMark string)
}

func NewAccountBalanceRequestV5(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV5{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV5 struct {
	ClientSegment
	AccountConnection     *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV5) Version() int         { return 5 }
func (a *AccountBalanceRequestSegmentV5) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV5) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV5) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV5) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV5) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountBalanceRequestV6(account domain.AccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV6{
		AccountConnection: element.NewAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV6 struct {
	ClientSegment
	AccountConnection     *element.AccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV6) Version() int         { return 6 }
func (a *AccountBalanceRequestSegmentV6) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV6) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV6) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV6) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountBalanceRequestV7(account domain.InternationalAccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV7{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV7 struct {
	ClientSegment
	AccountConnection     *element.InternationalAccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV7) Version() int         { return 7 }
func (a *AccountBalanceRequestSegmentV7) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV7) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV7) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV7) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

func NewAccountBalanceRequestV8(account domain.InternationalAccountConnection, allAccounts bool) AccountBalanceRequest {
	a := &AccountBalanceRequestSegmentV8{
		AccountConnection: element.NewInternationalAccountConnection(account),
		AllAccounts:       element.NewBoolean(allAccounts),
	}
	a.ClientSegment = NewBasicSegment(1, a)
	return a
}

type AccountBalanceRequestSegmentV8 struct {
	ClientSegment
	AccountConnection     *element.InternationalAccountConnectionDataElement
	AllAccounts           *element.BooleanDataElement
	MaxEntries            *element.NumberDataElement
	ContinuationReference *element.AlphaNumericDataElement
}

func (a *AccountBalanceRequestSegmentV8) Version() int         { return 8 }
func (a *AccountBalanceRequestSegmentV8) ID() string           { return "HKSAL" }
func (a *AccountBalanceRequestSegmentV8) referencedId() string { return "" }
func (a *AccountBalanceRequestSegmentV8) sender() string       { return senderUser }

func (a *AccountBalanceRequestSegmentV8) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AllAccounts,
		a.MaxEntries,
		a.ContinuationReference,
	}
}

func (a *AccountBalanceRequestSegmentV8) SetContinuationMark(continuationMark string) {
	a.ContinuationReference = element.NewAlphaNumeric(continuationMark, 35)
}

const AccountBalanceResponseID = "HISAL"

//go:generate go run ../cmd/unmarshaler/unmarshaler_generator.go -segment AccountBalanceResponseSegment -segment_interface AccountBalanceResponse -segment_versions="AccountBalanceResponseSegmentV5:5:Segment,AccountBalanceResponseSegmentV6:6:Segment,AccountBalanceResponseSegmentV7:7:Segment,AccountBalanceResponseSegmentV8:8:Segment"

type AccountBalanceResponse interface {
	BankSegment
	AccountBalance() domain.AccountBalance
	SepaAccountBalance() (domain.SepaAccountBalance, error)
}

type AccountBalanceResponseSegment struct {
	AccountBalanceResponse
}

type AccountBalanceResponseSegmentV5 struct {
	Segment
	AccountConnection  *element.AccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarmarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	DueDate            *element.DateDataElement
}

func (a *AccountBalanceResponseSegmentV5) Version() int         { return 5 }
func (a *AccountBalanceResponseSegmentV5) ID() string           { return AccountBalanceResponseID }
func (a *AccountBalanceResponseSegmentV5) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegmentV5) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegmentV5) AccountBalance() domain.AccountBalance {
	balance := domain.AccountBalance{
		Account:       a.AccountConnection.Val(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}
func (a *AccountBalanceResponseSegmentV5) SepaAccountBalance() (domain.SepaAccountBalance, error) {
	return domain.SepaAccountBalance{}, fmt.Errorf("not implemented")
}
func (a *AccountBalanceResponseSegmentV5) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarmarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.DueDate,
	}
}

type AccountBalanceResponseSegmentV6 struct {
	Segment
	AccountConnection  *element.AccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarmarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	DueDate            *element.DateDataElement
}

func (a *AccountBalanceResponseSegmentV6) Version() int         { return 6 }
func (a *AccountBalanceResponseSegmentV6) ID() string           { return AccountBalanceResponseID }
func (a *AccountBalanceResponseSegmentV6) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegmentV6) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegmentV6) AccountBalance() domain.AccountBalance {
	balance := domain.AccountBalance{
		Account:       a.AccountConnection.Val(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}
func (a *AccountBalanceResponseSegmentV6) SepaAccountBalance() (domain.SepaAccountBalance, error) {
	return domain.SepaAccountBalance{}, fmt.Errorf("not implemented")
}
func (a *AccountBalanceResponseSegmentV6) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarmarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.DueDate,
	}
}

type AccountBalanceResponseSegmentV7 struct {
	Segment
	AccountConnection  *element.InternationalAccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarmarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	DueDate            *element.DateDataElement
}

func (a *AccountBalanceResponseSegmentV7) Version() int         { return 7 }
func (a *AccountBalanceResponseSegmentV7) ID() string           { return AccountBalanceResponseID }
func (a *AccountBalanceResponseSegmentV7) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegmentV7) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegmentV7) AccountBalance() domain.AccountBalance {
	balance := domain.AccountBalance{
		Account:       a.AccountConnection.Val().ToAccountConnection(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}

func (a *AccountBalanceResponseSegmentV7) SepaAccountBalance() (domain.SepaAccountBalance, error) {
	balance := domain.SepaAccountBalance{
		Account:       a.AccountConnection.Val(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance, nil
}
func (a *AccountBalanceResponseSegmentV7) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarmarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.DueDate,
	}
}

type AccountBalanceResponseSegmentV8 struct {
	Segment
	AccountConnection  *element.InternationalAccountConnectionDataElement
	AccountProductName *element.AlphaNumericDataElement
	AccountCurrency    *element.CurrencyDataElement
	BookedBalance      *element.BalanceDataElement
	EarmarkedBalance   *element.BalanceDataElement
	CreditLimit        *element.AmountDataElement
	AvailableAmount    *element.AmountDataElement
	UsedAmount         *element.AmountDataElement
	BookingDate        *element.DateDataElement
	BookingTime        *element.TimeDataElement
	DueDate            *element.DateDataElement
	SeizableAfterMonth *element.AmountDataElement
}

func (a *AccountBalanceResponseSegmentV8) Version() int         { return 8 }
func (a *AccountBalanceResponseSegmentV8) ID() string           { return AccountBalanceResponseID }
func (a *AccountBalanceResponseSegmentV8) referencedId() string { return "HKSAL" }
func (a *AccountBalanceResponseSegmentV8) sender() string       { return senderBank }

func (a *AccountBalanceResponseSegmentV8) AccountBalance() domain.AccountBalance {
	balance := domain.AccountBalance{
		Account:       a.AccountConnection.Val().ToAccountConnection(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	return balance
}

func (a *AccountBalanceResponseSegmentV8) SepaAccountBalance() (domain.SepaAccountBalance, error) {
	balance := domain.SepaAccountBalance{
		Account:       a.AccountConnection.Val(),
		ProductName:   a.AccountProductName.Val(),
		Currency:      a.AccountCurrency.Val(),
		BookedBalance: a.BookedBalance.Balance(),
	}
	if earmarked := a.EarmarkedBalance; earmarked != nil {
		val := earmarked.Balance()
		balance.EarmarkedBalance = &val
	}
	if credit := a.CreditLimit; credit != nil {
		val := credit.Val()
		balance.CreditLimit = &val
	}
	if available := a.AvailableAmount; available != nil {
		val := available.Val()
		balance.AvailableAmount = &val
	}
	if used := a.UsedAmount; used != nil {
		val := used.Val()
		balance.UsedAmount = &val
	}
	if date := a.BookingDate; date != nil {
		val := date.Val()
		balance.BookingDate = &val
	}
	if t := a.BookingTime; t != nil {
		val := t.Val()
		*balance.BookingDate = balance.BookingDate.Add(val.Sub(time.Time{}))
	}
	if dueDate := a.DueDate; dueDate != nil {
		val := dueDate.Val()
		balance.DueDate = &val
	}
	if seizable := a.SeizableAfterMonth; seizable != nil {
		val := seizable.Val()
		balance.SeizableAfterMonth = &val
	}
	return balance, nil
}
func (a *AccountBalanceResponseSegmentV8) elements() []element.DataElement {
	return []element.DataElement{
		a.AccountConnection,
		a.AccountProductName,
		a.AccountCurrency,
		a.BookedBalance,
		a.EarmarkedBalance,
		a.CreditLimit,
		a.AvailableAmount,
		a.UsedAmount,
		a.BookingDate,
		a.BookingTime,
		a.DueDate,
		a.SeizableAfterMonth,
	}
}
