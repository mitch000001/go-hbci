package hbci

import "github.com/mitch000001/go-hbci/ast"
import "github.com/mitch000001/go-hbci/token"

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

// Phase1 reads the tokens from the lexer and transforms them into Segment-AST-Objects
func (p *Parser) Phase1(l *StringLexer) ([]*ast.Segment, error) {
	var segments []*ast.Segment
	var currentSegment *ast.Segment = ast.NewSegment(make([]token.Token, 0), make([]*ast.DataElement, 0))
	var currentDataElement *ast.DataElement = ast.NewDataElement(make([]token.Token, 0), nil)
	var tokenBuf []token.Token
	for l.HasNext() {
		tok := l.Next()
		tokenBuf = append(tokenBuf, tok)
		switch tok.Type() {
		case token.GROUP_DATA_ELEMENT_SEPARATOR:
			groupDataElement := ast.NewGroupDataElement(tokenBuf)
			currentDataElement.AddGroupDataElement(groupDataElement)
			tokenBuf = make([]token.Token, 0)
		case token.DATA_ELEMENT_SEPARATOR:
			currentDataElement.AppendTokens(tokenBuf...)
			if currentDataElement.IsDataElementGroup() {
				groupDataElement := ast.NewGroupDataElement(tokenBuf)
				currentDataElement.AddGroupDataElement(groupDataElement)
			}
			tokenBuf = make([]token.Token, 0)
			currentSegment.AddDataElement(currentDataElement)
			currentDataElement = ast.NewDataElement(make([]token.Token, 0), nil)
		case token.SEGMENT_END_MARKER:
			currentDataElement.AppendTokens(tokenBuf...)
			if currentDataElement.IsDataElementGroup() {
				groupDataElement := ast.NewGroupDataElement(tokenBuf)
				currentDataElement.AddGroupDataElement(groupDataElement)
			}
			currentSegment.AddDataElement(currentDataElement)
			segments = append(segments, currentSegment)
			currentSegment = ast.NewSegment(make([]token.Token, 0), make([]*ast.DataElement, 0))
			currentDataElement = ast.NewDataElement(make([]token.Token, 0), nil)
			tokenBuf = make([]token.Token, 0)
		}
	}
	return segments, nil
}
