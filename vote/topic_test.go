package vote

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestVote(t *testing.T) {
	topic := &Topic{
		TopicName: "hyperleder",
	}

	walletSk, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("产生钱包密钥对错误,错误原因，err:%v\n", err)
		t.Error("产生钱包密钥对错误")
	}
	voter, _ := GenerateVoter("voter1", walletSk, topic)
	walletSk1, _ := crypto.GenerateKey()
	_, _ = GenerateVoter("voter2", walletSk1, topic)
	walletSk2, _ := crypto.GenerateKey()
	_, _ = GenerateVoter("voter3", walletSk2, topic)
	walletSk3, _ := crypto.GenerateKey()
	_, _ = GenerateVoter("voter4", walletSk3, topic)
	// 投票阶段
	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
	if err != nil {
		t.Error("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil || !success {
		t.Error("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		t.Error("计票错误")
	}
}

//func TestVoteFromFile(t *testing.T) {
//	topic := &Topic{
//		TopicName: "hyperleder",
//	}
//
//	voter, _ := GenerateVoter("voter1", topic)
//	_, _ = GenerateVoter("voter2", topic)
//	_, _ = GenerateVoter("voter3", topic)
//	_, _ = GenerateVoter("voter4", topic)
//	file, err := ioutil.ReadFile("./message.txt")
//	// 投票阶段
//
//	cryptText, oneTimeSignature, err := voter.Vote(string(file), topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	success, err := topic.VerifyRingSignature(oneTimeSignature)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误")
//	}
//	// 计票阶段
//	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
//	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
//	if err != nil {
//		t.Error("计票错误")
//	}
//}
//
//func TestVoteAgain(t *testing.T) {
//	topic := &Topic{
//		TopicName: "hyperleder",
//	}
//
//	voter, _ := GenerateVoter("voter1", topic)
//	_, _ = GenerateVoter("voter2", topic)
//	_, _ = GenerateVoter("voter3", topic)
//	_, _ = GenerateVoter("voter4", topic)
//	// 投票阶段
//	_, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	success, err := topic.VerifyRingSignature(oneTimeSignature)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误")
//	}
//
//	_, oneTimeSignature1, err := voter.Vote("I will not approve it", topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	success, err = topic.VerifyRingSignature(oneTimeSignature1)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误，符合预期")
//	}
//}
//
//func TestVoteWithoutVote(t *testing.T) {
//	topic := &Topic{
//		TopicName: "hyperleder",
//	}
//
//	voter, _ := GenerateVoter("voter1", topic)
//	_, _ = GenerateVoter("voter2", topic)
//	voter3, _ := GenerateVoter("voter3", topic)
//	_, _ = GenerateVoter("voter4", topic)
//	// 投票阶段
//	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	success, err := topic.VerifyRingSignature(oneTimeSignature)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误")
//	}
//	// 计票阶段
//	publishSignature, votePrivateKeyBytes, _ := voter3.PublishResult()
//	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
//	if err != nil {
//		t.Error("计票错误,符合预期")
//	}
//}
//
//func TestVoteWrongCryptText(t *testing.T) {
//	topic := &Topic{
//		TopicName: "hyperleder",
//	}
//
//	voter, _ := GenerateVoter("voter1", topic)
//	_, _ = GenerateVoter("voter2", topic)
//	voter3, _ := GenerateVoter("voter3", topic)
//	_, _ = GenerateVoter("voter4", topic)
//	// 投票阶段
//	cryptText, oneTimeSignature, err := voter.Vote("I will not approve it", topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	cryptText3, oneTimeSignature3, err := voter3.Vote("I will approve it", topic.GetRing())
//	if err != nil {
//		t.Error("一次性环签名签名错误")
//	}
//	success, err := topic.VerifyRingSignature(oneTimeSignature)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误")
//	}
//	success, err = topic.VerifyRingSignature(oneTimeSignature3)
//	if err != nil || !success {
//		t.Error("验证一次性环签名签名错误")
//	}
//	// 计票阶段
//	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
//	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
//	if err != nil {
//		t.Error("计票错误")
//	}
//
//	// 计票阶段
//	publishSignature3, votePrivateKeyBytes3, _ := voter3.PublishResult()
//	err = topic.CaliVoterSignature(&voter3.WalletSk.PublicKey, publishSignature3, votePrivateKeyBytes3, cryptText3)
//	if err != nil {
//		t.Error("计票错误")
//	}
//
//	err = topic.CaliVoterSignature(&voter3.WalletSk.PublicKey, publishSignature3, votePrivateKeyBytes3, cryptText)
//	if err != nil {
//		t.Error("计票错误，符合预期")
//	}
//}
