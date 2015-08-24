package segment

import "fmt"

type VersionedSegment struct {
	ID      string
	Version int
}

func (v VersionedSegment) String() string {
	return fmt.Sprintf("%s:%d", v.ID, v.Version)
}

type SegmentIndex map[VersionedSegment]func() Segment

func (s SegmentIndex) UnmarshalerForSegment(segmentId VersionedSegment) (Unmarshaler, error) {
	segmentFn, ok := s[segmentId]
	if ok {
		unmarshaler, ok := segmentFn().(Unmarshaler)
		if ok {
			return unmarshaler, nil
		} else {
			return nil, fmt.Errorf("Segment does not implement the Unmarshaler interface")
		}
	} else {
		return nil, fmt.Errorf("Segment not in index: %q", segmentId)
	}
}

func (s SegmentIndex) IsIndexed(segmentId VersionedSegment) bool {
	_, ok := s[segmentId]
	return ok
}

func (s SegmentIndex) IsUnmarshaler(segmentId VersionedSegment) bool {
	segmentFn, ok := s[segmentId]
	if ok {
		_, ok := segmentFn().(Unmarshaler)
		return ok
	} else {
		return false
	}
}

var KnownSegments = SegmentIndex{
	VersionedSegment{"HNHBK", 3}:  func() Segment { return &MessageHeaderSegment{} },
	VersionedSegment{"HNHBS", 1}:  func() Segment { return &MessageEndSegment{} },
	VersionedSegment{"HNVSK", 2}:  func() Segment { return &EncryptionHeaderV2{} },
	VersionedSegment{"HNVSK", 3}:  func() Segment { return &EncryptionHeaderSegmentV3{} },
	VersionedSegment{"HNVSD", 1}:  func() Segment { return &EncryptedDataSegment{} },
	VersionedSegment{"HIRMG", 2}:  func() Segment { return &MessageAcknowledgement{} },
	VersionedSegment{"HIRMS", 2}:  func() Segment { return &SegmentAcknowledgement{} },
	VersionedSegment{"HISYN", 3}:  func() Segment { return &SynchronisationResponseSegment{} },
	VersionedSegment{"HIKIM", 2}:  func() Segment { return &BankAnnouncementSegment{} },
	VersionedSegment{"HIBPA", 2}:  func() Segment { return &CommonBankParameterV2{} },
	VersionedSegment{"HIBPA", 3}:  func() Segment { return &CommonBankParameterV3{} },
	VersionedSegment{"DIPINS", 1}: func() Segment { return &PinTanBusinessTransactionParamsSegment{} },
	VersionedSegment{"HIUPA", 2}:  func() Segment { return &CommonUserParameterDataV2{} },
	VersionedSegment{"HIUPA", 3}:  func() Segment { return &CommonUserParameterDataV3{} },
	VersionedSegment{"HIUPA", 4}:  func() Segment { return &CommonUserParameterDataV4{} },
	VersionedSegment{"HIUPD", 4}:  func() Segment { return &AccountInformationV4{} },
	VersionedSegment{"HIUPD", 5}:  func() Segment { return &AccountInformationV5{} },
	VersionedSegment{"HIUPD", 6}:  func() Segment { return &AccountInformationV6{} },
}
