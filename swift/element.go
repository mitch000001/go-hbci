package swift

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
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
	elements, err := ExtractTagElements(value)
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
	elements, err := ExtractTagElements(value)
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
	elements, err := ExtractTagElements(value)
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
	BankID             int
	AccountID          int
	Name               string
	MessageKeyAddition int
	Purpose2           string
}

var customFieldTagFieldKeys = [][]byte{
	[]byte("?00"),
	[]byte("?10"),
	[]byte("?20"),
	[]byte("?21"),
	[]byte("?22"),
	[]byte("?23"),
	[]byte("?24"),
	[]byte("?25"),
	[]byte("?26"),
	[]byte("?27"),
	[]byte("?28"),
	[]byte("?29"),
	[]byte("?30"),
	[]byte("?31"),
	[]byte("?32"),
	[]byte("?33"),
	[]byte("?34"),
	[]byte("?60"),
	[]byte("?61"),
	[]byte("?62"),
	[]byte("?63"),
}

func (c *CustomFieldTag) Unmarshal(value []byte) error {
	elements, err := ExtractTagElements(value)
	if err != nil {
		return err
	}
	if len(elements) != 2 {
		return fmt.Errorf("%T: Malformed marshaled value", c)
	}
	c.Tag = string(elements[0])
	tId, err := strconv.Atoi(string(elements[1][:3]))
	if err != nil {
		return err
	}
	c.TransactionID = tId
	marshaledFields := elements[1][3:]
	var fields []fieldKeyIndex
	for _, fieldKey := range customFieldTagFieldKeys {
		if idx := bytes.Index(marshaledFields, fieldKey); idx != -1 {
			fields = append(fields, fieldKeyIndex{string(fieldKey), idx})
		}
	}
	getFieldValue := func(currentFieldKeyIndex, nextFieldKeyIndex int) []byte {
		return marshaledFields[currentFieldKeyIndex+3 : nextFieldKeyIndex]
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
		case strings.HasPrefix(fieldKey, "?00"):
			c.BookingText = string(fieldValue)
		case strings.HasPrefix(fieldKey, "?10"):
			c.PrimanotenNumber = string(fieldValue)
		case strings.HasPrefix(fieldKey, "?2"):
			c.Purpose += strings.Replace(string(fieldValue), "\r\n", "", -1)
		case strings.HasPrefix(fieldKey, "?30"):
			bankId, err := strconv.Atoi(string(fieldValue))
			if err != nil {
				return err
			}
			c.BankID = bankId
		case strings.HasPrefix(fieldKey, "?31"):
			accountId, err := strconv.Atoi(string(fieldValue))
			if err != nil {
				return err
			}
			c.AccountID = accountId
		case strings.HasPrefix(fieldKey, "?32"):
			c.Name = string(fieldValue)
		case strings.HasPrefix(fieldKey, "?33"):
			c.Name += " " + string(fieldValue)
		case strings.HasPrefix(fieldKey, "?34"):
			messageKeyAddition, err := strconv.Atoi(string(fieldValue))
			if err != nil {
				return err
			}
			c.MessageKeyAddition = messageKeyAddition
		case strings.HasPrefix(fieldKey, "?6"):
			c.Purpose2 += strings.Replace(string(fieldValue), "\r\n", "", -1)
		}
	}
	return nil
}

type fieldKeyIndex struct {
	fieldKey string
	index    int
}
