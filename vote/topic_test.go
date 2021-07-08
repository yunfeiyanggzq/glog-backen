package vote

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noot/ring-go/cryptography"
	"io/ioutil"
	"testing"
)

func generateVoter(name string, topic *Topic) (voter *Voter, err error) {
	voter = &Voter{
		Name: name,
	}
	walletSk, err := voter.GenKey()
	if err != nil {
		fmt.Printf("产生钱包密钥对错误,错误原因，err:%v\n", err)
		return nil, err
	}
	voter.WalletSk = walletSk
	voterSk, err := voter.GenKey()
	if err != nil {
		fmt.Printf("产生一次性投票密钥对错误,错误原因，err:%v\n", err)
		return nil, err
	}
	voter.VoterSk = voterSk
	topic.AddAddressList(&voter.WalletSk.PublicKey)
	votePublicKeyBytes := crypto.FromECDSAPub(&voter.VoterSk.PublicKey)
	signature, err := cryptography.EcdsaSign(rand.Reader, voter.WalletSk, votePublicKeyBytes)
	if err != nil {
		fmt.Printf("发送投票公钥错误,错误原因，err:%v\n", err)
	}

	index, err := topic.Add2VoterPublicKeyList(&voter.WalletSk.PublicKey, signature, votePublicKeyBytes)
	if err != nil {
		fmt.Printf("加入环失败,错误原因，err:%v\n", err)
	}
	voter.Index = index

	return voter, nil
}

func TestVote(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	_, _ = generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	// 投票阶段
	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		t.Error("计票错误")
	}
}

func TestVoteFromFile(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	_, _ = generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	file, err := ioutil.ReadFile("./message.txt")
	// 投票阶段

	cryptText, oneTimeSignature, err := voter.Vote(string(file), topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		t.Error("计票错误")
	}
}

func TestVoteAgain(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	_, _ = generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	// 投票阶段
	_, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}

	_, oneTimeSignature1, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err = topic.VerifyRingSignature(oneTimeSignature1)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误，符合预期")
	}
}

func TestVoteWithoutVote(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	voter3, _ := generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	// 投票阶段
	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter3.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		t.Error("计票错误,符合预期")
	}
}

func TestVoteWrongCryptText(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	voter3, _ := generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	// 投票阶段
	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	cryptText3, oneTimeSignature3, err := voter3.Vote("I will approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}
	success, err = topic.VerifyRingSignature(oneTimeSignature3)
	if err != nil||!success {
		t.Error("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		t.Error("计票错误")
	}

	// 计票阶段
	publishSignature3, votePrivateKeyBytes3, _ := voter3.PublishResult()
	err = topic.CaliVoterSignature(&voter3.WalletSk.PublicKey, publishSignature3, votePrivateKeyBytes3, cryptText3)
	if err != nil {
		t.Error("计票错误")
	}

	err = topic.CaliVoterSignature(&voter3.WalletSk.PublicKey, publishSignature3, votePrivateKeyBytes3, cryptText)
	if err != nil {
		t.Error("计票错误，符合预期")
	}
}