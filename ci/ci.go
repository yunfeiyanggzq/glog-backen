package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noot/ring-go/cryptography"
	"github.com/noot/ring-go/vote"
	"io"
	"os"
	"os/exec"
	"strings"
)

func generateVoter(name string, topic *vote.Topic) (voter *vote.Voter, err error) {
	voter = &vote.Voter{
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

// 获取文件的md5码
func getFileMd5(filename string) string {
	// 文件全路径名
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		_ = fmt.Errorf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	_, _ = io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}

func doRunCi(fileName string, path string) (string, error) {
	md5Val := getFileMd5(fileName)
	cmd := exec.Command("/bin/sh", "./check.sh", path, fileName, fmt.Sprintf("MD5 (%s) = %s", fileName, md5Val))
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output:", err)
		return "", err
	}
	return string(bytes), nil
}


func main() {
	res, err := doRunCi("ci", "./")
	if err != nil || strings.Contains(res, "00失败00") {
		fmt.Println("验证结果为失败")
	} else {
		fmt.Println("验证结果为成功")
	}


	topic := &vote.Topic{
		TopicName: "hyperleder",
	}

	voter, _ := generateVoter("voter1", topic)
	_, _ = generateVoter("voter2", topic)
	_, _ = generateVoter("voter3", topic)
	_, _ = generateVoter("voter4", topic)
	// 投票阶段
	cryptText, oneTimeSignature, err := voter.Vote(res, topic.GetRing())
	if err != nil {
		fmt.Println("一次性环签名签名错误")
	}
	success, err := topic.VerifyRingSignature(oneTimeSignature)
	if err != nil || !success {
		fmt.Println("验证一次性环签名签名错误")
	}
	// 计票阶段
	publishSignature, votePrivateKeyBytes, _ := voter.PublishResult()
	err = topic.CaliVoterSignature(&voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		fmt.Println("计票错误")
	}
}
