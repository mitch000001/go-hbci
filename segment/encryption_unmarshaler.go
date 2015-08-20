package segment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (e *EncryptionHeaderSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], e)
	if err != nil {
		return err
	}
	e.Segment = seg
	if len(elements) > 1 {
		e.SecurityFunction = &element.AlphaNumericDataElement{}
		err = e.SecurityFunction.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		e.SecuritySupplierRole = &element.AlphaNumericDataElement{}
		err = e.SecuritySupplierRole.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 {
		e.SecurityID = &element.SecurityIdentificationDataElement{}
		err = e.SecurityID.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 {
		e.SecurityDate = &element.SecurityDateDataElement{}
		err = e.SecurityDate.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 {
		e.EncryptionAlgorithm = &element.EncryptionAlgorithmDataElement{}
		err = e.EncryptionAlgorithm.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 {
		e.KeyName = &element.KeyNameDataElement{}
		err = e.KeyName.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) > 7 {
		e.CompressionFunction = &element.AlphaNumericDataElement{}
		err = e.CompressionFunction.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	if len(elements) > 8 {
		e.Certificate = &element.CertificateDataElement{}
		err = e.Certificate.UnmarshalHBCI(elements[8])
		if err != nil {
			return err
		}
	}
	return nil
}
