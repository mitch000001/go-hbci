package hbci

import (
	"bytes"
	"fmt"
	"strings"
)

func NewDataElement(tokens []Token) DataElement {
	return DataElement{
		tokens: tokens,
	}
}

type DataElement struct {
	*DataElementGroup
	tokens []Token
}

func (d *DataElement) Tokens() []Token {
	if d.IsDataElementGroup() {
		var tokens []Token
		for _, gd := range d.groupDataElements {
			tokens = append(tokens, gd.tokens...)
		}
		return tokens
	} else {
		return d.tokens
	}
}

func (d *DataElement) IsDataElementGroup() bool {
	return d.DataElementGroup != nil
}

func (d *DataElement) AddGroupDataElement(groupDataElement GroupDataElement) {
	if d.DataElementGroup == nil {
		d.DataElementGroup = &DataElementGroup{}
	}
	d.groupDataElements = append(d.groupDataElements, groupDataElement)
	d.tokens = nil
}

type BinaryDataElement struct {
	token   Token
	content []byte
}

type AlphaNumericDataElement struct {
	token Token
	text  string
}

type TextDataElement struct {
	token Token
	text  string
}

type DigitDataElement struct {
	token Token
	value int64
}

type NumericDataElement struct {
	token Token
	value int64
}

type FloatDataElement struct {
	token Token
	value float64
}

type DataElementGroup struct {
	tokens            []Token
	groupDataElements []GroupDataElement
}

type GroupDataElement struct {
	tokens []Token
}

func NewSegment(tokens []Token, dataElements []DataElement) *Segment {
	return &Segment{
		tokens:       tokens,
		dataElements: dataElements,
	}
}

type Segment struct {
	tokens       []Token
	dataElements []DataElement
}

func (s *Segment) AddDataElement(dataElement DataElement) {
	s.dataElements = append(s.dataElements, dataElement)
	s.tokens = append(s.tokens, dataElement.Tokens()...)
}

func (s *Segment) GoString() string {
	var buf bytes.Buffer
	buf.WriteString("hbci.Segment{tokens:[")
	var tokenBuf []string
	for _, tok := range s.tokens {
		tokenBuf = append(tokenBuf, fmt.Sprintf("{%s}", tok))
	}
	fmt.Fprintf(&buf, "%s],dataElements:[", strings.Join(tokenBuf, " "))
	var dataElemBuf []string
	for _, elem := range s.dataElements {
		dataElemBuf = append(dataElemBuf, fmt.Sprintf("{%s}", elem))
	}
	fmt.Fprintf(&buf, "%s]}", strings.Join(dataElemBuf, " "))
	return buf.String()
}
