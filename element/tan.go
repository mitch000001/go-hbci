package element

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/charset"
	"github.com/mitch000001/go-hbci/internal"
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
	DataElement `yaml:"-"`
	// Ein-Schritt-Verfahren erlaubt
	//
	// Angabe, ob Ein-Schritt-Verfahren erlaubt ist oder nicht. Darüber wird das Kundenprodukt informiert,
	// ob die Einreichung von Aufträgen im Ein-Schritt- Verfahren zusätzlich zu den definierten
	// Zwei-Schritt-Verfahren zugelassen ist.
	OneStepProcessAllowed *BooleanDataElement `yaml:"OneStepProcessAllowed"`
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
	MoreThanOneObligatoryTanJobAllowed *BooleanDataElement `yaml:"MoreThanOneObligatoryTanJobAllowed"`
	// Auftrags-Hashwertverfahren
	//
	// Information, welches Verfahren für die Hashwertbildung über den Kunden- auftrag verwendet werden soll.
	// Es sind nur die in [HBCI] beschriebenen Verfahren und deren Parametrisierung (Initialisierungsvektor, etc.) zulässig.
	// Codierung:
	// 0: Auftrags-Hashwert nicht unterstützt
	// 1: RIPEMD-160
	// 2: SHA-1
	JobHashMethod *CodeDataElement `yaml:"JobHashMethod"`
	// FIXME: docs
	ProcessParameters *Tan2StepSubmissionProcessParametersV6 `yaml:"ProcessParameters"`
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

// Tan2StepSubmissionParametersV6 represents a slice of
// Tan2StepSubmissionParameterV6 DataElements
type Tan2StepSubmissionProcessParametersV6 struct {
	*arrayElementGroup
}

// UnmarshalHBCI unmarshals value into the Tan2StepSubmissionParameters
func (t *Tan2StepSubmissionProcessParametersV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	var param Tan2StepSubmissionProcessParameterV6
	paramElements := len(param.Elements())
	if len(elements)%paramElements != 0 {
		return fmt.Errorf("malformed marshaled value: value pairs not even: %d/%d", len(elements), paramElements)
	}
	dataElements := make([]DataElement, len(elements)/paramElements)
	for i := 0; i < len(elements); i += paramElements {
		elem := bytes.Join(elements[i:i+paramElements], []byte(":"))
		param := &Tan2StepSubmissionProcessParameterV6{}
		err := param.UnmarshalHBCI(elem)
		if err != nil {
			return err
		}
		dataElements[i/paramElements] = param
	}
	t.arrayElementGroup = newArrayElementGroup(tan2StepSubmissionProcessParameterDEG, len(dataElements), len(dataElements), dataElements)
	return nil
}

type Tan2StepSubmissionProcessParameterV6 struct {
	DataElement                            `yaml:"-"`
	SecurityFunction                       *CodeDataElement           `yaml:"SecurityFunction"`
	TanProcess                             *CodeDataElement           `yaml:"TanProcess"`
	TechnicalIDTanProcess                  *IdentificationDataElement `yaml:"TechnicalIDTanProcess"`
	ZKATanProcess                          *AlphaNumericDataElement   `yaml:"ZKATanProcess"`
	ZKATanProcessVersion                   *AlphaNumericDataElement   `yaml:"ZKATanProcessVersion"`
	TwoStepProcessName                     *AlphaNumericDataElement   `yaml:"TwoStepProcessName"`
	TwoStepProcessMaxInputValue            *NumberDataElement         `yaml:"TwoStepProcessMaxInputValue"`
	TwoStepProcessAllowedFormat            *CodeDataElement           `yaml:"TwoStepProcessAllowedFormat"`
	TwoStepProcessReturnValueText          *AlphaNumericDataElement   `yaml:"TwoStepProcessReturnValueText"`
	TwoStepProcessReturnValueTextMaxLength *NumberDataElement         `yaml:"TwoStepProcessReturnValueTextMaxLength"`
	MultiTANAllowed                        *BooleanDataElement        `yaml:"MultiTANAllowed"`
	TanTimeAndDialogReference              *CodeDataElement           `yaml:"TanTimeAndDialogReference"`
	JobCancellationAllowed                 *BooleanDataElement        `yaml:"JobCancellationAllowed"`
	SMSAccountRequired                     *CodeDataElement           `yaml:"SMSAccountRequired"`
	IssuerAccountRequired                  *CodeDataElement           `yaml:"IssuerAccountRequired"`
	ChallengeClassRequired                 *BooleanDataElement        `yaml:"ChallengeClassRequired"`
	ChallengeStructured                    *BooleanDataElement        `yaml:"ChallengeStructured"`
	InitializationMode                     *CodeDataElement           `yaml:"InitializationMode"`
	TanMediumDescriptionRequired           *CodeDataElement           `yaml:"TanMediumDescriptionRequired"`
	HHD_UCResponseRequired                 *BooleanDataElement        `yaml:"HHD_UCResponseRequired"`
	SupportedActiveTanMedia                *NumberDataElement         `yaml:"SupportedActiveTanMedia"`
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

