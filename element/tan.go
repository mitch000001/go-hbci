package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/internal"
	"gopkg.in/yaml.v3"
)

// Gültigkeitsdatum und –uhrzeit für Challenge
//
// Datum und Uhrzeit, bis zu welchem Zeitpunkt eine TAN auf Basis der ge-
// sendeten Challenge gültig ist. Nach Ablauf der Gültigkeitsdauer wird die ent-
// sprechende TAN entwertet.
type TanChallengeExpiryDate struct {
	DataElement
	Date *DateDataElement
	Time *TimeDataElement
}

// GroupDataElements returns the grouped DataElements
func (t *TanChallengeExpiryDate) GroupDataElements() []DataElement {
	return []DataElement{
		t.Date,
		t.Time,
	}
}

// Tan2StepSubmissionParameterV6
//
// Parameter Zwei-Schritt-TAN-Einreichung, Elementversion #6
//
// Auftragsspezifische Bankparameterdaten für den Geschäftsvorfall „Zwei- Schritt-TAN-Einreichung“.
type Tan2StepSubmissionParameterV6 struct {
	DataElement
	// Ein-Schritt-Verfahren erlaubt
	//
	// Angabe, ob Ein-Schritt-Verfahren erlaubt ist oder nicht. Darüber wird das Kundenprodukt informiert,
	// ob die Einreichung von Aufträgen im Ein-Schritt- Verfahren zusätzlich zu den definierten
	// Zwei-Schritt-Verfahren zugelassen ist.
	OneStepProcessAllowed *BooleanDataElement
	// Mehr als ein TAN-pflichtiger Auftrag pro Nachricht erlaubt
	//
	// Angabe, ob in einer FinTS-Nachricht mehr als ein TAN-pflichtiger Auftrag gesendet werden darf.
	// Bei Angabe von „N“ darf in einer FinTS-Nachricht nur ein TAN-pflichtiger Auftrag enthalten sein.
	// Bei Angabe von „J“ wird die ma- ximale Anzahl der TAN-pflichtigen Aufträge analog dem
	// Geschäftsvorfallparameter „Maximale Anzahl Aufträge“ in der BPD bestimmt (vgl. [Formals], Kapitel D.6).
	// Die Option bezieht sich auf die Anzahl der in der Nachricht ent- haltenen Aufträge, nicht auf die
	// Anzahl der TANs, d. h. es ist pro Signaturab- schluss nur eine TAN erlaubt, die bei Angabe von „J“
	// aber ggf. für mehrere Aufträge gilt. Dieser Parameter gilt sowohl für das Einschritt- als auch das
	// Zwei-Schritt-Verfahren.
	MoreThanOneObligatoryTanJobAllowed *BooleanDataElement
	// Auftrags-Hashwertverfahren
	//
	// Information, welches Verfahren für die Hashwertbildung über den Kunden- auftrag verwendet werden soll.
	// Es sind nur die in [HBCI] beschriebenen Verfahren und deren Parametrisierung (Initialisierungsvektor, etc.) zulässig.
	// Codierung:
	// 0: Auftrags-Hashwert nicht unterstützt
	// 1: RIPEMD-160
	// 2: SHA-1
	JobHashMethod *CodeDataElement
	// FIXME: docs
	ProcessParameters *Tan2StepSubmissionProcessParametersV6
}

// Elements returns the elements of this DataElement.
func (t *Tan2StepSubmissionParameterV6) Elements() []DataElement {
	return []DataElement{
		t.OneStepProcessAllowed,
		t.MoreThanOneObligatoryTanJobAllowed,
		t.JobHashMethod,
		t.ProcessParameters,
	}
}

// UnmarshalHBCI unmarshals value
func (t *Tan2StepSubmissionParameterV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	oneStepProcessAllowed := &BooleanDataElement{}
	err = oneStepProcessAllowed.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	t.OneStepProcessAllowed = oneStepProcessAllowed
	moreThanOneObligatoryTanJobAllowed := &BooleanDataElement{}
	err = moreThanOneObligatoryTanJobAllowed.UnmarshalHBCI(elements[1])
	if err != nil {
		return err
	}
	t.MoreThanOneObligatoryTanJobAllowed = moreThanOneObligatoryTanJobAllowed
	t.JobHashMethod = NewCode(charset.ToUTF8(elements[2]), 1, []string{"0", "1", "2"})
	processParams := &Tan2StepSubmissionProcessParametersV6{}
	err = processParams.UnmarshalHBCI(bytes.Join(elements[3:], []byte(":")))
	if err != nil {
		return err
	}
	t.ProcessParameters = processParams
	t.DataElement = NewGroupDataElementGroup(tan2StepSubmissionParameterDEG, 4, t)
	return nil
}

func (t *Tan2StepSubmissionParameterV6) MarshalYAML() (interface{}, error) {
	return map[string]yaml.Marshaler{
		"OneStepProcessAllowed":              t.OneStepProcessAllowed,
		"MoreThanOneObligatoryTanJobAllowed": t.MoreThanOneObligatoryTanJobAllowed,
		"JobHashMethod":                      t.JobHashMethod,
		"ProcessParameters":                  t.ProcessParameters,
	}, nil
}

// Tan2StepSubmissionParametersV6 represents a slice of
// Tan2StepSubmissionParameterV6 DataElements
type Tan2StepSubmissionProcessParametersV6 struct {
	*arrayElementGroup
}

// func (t *Tan2StepSubmissionProcessParametersV6) MarshalYAML() (interface{}, error) {
// 	return t.GroupDataElements(), nil
// }

// // Val returns the underlying Tan2StepSubmissionParameters
// func (p *Tan2StepSubmissionParameters) Val() []domain.PinTanBusinessTransaction {
// 	transactions := make([]domain.PinTanBusinessTransaction, len(p.array))
// 	for i, elem := range p.array {
// 		transactions[i] = elem.(*PinTanBusinessTransactionParameter).Val()
// 	}
// 	return transactions
// }

// UnmarshalHBCI unmarshals value into the Tan2StepSubmissionParameters
func (t *Tan2StepSubmissionProcessParametersV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements)%21 != 0 {
		return fmt.Errorf("malformed marshaled value: value pairs not even")
	}
	dataElements := make([]DataElement, len(elements)/21)
	for i := 0; i < len(elements); i += 21 {
		elem := bytes.Join(elements[i:i+21], []byte(":"))
		param := &Tan2StepSubmissionProcessParameterV6{}
		err := param.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		dataElements[i/21] = param
	}
	t.arrayElementGroup = newArrayElementGroup(tan2StepSubmissionProcessParameterDEG, len(dataElements), len(dataElements), dataElements)
	return nil
}

type Tan2StepSubmissionProcessParameterV6 struct {
	DataElement
	SecurityFunction                       *CodeDataElement
	TanProcess                             *CodeDataElement
	TechnicalIDTanProcess                  *IdentificationDataElement
	ZKATanProcess                          *AlphaNumericDataElement
	ZKATanProcessVersion                   *AlphaNumericDataElement
	TwoStepProcessName                     *AlphaNumericDataElement
	TwoStepProcessMaxInputValue            *NumberDataElement
	TwoStepProcessAllowedFormat            *CodeDataElement
	TwoStepProcessReturnValueText          *AlphaNumericDataElement
	TwoStepProcessReturnValueTextMaxLength *NumberDataElement
	MultiTANAllowed                        *BooleanDataElement
	TanTimeAndDialogReference              *CodeDataElement
	JobCancellationAllowed                 *BooleanDataElement
	SMSAccountRequired                     *CodeDataElement
	IssuerAccountRequired                  *CodeDataElement
	ChallengeClassRequired                 *BooleanDataElement
	ChallengeStructured                    *BooleanDataElement
	InitializationMode                     *CodeDataElement
	TanMediumDescriptionRequired           *CodeDataElement
	HHD_UCResponseRequired                 *BooleanDataElement
	SupportedActiveTanMedia                *NumberDataElement
}

// Elements returns the elements of this DataElement.
func (t *Tan2StepSubmissionProcessParameterV6) Elements() []DataElement {
	return []DataElement{
		t.SecurityFunction,
		t.TanProcess,
		t.TechnicalIDTanProcess,
		t.ZKATanProcess,
		t.ZKATanProcessVersion,
		t.TwoStepProcessName,
		t.TwoStepProcessMaxInputValue,
		t.TwoStepProcessAllowedFormat,
		t.TwoStepProcessReturnValueText,
		t.TwoStepProcessReturnValueTextMaxLength,
		t.MultiTANAllowed,
		t.TanTimeAndDialogReference,
		t.JobCancellationAllowed,
		t.SMSAccountRequired,
		t.IssuerAccountRequired,
		t.ChallengeClassRequired,
		t.ChallengeStructured,
		t.InitializationMode,
		t.TanMediumDescriptionRequired,
		t.HHD_UCResponseRequired,
		t.SupportedActiveTanMedia,
	}
}

// // Val returns the underlying PinTanBusinessTransaction
// func (t *Tan2StepSubmissionProcessParameterV6) Val() interface{} {
// 	return domain.PinTanBusinessTransaction{
// 	}
// }

// UnmarshalHBCI unmarshals value
func (t *Tan2StepSubmissionProcessParameterV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	iter := internal.NewIterator(elements)
	t.SecurityFunction = NewCode(iter.NextString(), 3, nil)
	t.TanProcess = NewCode(iter.NextString(), 1, []string{"1", "2"})
	t.TechnicalIDTanProcess = NewIdentification(iter.NextString())
	t.ZKATanProcess = NewAlphaNumeric(iter.NextString(), 32)
	t.ZKATanProcessVersion = NewAlphaNumeric(iter.NextString(), 10)
	t.TwoStepProcessName = NewAlphaNumeric(iter.NextString(), 30)
	var twoStepProcessMaxInputValue NumberDataElement
	if err := twoStepProcessMaxInputValue.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling TwoStepProcessMaxInputValue: %v", err)
	}
	t.TwoStepProcessMaxInputValue = &twoStepProcessMaxInputValue
	t.TwoStepProcessAllowedFormat = NewCode(iter.NextString(), 1, nil)
	t.TwoStepProcessReturnValueText = NewAlphaNumeric(iter.NextString(), 30)
	var TwoStepProcessReturnValueTextMaxLength NumberDataElement
	if err := TwoStepProcessReturnValueTextMaxLength.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling TwoStepProcessReturnValueTextMaxLength: %v", err)
	}
	t.TwoStepProcessReturnValueTextMaxLength = &TwoStepProcessReturnValueTextMaxLength
	var MultiTANAllowed BooleanDataElement
	if err := MultiTANAllowed.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling MultiTANAllowed: %v", err)
	}
	t.MultiTANAllowed = &MultiTANAllowed
	t.TanTimeAndDialogReference = NewCode(iter.NextString(), 1, nil)
	var JobCancellationAllowed BooleanDataElement
	if err := JobCancellationAllowed.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling JobCancellationAllowed: %v", err)
	}
	t.JobCancellationAllowed = &JobCancellationAllowed
	t.SMSAccountRequired = NewCode(iter.NextString(), 1, []string{"0", "1", "2"})
	t.IssuerAccountRequired = NewCode(iter.NextString(), 1, []string{"0", "2"})
	var ChallengeClassRequired BooleanDataElement
	if err := ChallengeClassRequired.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling ChallengeClassRequired: %v", err)
	}
	t.ChallengeClassRequired = &ChallengeClassRequired
	var ChallengeStructured BooleanDataElement
	if err := ChallengeStructured.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling ChallengeStructured: %v", err)
	}
	t.ChallengeStructured = &ChallengeStructured
	t.InitializationMode = NewCode(iter.NextString(), -1, []string{"00", "01", "02"})
	t.TanMediumDescriptionRequired = NewCode(iter.NextString(), 1, []string{"0", "1", "2"})
	var HHD_UCResponseRequired BooleanDataElement
	if err := HHD_UCResponseRequired.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling HHD_UCResponseRequired: %v", err)
	}
	t.HHD_UCResponseRequired = &HHD_UCResponseRequired
	var SupportedActiveTanMedia NumberDataElement
	if err := SupportedActiveTanMedia.UnmarshalHBCI(iter.Next()); err != nil {
		return fmt.Errorf("error unmarshaling SupportedActiveTanMedia: %v", err)
	}
	t.SupportedActiveTanMedia = &SupportedActiveTanMedia
	t.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, t)
	return nil
}

func (t Tan2StepSubmissionProcessParameterV6) MarshalYAML() (interface{}, error) {
	return map[string]yaml.Marshaler{
		"SecurityFunction":                       t.SecurityFunction,
		"TanProcess":                             t.TanProcess,
		"TechnicalIDTanProcess":                  t.TechnicalIDTanProcess,
		"ZKATanProcess":                          t.ZKATanProcess,
		"ZKATanProcessVersion":                   t.ZKATanProcessVersion,
		"TwoStepProcessName":                     t.TwoStepProcessName,
		"TwoStepProcessMaxInputValue":            t.TwoStepProcessMaxInputValue,
		"TwoStepProcessAllowedFormat":            t.TwoStepProcessAllowedFormat,
		"TwoStepProcessReturnValueText":          t.TwoStepProcessReturnValueText,
		"TwoStepProcessReturnValueTextMaxLength": t.TwoStepProcessReturnValueTextMaxLength,
		"MultiTANAllowed":                        t.MultiTANAllowed,
		"TanTimeAndDialogReference":              t.TanTimeAndDialogReference,
		"JobCancellationAllowed":                 t.JobCancellationAllowed,
		"SMSAccountRequired":                     t.SMSAccountRequired,
		"IssuerAccountRequired":                  t.IssuerAccountRequired,
		"ChallengeClassRequired":                 t.ChallengeClassRequired,
		"ChallengeStructured":                    t.ChallengeStructured,
		"InitializationMode":                     t.InitializationMode,
		"TanMediumDescriptionRequired":           t.TanMediumDescriptionRequired,
		"HHD_UCResponseRequired":                 t.HHD_UCResponseRequired,
		"SupportedActiveTanMedia":                t.SupportedActiveTanMedia,
	}, nil
}
