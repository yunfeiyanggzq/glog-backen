package vote

type File struct {
	Name              string        `json:"name"`
	Introduction      string        `json:"introduction"`
	Usage             string        `json:"usage"`
	Extra             string        `json:"extra"`
	CiResult          *VerifyResult `json:"ciResult"`
	CloseCheckResult  *VerifyResult `json:"closeResult"`
	OpenCheckResult   *VerifyResult `json:"openResult"`
	VoteTopic         Topic         `json:"voteTopic"`
	CiProgress        int           `json:"ciProgress"`
	CloseVoteProgress int           `json:"closeVoteProgress"`
	Install           string        `json:"install"`
}

var FileMap = make(map[string]*File)

type VerifyResult struct {
	FinalResult        bool              `json:"finalResult"`
	VerifyResultDetail map[string]string `json:"finalResultDetail"`
}
