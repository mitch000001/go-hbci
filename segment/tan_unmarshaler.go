// Code generated by *generator.VersionedSegmentUnmarshalerGenerator; DO NOT EDIT.

package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

var (
	_	BankSegment	= &TanResponseSegmentV6{}
	_	BankSegment	= &TanResponseSegmentV7{}
)

func init() {
	v6 := TanResponseSegmentV6{}
	KnownSegments.mustAddToIndex(VersionedSegment{v6.ID(), v6.Version()}, func() Segment { return &TanResponseSegmentV6{} })
	v7 := TanResponseSegmentV7{}
	KnownSegments.mustAddToIndex(VersionedSegment{v7.ID(), v7.Version()}, func() Segment { return &TanResponseSegmentV7{} })
}

func (t *TanResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment TanResponse
	switch header.Version.Val() {
	case 6:
		segment = &TanResponseSegmentV6{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 7:
		segment = &TanResponseSegmentV7{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown segment version: %d", header.Version.Val())
	}
	t.TanResponse = segment
	return nil
}

func (t *TanResponseSegmentV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], t)
	if err != nil {
		return err
	}
	t.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		t.TANProcess = &element.AlphaNumericDataElement{}
		err = t.TANProcess.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling TANProcess: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		t.JobHash = &element.BinaryDataElement{}
		err = t.JobHash.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling JobHash: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		t.JobReference = &element.AlphaNumericDataElement{}
		err = t.JobReference.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling JobReference: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		t.Challenge = &element.AlphaNumericDataElement{}
		err = t.Challenge.UnmarshalHBCI(elements[4])
		if err != nil {
			return fmt.Errorf("error unmarshaling Challenge: %w", err)
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		t.ChallengeHHD_UC = &element.BinaryDataElement{}
		err = t.ChallengeHHD_UC.UnmarshalHBCI(elements[5])
		if err != nil {
			return fmt.Errorf("error unmarshaling ChallengeHHD_UC: %w", err)
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		t.ChallengeExpiryDate = &element.TanChallengeExpiryDate{}
		err = t.ChallengeExpiryDate.UnmarshalHBCI(elements[6])
		if err != nil {
			return fmt.Errorf("error unmarshaling ChallengeExpiryDate: %w", err)
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		t.TANMediumDescription = &element.AlphaNumericDataElement{}
		if len(elements)+1 > 7 {
			err = t.TANMediumDescription.UnmarshalHBCI(bytes.Join(elements[7:], []byte("+")))
		} else {
			err = t.TANMediumDescription.UnmarshalHBCI(elements[7])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling TANMediumDescription: %w", err)
		}
	}
	return nil
}

func (t *TanResponseSegmentV7) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("malformed marshaled value: no elements extracted")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], t)
	if err != nil {
		return err
	}
	t.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		t.TANProcess = &element.AlphaNumericDataElement{}
		err = t.TANProcess.UnmarshalHBCI(elements[1])
		if err != nil {
			return fmt.Errorf("error unmarshaling TANProcess: %w", err)
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		t.JobHash = &element.BinaryDataElement{}
		err = t.JobHash.UnmarshalHBCI(elements[2])
		if err != nil {
			return fmt.Errorf("error unmarshaling JobHash: %w", err)
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		t.JobReference = &element.AlphaNumericDataElement{}
		err = t.JobReference.UnmarshalHBCI(elements[3])
		if err != nil {
			return fmt.Errorf("error unmarshaling JobReference: %w", err)
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		t.Challenge = &element.AlphaNumericDataElement{}
		err = t.Challenge.UnmarshalHBCI(elements[4])
		if err != nil {
			return fmt.Errorf("error unmarshaling Challenge: %w", err)
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		t.ChallengeHHD_UC = &element.BinaryDataElement{}
		err = t.ChallengeHHD_UC.UnmarshalHBCI(elements[5])
		if err != nil {
			return fmt.Errorf("error unmarshaling ChallengeHHD_UC: %w", err)
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		t.ChallengeExpiryDate = &element.TanChallengeExpiryDate{}
		err = t.ChallengeExpiryDate.UnmarshalHBCI(elements[6])
		if err != nil {
			return fmt.Errorf("error unmarshaling ChallengeExpiryDate: %w", err)
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		t.TANMediumDescription = &element.AlphaNumericDataElement{}
		if len(elements)+1 > 7 {
			err = t.TANMediumDescription.UnmarshalHBCI(bytes.Join(elements[7:], []byte("+")))
		} else {
			err = t.TANMediumDescription.UnmarshalHBCI(elements[7])
		}
		if err != nil {
			return fmt.Errorf("error unmarshaling TANMediumDescription: %w", err)
		}
	}
	return nil
}
