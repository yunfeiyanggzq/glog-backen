package vote

type VoteInfo struct {
	UserName    string   `json:"userName"`
	FileName    string   `json:"fileName"`
	VoteContent *Content `json:"content"`
}

var VoteContentCryptMap = make(map[string][]byte)

type EnsureVote struct {
	UserName string `json:"userName"`
	FileName string `json:"fileName"`
}

type Content struct {
	Comment string `json:"comment"`
	Score   int `json:"score"`
}
