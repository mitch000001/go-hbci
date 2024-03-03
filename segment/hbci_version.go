package segment

import (
	"fmt"
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

// SupportedHBCIVersions maps integer version codes to HBCIVersions
var SupportedHBCIVersions = map[int]HBCIVersion{
	220: HBCI220,
	300: FINTS300,
}

// HBCIVersion defines segment constructors for a HBCI version
type HBCIVersion struct {
	version                int
	PinTanEncryptionHeader func(clientSystemId string, keyName domain.KeyName) *EncryptionHeaderSegment
	RDHEncryptionHeader    func(clientSystemId string, keyName domain.KeyName, key []byte) *EncryptionHeaderSegment
	SignatureHeader        func() *SignatureHeaderSegment
	PinTanSignatureHeader  func(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	RDHSignatureHeader     func(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment
	SignatureEnd           func() *SignatureEndSegment
	SynchronisationRequest func(modus SyncMode) *SynchronisationRequestSegment
}

// Version returns the HBCI version as integer
func (v HBCIVersion) Version() int {
	return v.version
}

// Builder represents a builder which returns certain builders based on the
// provided versions
type Builder interface {
	AccountBalanceRequest(account domain.AccountConnection, allAccounts bool) (AccountBalanceRequest, error)
	SepaAccountBalanceRequest(account domain.InternationalAccountConnection, allAccounts bool) (AccountBalanceRequest, error)
	AccountTransactionRequest(account domain.AccountConnection, allAccounts bool) (*AccountTransactionRequestSegment, error)
	SepaAccountTransactionRequest(account domain.InternationalAccountConnection, allAccounts bool) (*AccountTransactionRequestSegment, error)
	StatusProtocolRequest(from, to time.Time, maxEntries int, continuationReference string) (StatusProtocolRequest, error)
	TanProcessV4Request(referencingSegmentID string) (*TanRequestSegment, error)
}

// NewBuilder returns a new Builder which uses the supported segments to
// identify which segment to use
func NewBuilder(supportedSegments []VersionedSegment) Builder {
	segments := make(map[string][]int)
	for _, s := range supportedSegments {
		segments[s.ID] = append(segments[s.ID], s.Version)
	}
	return &builder{segments}
}

type builder struct {
	supportedSegments map[string][]int
}

func (b *builder) AccountBalanceRequest(account domain.AccountConnection, allAccounts bool) (AccountBalanceRequest, error) {
	versions, ok := b.supportedSegments["HISALS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKSAL")
	}
	request, err := AccountBalanceRequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building account balance request (HKSAL): %w", err)
	}
	return request(account, allAccounts), nil
}
func (b *builder) SepaAccountBalanceRequest(account domain.InternationalAccountConnection, allAccounts bool) (AccountBalanceRequest, error) {
	versions, ok := b.supportedSegments["HISALS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKSAL")
	}
	request, err := SepaAccountBalanceRequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building account balance request (HKSAL): %w", err)
	}
	return request(account, allAccounts), nil
}
func (b *builder) AccountTransactionRequest(account domain.AccountConnection, allAccounts bool) (*AccountTransactionRequestSegment, error) {
	versions, ok := b.supportedSegments["HIKAZS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKKAZ")
	}
	request, err := AccountTransactionRequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building account transaction request (HKKAZ): %w", err)
	}
	return request(account, allAccounts), nil
}
func (b *builder) SepaAccountTransactionRequest(account domain.InternationalAccountConnection, allAccounts bool) (*AccountTransactionRequestSegment, error) {
	versions, ok := b.supportedSegments["HIKAZS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKKAZ")
	}
	request, err := SepaAccountTransactionRequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building sepa account transaction request (HKKAZ): %w", err)
	}
	return request(account, allAccounts), nil
}
func (b *builder) StatusProtocolRequest(from, to time.Time, maxEntries int, continuationReference string) (StatusProtocolRequest, error) {
	versions, ok := b.supportedSegments["HIPROS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKPRO")
	}
	request, err := StatusProtocolRequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building status protocol request (HKPRO): %w", err)
	}
	return request(from, to, maxEntries, continuationReference), nil
}
func (b *builder) TanProcessV4Request(referencingSegmentID string) (*TanRequestSegment, error) {
	versions, ok := b.supportedSegments["HITANS"]
	if !ok {
		return nil, fmt.Errorf("Segment %s not supported", "HKTAN")
	}
	request, err := TanProcess4RequestBuilder(versions)
	if err != nil {
		return nil, fmt.Errorf("error building tan request (HKTAN): %w", err)
	}
	return request(referencingSegmentID), nil
}
