package ast

import (
	"bytes"
	"fmt"
	goast "go/ast"
	gotoken "go/token"
	"strings"

	"github.com/mitch000001/go-hbci/token"
)

type Node interface {
	RawTokens() token.Tokens
	Pos() int
	End() int
}

// astNode implements the method set nescessary to use the node with ast.Walk
type astNode struct {
	node Node
}

func (a astNode) Pos() gotoken.Pos {
	return gotoken.Pos(a.node.Pos() + 1)
}

func (a astNode) End() gotoken.Pos {
	return gotoken.Pos(a.node.End() + 1)
}

type astVisitor struct {
	visitor Visitor
}

func (a astVisitor) Visit(node goast.Node) goast.Visitor {

	return nil
}

//  Walk traverses an AST in depth-first order:
//  It starts by calling v.Visit(node); node must not be nil.
//  If the visitor w returned by v.Visit(node) is not nil,
//  Walk is invoked recursively with visitor w for each of the non-nil children of node,
//  followed by a call of w.Visit(nil).
func Walk(v Visitor, node Node) {
}

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func NewDataElement(tokens []token.Token, dataElementGroup *DataElementGroup) *DataElement {
	return &DataElement{
		DataElementGroup: dataElementGroup,
		tokens:           tokens,
	}
}

type DataElement struct {
	*DataElementGroup
	tokens []token.Token
}

func (d *DataElement) Tokens() []token.Token {
	if d.IsDataElementGroup() {
		var tokens []token.Token
		for _, gd := range d.groupDataElements {
			tokens = append(tokens, gd.tokens...)
		}
		return tokens
	} else {
		return d.tokens
	}
}

func (d *DataElement) AppendTokens(tokens ...token.Token) {
	d.tokens = append(d.tokens, tokens...)
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

func (d *DataElement) GroupDataElements() []GroupDataElement {
	return d.DataElementGroup.groupDataElements
}

type BinaryDataElement struct {
	token   token.Token
	content []byte
}

type AlphaNumericDataElement struct {
	token token.Token
	text  string
}

type TextDataElement struct {
	token token.Token
	text  string
}

type DigitDataElement struct {
	token token.Token
	value int64
}

type NumericDataElement struct {
	token token.Token
	value int64
}

type FloatDataElement struct {
	token token.Token
	value float64
}

func NewDataElementGroup(tokens token.Tokens, groupDataElements ...GroupDataElement) *DataElementGroup {
	return &DataElementGroup{tokens: tokens, groupDataElements: groupDataElements}
}

type DataElementGroup struct {
	tokens            []token.Token
	groupDataElements []GroupDataElement
}

func NewGroupDataElement(tokens []token.Token) GroupDataElement {
	return GroupDataElement{tokens: tokens}
}

type GroupDataElement struct {
	tokens []token.Token
}

type EncryptedMessage struct {
	*Message
	EncryptionHeader *EncryptionHeader
	EncryptedData    *EncryptedData
}

type EncryptionHeader struct{}

type EncryptedData struct{}

type Message struct {
	Header *MessageHeader
	End    *MessageEnd
}

type BankMessage struct {
	*Message
	SignatureBegin         *SignatureHeader
	SignatureEnd           *SignatureEnd
	MessageAcknowledgement *MessageAcknowledgement
	SegmentAcknowledgement *SegmentAcknowledgement
	DataSegments           []*DataSegment
}

type MessageAcknowledgement struct{}
type SegmentAcknowledgement struct{}
type DataSegment struct{}

type ClientMessage struct {
	*Message
	SignatureBegin *SignatureHeader
	SignatureEnd   *SignatureEnd
	Jobs           []*Job
}

type MessageHeader struct {
	*Segment
}

type MessageHeaderSegment struct {
	Header      *SegmentHeader
	Size        *MessageSize
	HBCIVersion *HBCIVersion
	DialogID    *DialogIdentifier
	Number      *MessageNumber
	Ref         *ReferenceMessage
}

type MessageSize struct{}
type HBCIVersion struct{}
type DialogIdentifier struct{}
type MessageNumber struct{}
type ReferenceMessage struct{}

type MessageEnd struct{}
type SignatureHeader struct{}
type SignatureEnd struct{}
type Job struct{}

type SegmentSequence struct{}

type SegmentHeader struct {
	*DataElementGroup
	ID      *SegmentIdentifier
	Number  *SegmentNumber
	Version *SegmentVersion
	Ref     *ReferenceSegment
}

type ReferenceSegment struct{}
type SegmentIdentifier struct{}
type SegmentNumber struct{}
type SegmentVersion struct{}

func NewSegment(tokens []token.Token, dataElements []*DataElement) *Segment {
	return &Segment{
		tokens:       tokens,
		dataElements: dataElements,
	}
}

type Segment struct {
	Header       *SegmentHeader
	tokens       []token.Token
	dataElements []*DataElement
}

func (s *Segment) DataElements() []*DataElement {
	return s.dataElements
}

func (s *Segment) AddDataElement(dataElement *DataElement) {
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
