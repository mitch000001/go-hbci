package swift

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

func NewMT940Messages(data []byte) *MT940Messages {
	return &MT940Messages{
		Data:      data,
		timestamp: time.Now().UnixNano(),
	}
}

type MT940Messages struct {
	Data      []byte
	timestamp int64
}

func MergeMT940Messages(messages ...*MT940Messages) *MT940Messages {
	sortedMessages := sortedByTimestamp(messages)
	sort.Sort(sortedMessages)
	merged := &MT940Messages{
		Data:      []byte{},
		timestamp: time.Now().UnixNano(),
	}
	for _, msg := range sortedMessages {
		if msg == nil {
			continue
		}
		merged.Data = append(merged.Data, msg.Data...)
	}
	return merged
}

func NewMT942Messages(data []byte) *MT942Messages {
	return &MT942Messages{
		Data:      data,
		timestamp: time.Now().UnixNano(),
	}
}

type MT942Messages struct {
	Data      []byte
	timestamp int64
}

func MergeMT942Messages(messages ...*MT942Messages) *MT942Messages {
	sortedMessages := mt942sortedByTimestamp(messages)
	sort.Sort(sortedMessages)
	merged := &MT942Messages{
		Data:      []byte{},
		timestamp: time.Now().UnixNano(),
	}
	for _, msg := range sortedMessages {
		if msg == nil {
			continue
		}
		merged.Data = append(merged.Data, msg.Data...)
	}
	return merged
}

type sortedByTimestamp []*MT940Messages

// Len is the number of elements in the collection.
func (s sortedByTimestamp) Len() int {
	return len(s)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (s sortedByTimestamp) Less(i int, j int) bool {
	return s[i].timestamp < s[j].timestamp
}

// Swap swaps the elements with indexes i and j.
func (s sortedByTimestamp) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type mt942sortedByTimestamp []*MT942Messages

// Len is the number of elements in the collection.
func (s mt942sortedByTimestamp) Len() int {
	return len(s)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (s mt942sortedByTimestamp) Less(i int, j int) bool {
	return s[i].timestamp < s[j].timestamp
}

// Swap swaps the elements with indexes i and j.
func (s mt942sortedByTimestamp) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type MT940Unmarshaler interface {
	UnmarshalMT940([]byte) ([]domain.BookedAccountTransactions, error)
}

func NewMT940MessagesUnmarshaler() MT940Unmarshaler {
	return &mt940MessagesUnmarshaler{}
}

type mt940MessagesUnmarshaler struct{}

func (m *mt940MessagesUnmarshaler) UnmarshalMT940(value []byte) ([]domain.BookedAccountTransactions, error) {
	messageExtractor := NewMessageExtractor(value)
	messages, err := messageExtractor.Extract()
	if err != nil {
		return nil, fmt.Errorf("error extracting messages: %w", err)
	}
	var errors errorList
	var transactions []domain.BookedAccountTransactions
	for _, message := range messages {
		tr := &MT940{}
		err = tr.Unmarshal(message)
		if err != nil {
			errors = append(errors, fmt.Errorf("error unmarshaling MT940: %w", err))
		}
		transactions = append(transactions, tr.BookedAccountTransactions())
	}
	if len(errors) != 0 {
		return nil, errors
	}
	return transactions, nil
}

type MT942Unmarshaler interface {
	UnmarshalMT942([]byte) ([]domain.UnbookedAccountTransactions, error)
}

func NewMT942MessagesUnmarshaler() MT942Unmarshaler {
	return &mt942MessagesUnmarshaler{}
}

type mt942MessagesUnmarshaler struct{}

func (m *mt942MessagesUnmarshaler) UnmarshalMT942(value []byte) ([]domain.UnbookedAccountTransactions, error) {
	messageExtractor := NewMessageExtractor(value)
	messages, err := messageExtractor.Extract()
	if err != nil {
		return nil, fmt.Errorf("error extracting messages: %w", err)
	}
	var errors errorList
	var transactions []domain.UnbookedAccountTransactions
	for _, message := range messages {
		tr := &MT942{}
		err = tr.Unmarshal(message)
		if err != nil {
			errors = append(errors, fmt.Errorf("error unmarshaling MT942: %w", err))
		}
		transactions = append(transactions, tr.UnbookedAccountTransactions())
	}
	if len(errors) != 0 {
		return nil, errors
	}
	return transactions, nil
}

type errorList []error

func (e errorList) Error() string {
	errs := make([]string, len(e))
	for i, err := range e {
		errs[i] = err.Error()
	}
	return strings.Join(errs, ",")
}
