package element

import "testing"

func TestTan2StepSubmissionParameterV6_UnmarshalHBCI(t *testing.T) {
	type fields struct {
		DataElement                        DataElement
		OneStepProcessAllowed              *BooleanDataElement
		MoreThanOneObligatoryTanJobAllowed *BooleanDataElement
		JobHashMethod                      *CodeDataElement
		ProcessParameters                  *Tan2StepSubmissionProcessParametersV6
	}
	tests := []struct {
		name    string
		fields  fields
		value   []byte
		wantErr bool
	}{
		{
			name:   "valid params",
			fields: fields{},
			value: []byte(
				"J:N:0:910:2:HHD1.3.0:::chipTAN manuell:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:0:N:1:" +
					"911:2:HHD1.3.2OPT:HHDOPT1:1.3.2:chipTAN optisch:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:0:N:1:" +
					"912:2:HHD1.3.2USB:HHDUSB1:1.3.2:chipTAN-USB:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:0:N:1:" +
					"913:2:Q1S:Secoder_UC:1.2.0:chipTAN-QR:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:0:N:1:" +
					"920:2:smsTAN:::smsTAN:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:2:N:5:" +
					"921:2:pushTAN:::pushTAN:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:2:N:2:" +
					"900:2:iTAN:::iTAN:6:1:TAN-Nummer:3:J:2:N:0:0:N:N:00:0:N:0",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tan2StepSubmissionParameterV6{
				DataElement:                        tt.fields.DataElement,
				OneStepProcessAllowed:              tt.fields.OneStepProcessAllowed,
				MoreThanOneObligatoryTanJobAllowed: tt.fields.MoreThanOneObligatoryTanJobAllowed,
				JobHashMethod:                      tt.fields.JobHashMethod,
				ProcessParameters:                  tt.fields.ProcessParameters,
			}
			if err := tr.UnmarshalHBCI(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Tan2StepSubmissionParameterV6.UnmarshalHBCI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
