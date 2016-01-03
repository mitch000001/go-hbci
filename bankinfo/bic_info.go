package bankinfo

type BicInfo struct {
	BIC    string
	BankId string
}

type SortableBicInfos []BicInfo

func (s SortableBicInfos) Len() int           { return len(s) }
func (s SortableBicInfos) Swap(a, b int)      { s[a], s[b] = s[b], s[a] }
func (s SortableBicInfos) Less(a, b int) bool { return s[a].BankId < s[b].BankId }
