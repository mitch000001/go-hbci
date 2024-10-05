package domain

import "time"

type TanProcessParameter struct {
	Steps                 int
	SecurityFunction      string
	TanProcess            string
	TanProcessTechnicalID string
	TanProcessName        string
	TanProcessVersion     string
	TwoStepProcessName    string
}

type TanParams struct {
	TANProcess           string
	JobHash              []byte
	JobReference         string
	Challenge            string
	ChallengeHHD_UC      []byte
	TANMediumDescription string
	ChallengeExpiryDate  time.Time
}
