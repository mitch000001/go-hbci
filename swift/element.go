package swift

import (
	"fmt"
	"strconv"
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

type CustomField struct {
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
