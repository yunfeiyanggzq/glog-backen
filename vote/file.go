package vote

import "github.com/noot/ring-go/cryptography"

type FileExternal struct {
	Name              string `json:"name"`
	Introduction      string `json:"introduction"`
	Usage             string `json:"usage"`
	Extra             string `json:"extra"`
	CiProgress        int    `json:"ciProgress"`
	CloseVoteProgress int    `json:"closeVoteProgress"`
	Install           string `json:"install"`
	CreateTime        string `json:"createTime"`
	ViewCount         int    `json:"viewCount"`
	FinalResult       int    `json:"finalResult"`
}

type File struct {
	Name              string        `json:"name"`
	Introduction      string        `json:"introduction"`
	Usage             string        `json:"usage"`
	Extra             string        `json:"extra"`
	CiResult          *VerifyResult `json:"ciResult"`
	CloseCheckResult  *VerifyResult `json:"closeResult"`
	OpenCheckResult   *VerifyResult `json:"openResult"`
	FinalResult       int           `json:"finalResult"`
	CiVoteTopic       *Topic        `json:"ciVoteTopic"`
	CloseVoteTopic    *Topic        `json:"closeVoteTopic"`
	CiProgress        int           `json:"ciProgress"`
	CloseVoteProgress int           `json:"closeVoteProgress"`
	Install           string        `json:"install"`
	CreateTime        string        `json:"createTime"`
	ViewCount         int           `json:"viewCount"`
}

var FileMap = make(map[string]*File)

type VerifyResult struct {
	FinalResult      string                            `json:"finalResult"`
	VoteResultDetail map[string]*cryptography.RingSign `json:"voteResultDetail"`
	CaliResultDetail map[string]string                 `json:"caliVoteResultDetail"`
}
