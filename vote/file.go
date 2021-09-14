package vote

import "github.com/noot/ring-go/cryptography"

type FileExternal struct {
	Name                        string            `json:"name"`
	Introduction                string            `json:"introduction"`
	Usage                       string            `json:"usage"`
	Extra                       string            `json:"extra"`
	CiProgress                  int               `json:"ciProgress"`
	CloseVoteProgress           int               `json:"closeVoteProgress"`
	Install                     string            `json:"install"`
	CreateTime                  string            `json:"createTime"`
	ViewCount                   int               `json:"viewCount"`
	FinalResult                 int               `json:"finalResult"`
	OwnerUserName               string            `json:"ownerUserName"`
	CiVoteUserNameList          []string          `json:"ciVoteUserNameList"`
	CiVoteCommentList           map[string]string `json:"ciVoteCommentList"`
	CiVoteScoreList             map[string]string `json:"ciVoteScoreList"`
	CloseVoteUserNameList       []string          `json:"closeVoteUserNameList"`
	CloseVoteRandomUserNameList []string          `json:"closeVoteRandomUserNameList"`
	CommentList                 []string          `json:"commentList"`
	CloseVoteCommentList        map[string]string `json:"closeVoteCommentList"`
	CloseVoteScoreList          map[string]string `json:"closeVoteScoreList"`
	OpenVoteCommentList         map[string]string `json:"openVoteCommentList"`
	OpenVoteScoreList           map[string]string `json:"openVoteScoreList"`
	Value                       int               `json:"value"`
}

type File struct {
	Name                        string        `json:"name"`
	OwnerUserName               string        `json:"ownerUserName"`
	Introduction                string        `json:"introduction"`
	Usage                       string        `json:"usage"`
	Extra                       string        `json:"extra"`
	CiResult                    *VerifyResult `json:"ciResult"`
	CloseCheckResult            *VerifyResult `json:"closeResult"`
	OpenCheckResult             *VerifyResult `json:"openResult"`
	FinalResult                 int           `json:"finalResult"` // -1 failed 0 wip 1 success
	CiVoteTopic                 *Topic        `json:"ciVoteTopic"`
	CloseVoteTopic              *Topic        `json:"closeVoteTopic"`
	CiProgress                  int           `json:"ciProgress"`
	CloseVoteProgress           int           `json:"closeVoteProgress"`
	Install                     string        `json:"install"`
	CreateTime                  string        `json:"createTime"`
	ViewCount                   int           `json:"viewCount"`
	CiVoteUserNameList          []string      `json:"ciVoteUserNameList"`
	CloseVoteUserNameList       []string      `json:"closeVoteUserNameList"`
	CommentList                 []string      `json:"commentList"`
	CloseVoteRandomUserNameList []string      `json:"closeVoteRandomUserNameList"`
	Value                       int           `json:"value"`
}

var FileMap = make(map[string]*File)

type VerifyResult struct {
	FinalResult      int64                             `json:"finalResult"`
	VoteResultDetail map[string]*cryptography.RingSign `json:"voteResultDetail"`
	CaliResultDetail map[string]string                 `json:"caliVoteResultDetail"`
}
