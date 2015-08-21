package message

import "github.com/mitch000001/go-hbci/domain"
import "github.com/mitch000001/go-hbci/segment"

func NewCommunicationAccessMessage(fromBank domain.BankId, toBank domain.BankId, maxEntries int, aufsetzpunkt string) *CommunicationAccessMessage {
	c := &CommunicationAccessMessage{
		Request: segment.NewCommunicationAccessRequestSegment(fromBank, toBank, maxEntries, aufsetzpunkt),
	}
	c.BasicMessage = NewBasicMessage(c)
	return c
}

type CommunicationAccessMessage struct {
	*BasicMessage
	Request *segment.CommunicationAccessRequestSegment
}

func (c *CommunicationAccessMessage) HBCISegments() []segment.ClientSegment {
	return []segment.ClientSegment{
		c.Request,
	}
}
