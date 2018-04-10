package swift

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/internal"
)

type Tag interface {
	Unmarshal([]byte) error
	Value() interface{}
	ID() string
}

type tag struct {
	id    string
	value interface{}
}

func (t *tag) ID() string         { return t.id }
func (t *tag) Value() interface{} { return t.value }

type AlphaNumericTag struct {
	*tag
}

func (a *AlphaNumericTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", a)
	}
	id := string(elements[0])
	val := string(elements[1])
	a.tag = &tag{id: id, value: val}
	return nil
}

func (a *AlphaNumericTag) Val() string {
	return a.value.(string)
}

type NumberTag struct {
	*tag
}

func (n *NumberTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", n)
	}
	id := string(elements[0])
	num, err := strconv.Atoi(string(elements[1]))
	if err != nil {
		return fmt.Errorf("%T: Error while unmarshaling: %v", n, err)
	}
	n.tag = &tag{id: id, value: num}
	return nil
}

func (n *NumberTag) Val() int {
	return n.value.(int)
}

type FloatTag struct {
	*tag
}

func (f *FloatTag) Unmarshal(value []byte) error {
	elements, err := extractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", f)
	}
	id := string(elements[0])
	num, err := strconv.ParseFloat(string(elements[1]), 64)
	if err != nil {
		return fmt.Errorf("%T: Error while unmarshaling: %v", f, err)
	}
	f.tag = &tag{id: id, value: num}
	return nil
}

func (f *FloatTag) Val() float64 {
	return f.value.(float64)
}

type CustomFieldTag struct {
	Tag                string
	TransactionID      int
	BookingText        string
	PrimanotenNumber   string
	Purpose            string
	BankID             string
	AccountID          string
	Name               string
	MessageKeyAddition int
	Purpose2           string
}

var customFieldTagFieldKeys = [][]byte{
	[]byte{'?', '0', '0'},
	[]byte{'?', '1', '0'},
	[]byte{'?', '2', '0'},
	[]byte{'?', '2', '1'},
	[]byte{'?', '2', '2'},
	[]byte{'?', '2', '3'},
	[]byte{'?', '2', '4'},
	[]byte{'?', '2', '5'},
	[]byte{'?', '2', '6'},
	[]byte{'?', '2', '7'},
	[]byte{'?', '2', '8'},
	[]byte{'?', '2', '9'},
	[]byte{'?', '3', '0'},
	[]byte{'?', '3', '1'},
	[]byte{'?', '3', '2'},
	[]byte{'?', '3', '3'},
	[]byte{'?', '3', '4'},
	[]byte{'?', '6', '0'},
	[]byte{'?', '6', '1'},
	[]byte{'?', '6', '2'},
	[]byte{'?', '6', '3'},
}

// Unmarshal unmarshals the tag bytes into c
func (c *CustomFieldTag) Unmarshal(value []byte) error {
	tag, err := extractRawTag(value)
	if err != nil {
		return err
	}
	c.Tag = tag.ID
	tID, err := strconv.Atoi(charset.ToUTF8(tag.Value[:3]))
	if err != nil {
		return err
	}
	c.TransactionID = tID
	marshaledFields := tag.Value[3:]
	marshaledFields = bytes.Replace(
		marshaledFields, []byte{'\r', '\n'}, []byte{}, -1,
	)
	var fields []fieldKeyIndex
	for _, fieldKey := range customFieldTagFieldKeys {
		if idx := bytes.Index(marshaledFields, fieldKey); idx != -1 {
			fields = append(fields, fieldKeyIndex{fieldKey, idx})
		}
	}

	getFieldValue := func(currentFieldKeyIndex, nextFieldKeyIndex int) string {
		return charset.ToUTF8(
			marshaledFields[currentFieldKeyIndex+3 : nextFieldKeyIndex],
		)
	}
	for i, fieldKeyIndex := range fields {
		var nextFieldKeyIndex int
		if len(fields)-1 == i {
			nextFieldKeyIndex = len(marshaledFields)
		} else {
			nextFieldKeyIndex = fields[i+1].index
		}

		fieldValue := getFieldValue(fieldKeyIndex.index, nextFieldKeyIndex)

		switch fieldKey := fieldKeyIndex.fieldKey; {
		case bytes.HasPrefix(fieldKey, []byte{'?', '0', '0'}):
			c.BookingText = fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '1', '0'}):
			c.PrimanotenNumber = fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '2'}):
			c.Purpose += fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '3', '0'}):
			c.BankID = fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '3', '1'}):
			c.AccountID = fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '3', '2'}):
			c.Name = fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '3', '3'}):
			c.Name += " " + fieldValue
		case bytes.HasPrefix(fieldKey, []byte{'?', '3', '4'}):
			messageKeyAddition, err := strconv.Atoi(fieldValue)
			if err != nil {
				return err
			}
			c.MessageKeyAddition = messageKeyAddition
		case bytes.HasPrefix(fieldKey, []byte{'?', '6'}):
			c.Purpose2 += fieldValue
		default:
			internal.Debug.Printf("Unmarshal CustomFieldTag: unknown fieldKey: %s\n", fieldKey)
		}
	}
	return nil
}

type fieldKeyIndex struct {
	fieldKey []byte
	index    int
}
