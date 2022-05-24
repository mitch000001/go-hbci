package segment

import "fmt"

type VersionedSegment struct {
	ID      string
	Version int
}

func (v VersionedSegment) String() string {
	return fmt.Sprintf("%s:%d", v.ID, v.Version)
}

type SegmentIndex struct {
	segmentMap map[VersionedSegment]func() Segment
}

func (s SegmentIndex) UnmarshalerForSegment(segmentId VersionedSegment) (Unmarshaler, error) {
	segmentFn, ok := s.segmentMap[segmentId]
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
	_, ok := s.segmentMap[segmentId]
	return ok
}

func (s SegmentIndex) IsUnmarshaler(segmentId VersionedSegment) bool {
	segmentFn, ok := s.segmentMap[segmentId]
	if ok {
		_, ok := segmentFn().(Unmarshaler)
		return ok
	} else {
		return false
	}
}

func (s *SegmentIndex) addToIndex(segmentIdentifier VersionedSegment, segmentProviderFn func() Segment) error {
	if s.IsIndexed(segmentIdentifier) {
		return fmt.Errorf("Segment already in index: %s", segmentIdentifier)
	}
	s.segmentMap[segmentIdentifier] = segmentProviderFn
	return nil
}

func (s *SegmentIndex) mustAddToIndex(segmentIdentifier VersionedSegment, segmentProviderFn func() Segment) {
	err := s.addToIndex(segmentIdentifier, segmentProviderFn)
	if err != nil {
		panic(err)
	}
}

var KnownSegments = SegmentIndex{segmentMap: make(map[VersionedSegment]func() Segment)}

func init() {
	KnownSegments.mustAddToIndex(VersionedSegment{MessageHeaderID, 3}, func() Segment { return &MessageHeaderSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HNHBS", 1}, func() Segment { return &MessageEndSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{EncryptionHeaderSegmentID, 2}, func() Segment { return &EncryptionHeaderV2{} })
	KnownSegments.mustAddToIndex(VersionedSegment{EncryptionHeaderSegmentID, 3}, func() Segment { return &EncryptionHeaderSegmentV3{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HNVSD", 1}, func() Segment { return &EncryptedDataSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIRMG", 2}, func() Segment { return &MessageAcknowledgement{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIRMS", 2}, func() Segment { return &SegmentAcknowledgement{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HISYN", 3}, func() Segment { return &SynchronisationResponseSegmentV3{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HISYN", 4}, func() Segment { return &SynchronisationResponseSegmentV4{} })
	KnownSegments.mustAddToIndex(VersionedSegment{BankAnnouncementID, 2}, func() Segment { return &BankAnnouncementSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{CommonBankParameterID, 2}, func() Segment { return &CommonBankParameterV2{} })
	KnownSegments.mustAddToIndex(VersionedSegment{CommonBankParameterID, 3}, func() Segment { return &CommonBankParameterV3{} })
	KnownSegments.mustAddToIndex(VersionedSegment{PinTanBusinessTransactionParamsID, 1}, func() Segment { return &PinTanBusinessTransactionParamsSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{PinTanBankParameterID, 1}, func() Segment { return &PinTanBankParameterV1{} })
	KnownSegments.mustAddToIndex(VersionedSegment{CommonUserParameterDataID, 2}, func() Segment { return &CommonUserParameterDataV2{} })
	KnownSegments.mustAddToIndex(VersionedSegment{CommonUserParameterDataID, 3}, func() Segment { return &CommonUserParameterDataV3{} })
	KnownSegments.mustAddToIndex(VersionedSegment{CommonUserParameterDataID, 4}, func() Segment { return &CommonUserParameterDataV4{} })
	KnownSegments.mustAddToIndex(VersionedSegment{AccountInformationID, 4}, func() Segment { return &AccountInformationV4{} })
	KnownSegments.mustAddToIndex(VersionedSegment{AccountInformationID, 5}, func() Segment { return &AccountInformationV5{} })
	KnownSegments.mustAddToIndex(VersionedSegment{AccountInformationID, 6}, func() Segment { return &AccountInformationV6{} })
	KnownSegments.mustAddToIndex(VersionedSegment{AccountInformationID, 7}, func() Segment { return &AccountInformationV7{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HISAL", 5}, func() Segment { return &AccountBalanceResponseSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIKIF", 1}, func() Segment { return &AccountInformationResponseSegment{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIKAZ", 5}, func() Segment { return &AccountTransactionResponseSegmentV5{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIKAZ", 6}, func() Segment { return &AccountTransactionResponseSegmentV6{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIKAZ", 7}, func() Segment { return &AccountTransactionResponseSegmentV7{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIPRO", 3}, func() Segment { return &StatusProtocolResponseSegmentV3{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HIPRO", 4}, func() Segment { return &StatusProtocolResponseSegmentV4{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HITANS", 6}, func() Segment { return &TanBankParameterV6{} })
	KnownSegments.mustAddToIndex(VersionedSegment{"HITAN", 6}, func() Segment { return &TanResponseSegmentV6{} })
}
