package hbci

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

func (p *Parser) Phase1(l *Lexer) ([]*Segment, error) {
	var segments []*Segment
	var currentSegment *Segment = NewSegment(make([]Token, 0), make([]DataElement, 0))
	var currentDataElement DataElement = NewDataElement(make([]Token, 0))
	var tokenBuf []Token
	for l.HasNext() {
		token := l.Next()
		tokenBuf = append(tokenBuf, token)
		switch token.typ {
		case GROUP_DATA_ELEMENT_SEPARATOR:
			groupDataElement := GroupDataElement{tokens: tokenBuf}
			currentDataElement.AddGroupDataElement(groupDataElement)
			tokenBuf = make([]Token, 0)
		case DATA_ELEMENT_SEPARATOR:
			currentDataElement.tokens = append(currentDataElement.tokens, tokenBuf...)
			if currentDataElement.IsDataElementGroup() {
				groupDataElement := GroupDataElement{tokens: tokenBuf}
				currentDataElement.AddGroupDataElement(groupDataElement)
			}
			tokenBuf = make([]Token, 0)
			currentSegment.AddDataElement(currentDataElement)
			currentDataElement = NewDataElement(make([]Token, 0))
		case SEGMENT_END_MARKER:
			currentDataElement.tokens = append(currentDataElement.tokens, tokenBuf...)
			if currentDataElement.IsDataElementGroup() {
				groupDataElement := GroupDataElement{tokens: tokenBuf}
				currentDataElement.AddGroupDataElement(groupDataElement)
			}
			currentSegment.AddDataElement(currentDataElement)
			segments = append(segments, currentSegment)
			currentSegment = NewSegment(make([]Token, 0), make([]DataElement, 0))
			currentDataElement = NewDataElement(make([]Token, 0))
			tokenBuf = make([]Token, 0)
		}
	}
	return segments, nil
}
