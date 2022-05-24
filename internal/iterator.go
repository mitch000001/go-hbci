package internal

func NewIterator(data [][]byte) Iterator {
	return &arrayIterator{data: data}
}

type Iterator interface {
	Next() []byte
	NextString() string
	HasNext() bool
	Remainder() [][]byte
}

type arrayIterator struct {
	data     [][]byte
	position int
}

func (a *arrayIterator) Next() []byte {
	if !a.HasNext() {
		return nil
	}
	val := a.data[a.position]
	a.position = a.position + 1
	return val
}

func (a *arrayIterator) NextString() string {
	return string(a.Next())
}

func (a *arrayIterator) HasNext() bool {
	return a.position < len(a.data)
}

func (a *arrayIterator) Remainder() [][]byte {
	return a.data[a.position:]
}
