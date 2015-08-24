package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (e *EncryptionHeaderSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment encryptionHeaderSegment
	switch header.Version.Val() {
	case 2:
		segment = &EncryptionHeaderV2{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 3:
		segment = &EncryptionHeaderSegmentV3{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	e.encryptionHeaderSegment = segment
	return nil
}

func (e *EncryptionHeaderV2) UnmarshalHBCI(value []byte) error {
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
	e.ClientSegment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		e.SecurityFunction = &element.AlphaNumericDataElement{}
		err = e.SecurityFunction.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		e.SecuritySupplierRole = &element.AlphaNumericDataElement{}
		err = e.SecuritySupplierRole.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		e.SecurityID = &element.SecurityIdentificationDataElement{}
		err = e.SecurityID.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		e.SecurityDate = &element.SecurityDateDataElement{}
		err = e.SecurityDate.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		e.EncryptionAlgorithm = &element.EncryptionAlgorithmDataElement{}
		err = e.EncryptionAlgorithm.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		e.KeyName = &element.KeyNameDataElement{}
		err = e.KeyName.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		e.CompressionFunction = &element.AlphaNumericDataElement{}
		err = e.CompressionFunction.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	if len(elements) > 8 && len(elements[8]) > 0 {
		e.Certificate = &element.CertificateDataElement{}
		if len(elements)+1 > 8 {
			err = e.Certificate.UnmarshalHBCI(bytes.Join(elements[8:], []byte("+")))
		} else {
			err = e.Certificate.UnmarshalHBCI(elements[8])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EncryptionHeaderSegmentV3) UnmarshalHBCI(value []byte) error {
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
	e.ClientSegment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		e.SecurityProfile = &element.SecurityProfileDataElement{}
		err = e.SecurityProfile.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		e.SecurityFunction = &element.CodeDataElement{}
		err = e.SecurityFunction.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		e.SecuritySupplierRole = &element.CodeDataElement{}
		err = e.SecuritySupplierRole.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		e.SecurityID = &element.SecurityIdentificationDataElement{}
		err = e.SecurityID.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		e.SecurityDate = &element.SecurityDateDataElement{}
		err = e.SecurityDate.UnmarshalHBCI(elements[5])
		if err != nil {
			return err
		}
	}
	if len(elements) > 6 && len(elements[6]) > 0 {
		e.EncryptionAlgorithm = &element.EncryptionAlgorithmDataElement{}
		err = e.EncryptionAlgorithm.UnmarshalHBCI(elements[6])
		if err != nil {
			return err
		}
	}
	if len(elements) > 7 && len(elements[7]) > 0 {
		e.KeyName = &element.KeyNameDataElement{}
		err = e.KeyName.UnmarshalHBCI(elements[7])
		if err != nil {
			return err
		}
	}
	if len(elements) > 8 && len(elements[8]) > 0 {
		e.CompressionFunction = &element.CodeDataElement{}
		err = e.CompressionFunction.UnmarshalHBCI(elements[8])
		if err != nil {
			return err
		}
	}
	if len(elements) > 9 && len(elements[9]) > 0 {
		e.Certificate = &element.CertificateDataElement{}
		if len(elements)+1 > 9 {
			err = e.Certificate.UnmarshalHBCI(bytes.Join(elements[9:], []byte("+")))
		} else {
			err = e.Certificate.UnmarshalHBCI(elements[9])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
