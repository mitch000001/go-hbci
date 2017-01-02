package bankinfo

// BicInfo holds the BIC data associated with a given bank institute.
type BicInfo struct {
	BIC    string
	BankID string
}

type sortableBicInfos []BicInfo

func (s sortableBicInfos) Len() int           { return len(s) }
func (s sortableBicInfos) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s sortableBicInfos) Less(a, b int) bool { return s[a].BankID < s[b].BankID }
