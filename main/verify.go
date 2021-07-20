package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noot/ring-go/cryptography"
	util "github.com/noot/ring-go/utils"
	"github.com/noot/ring-go/vote"
	"io"
	mathRand "math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const verifyFilePlacePrefix = "/Users/bytedance/code/ring-go/verifyFilePlace"
const saveFilePlace = "/Users/bytedance/code/ring-go/saveFilePlace"
const ciNum = 3

func generateMiner() error {
	for i := 0; i < 100; i++ {
		user := &vote.User{}
		sk, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		user.Name = fmt.Sprintf("%d", i)
		user.WalletSk = sk
		user.Address = util.GetMD5(sk.X.Bytes())
		user.Balance = 100
		vote.VerifyMiner[user.Name] = user
	}
	return nil
}

func randomChooseMiner() (users []*vote.User) {
	for i := 0; i < 3; i++ {
		users = append(users, vote.VerifyMiner[fmt.Sprintf("%d", mathRand.Intn(100))])
	}
	return users
}

func generateUser(name string) *vote.User {
	walletSk, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("产生钱包密钥对错误,错误原因，err:%v\n", err)
		return nil
	}
	user := &vote.User{
		Name:     name,
		WalletSk: walletSk,
	}
	return user
}

func generateVoter(user *vote.User, topic *vote.Topic) error {
	voter := &vote.Voter{
		Name: user.Name,
	}
	voter.WalletSk = user.WalletSk
	voterSk, err := voter.GenKey()
	if err != nil {
		return err
	}
	voter.VoterSk = voterSk
	topic.AddAddressList(&voter.WalletSk.PublicKey)
	votePublicKeyBytes := crypto.FromECDSAPub(&voter.VoterSk.PublicKey)
	signature, err := cryptography.EcdsaSign(rand.Reader, voter.WalletSk, votePublicKeyBytes)
	if err != nil {
		return err
	}

	index, err := topic.Add2VoterPublicKeyList(&voter.WalletSk.PublicKey, signature, votePublicKeyBytes)
	if err != nil {
		return err
	}
	voter.Index = index
	user.Voter = voter
	return nil
}

// 获取文件的md5码
func getFileMd5(verifyFilePlace string, filename string) string {
	// 文件全路径名
	path := fmt.Sprintf("%s/%s/%s", verifyFilePlace, filename, filename)
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

func doRunCi(userName string, fileName string) (string, error) {
	verifyFilePlace := fmt.Sprintf("%s/%s", verifyFilePlacePrefix, userName)
	unzipFile := fmt.Sprintf("%s/%s.zip", saveFilePlace, fileName)
	cmd := exec.Command("/bin/sh", "./unzip.sh", unzipFile, verifyFilePlace)
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output:", err)
		return "", err
	}
	md5Val := getFileMd5(verifyFilePlace, fileName)
	cmd = exec.Command("/bin/sh", "./check.sh", fmt.Sprintf("%s/%s", verifyFilePlace, fileName), fileName, fmt.Sprintf("MD5 (%s) = %s", fileName, md5Val))
	bytes, err = cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output:", err)
		return "", err
	}
	return string(bytes), nil
}

func verifyFile(file *vote.File, user *vote.User, voteWaiter *sync.WaitGroup, caliVoteWaiter *sync.WaitGroup, caliResultWaiter *sync.WaitGroup, mainWaiter *sync.WaitGroup) {
	//./check.sh ../saveFilePlace/test.zip ../verifyFilePlace  test  MD5 (test) = 19d73de51c06285ef9bc401e7f7dc778
	defer mainWaiter.Done()
	file.CiProgress = 1
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	res, err := doRunCi(user.Name, file.Name)
	if err != nil {
		res = "运行失败//00失败00"
	}

	// 投票阶段
	cryptText, oneTimeSignature, err := user.Voter.Vote(res, file.CiVoteTopic.GetRing())
	if err != nil {
		fmt.Printf("一次性环签名签名错误 %v\n", err)
		return
	}
	success, err := file.CiVoteTopic.VerifyRingSignature(oneTimeSignature)
	if err != nil || !success {
		fmt.Println("验证一次性环签名签名错误")
		return
	} else {
		file.CiResult.VoteResultDetail[string(cryptText)] = oneTimeSignature
	}
	voteWaiter.Done()
	voteWaiter.Wait()
	// 计票阶段
	file.CiProgress = 2
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	publishSignature, votePrivateKeyBytes, _ := user.Voter.PublishResult()
	caliVoteWaiter.Done()
	caliVoteWaiter.Wait()
	voterAddress, voterContent, err := file.CiVoteTopic.CaliVoterSignature(&user.Voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		fmt.Println("计票错误")
		return
	}
	file.CiResult.CaliResultDetail[voterAddress] = voterContent
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	caliResultWaiter.Done()
	caliResultWaiter.Wait()
	file.CiProgress = 3
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
}

func main() {
	topic := &vote.Topic{
		TopicName: "test",
	}
	file := &vote.File{
		Name:        "test",
		CiVoteTopic: topic,
		CiResult: &vote.VerifyResult{
			VoteResultDetail: make(map[string]*cryptography.RingSign),
			CaliResultDetail: make(map[string]string),
		},
	}

	user1 := generateUser("user1")
	user2 := generateUser("user2")
	user3 := generateUser("user3")
	user4 := generateUser("user4")
	generateVoter(user1, file.CiVoteTopic)
	generateVoter(user2, file.CiVoteTopic)
	generateVoter(user3, file.CiVoteTopic)
	generateVoter(user4, file.CiVoteTopic)
	voteWaiter := &sync.WaitGroup{}
	voteWaiter.Add(ciNum)
	caliVoteWaiter := &sync.WaitGroup{}
	caliVoteWaiter.Add(ciNum)
	caliResultWaiter := &sync.WaitGroup{}
	caliResultWaiter.Add(ciNum)
	mainWaiter := &sync.WaitGroup{}
	mainWaiter.Add(ciNum)
	go verifyFile(file, user1, voteWaiter, caliVoteWaiter, caliResultWaiter, mainWaiter)
	go verifyFile(file, user2, voteWaiter, caliVoteWaiter, caliResultWaiter, mainWaiter)
	go verifyFile(file, user3, voteWaiter, caliVoteWaiter, caliResultWaiter, mainWaiter)
	mainWaiter.Wait()
	file.CiResult.FinalResult = "success"
	for _, v := range file.CiResult.CaliResultDetail {
		if strings.Contains(v, "00失败00") {
			file.CiResult.FinalResult = "failed"
			break
		}
	}
}
