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
	"time"
)

const verifyFilePlacePrefix = "/Users/bytedance/code/ring-go/verifyFilePlace"
const saveFilePlace = "/Users/bytedance/code/ring-go/saveFilePlace"
const ciNum = 3
const ensureWaitSecond = 100
const caliWaitSecond = 100
const finishWaitSecond = 100

func generateMinerMock() error {
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

func generateUserMock() error {
	for i := 0; i < 1000; i++ {
		user := &vote.User{}
		sk, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		user.Name = fmt.Sprintf("user%d", i)
		user.Email = fmt.Sprintf("user%d@qq.com", i)
		user.Phone = "18811228731"
		user.Image = "default"
		user.Introduction = ""
		user.LoginPassword = fmt.Sprintf("user%d", i)
		user.WalletSk = sk
		user.Address = util.GetMD5(sk.X.Bytes())
		user.Balance = 100
		vote.UserMap[user.Name] = user
	}
	return nil
}

func randomChooseMiner() (users []*vote.User) {
	for i := 0; i < 3; i++ {
		users = append(users, vote.VerifyMiner[fmt.Sprintf("%d", mathRand.Intn(100))])
	}
	return users
}

func randomChooseUser() (users []*vote.User) {
	for i := 0; i < 10; i++ {
		users = append(users, vote.UserMap[fmt.Sprintf("user%d", mathRand.Intn(1000))])
	}
	return users
}

func generateVoter(user *vote.User, topic *vote.Topic) error {
	voter := &vote.Voter{
		Name: user.Name,
	}
	if _, ok := topic.VoterAddressMap[&user.WalletSk.PublicKey]; !ok {
		return fmt.Errorf("该用户没有被选中参与投票")
	}
	voter.WalletSk = user.WalletSk
	voterSk, err := voter.GenKey()
	if err != nil {
		return err
	}
	voter.VoterSk = voterSk
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
	time.Sleep(time.Second * time.Duration(mathRand.Intn(10)))
	voteWaiter.Done()
	voteWaiter.Wait()
	// 计票阶段
	file.CiProgress = 2
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	publishSignature, votePrivateKeyBytes, _ := user.Voter.PublishResult()
	time.Sleep(time.Second * time.Duration(mathRand.Intn(10)))
	caliVoteWaiter.Done()
	caliVoteWaiter.Wait()
	voterAddress, voterContent, err := file.CiVoteTopic.CaliVoterSignature(&user.Voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, cryptText)
	if err != nil {
		fmt.Println("计票错误")
		return
	}
	file.CiProgress = 3
	file.CiResult.CaliResultDetail[voterAddress] = voterContent
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	time.Sleep(time.Second * time.Duration(mathRand.Intn(10)))
	caliResultWaiter.Done()
	caliResultWaiter.Wait()
	file.CiProgress = 4
	fmt.Printf("userName:%s progress:%d\n", user.Name, file.CiProgress)
	time.Sleep(time.Second * time.Duration(mathRand.Intn(10)))
}

func StartCloseVoteByUser(file *vote.File) {
	fmt.Println("开始封闭式投票流程！" + file.Name)
	file.CloseVoteTopic = &vote.Topic{
		TopicName: file.Name,
	}
	file.CloseCheckResult = &vote.VerifyResult{
		VoteResultDetail: make(map[string]*cryptography.RingSign),
		CaliResultDetail: make(map[string]string),
	}
	file.CloseVoteTopic.VoteStartTime = time.Now().Add(time.Second * ensureWaitSecond).Unix()
	file.CloseVoteTopic.CaliVoteStartTime = time.Now().Add(time.Second * (caliWaitSecond + ensureWaitSecond)).Unix()
	// 选择投票人
	users := randomChooseUser()
	fmt.Println(users)
	for _, user := range users {
		fmt.Println(user)
		fmt.Println(file)
		file.CloseVoteTopic.AddAddressList(&user.WalletSk.PublicKey)
		fmt.Println("被选中参与投票的用户" + user.Name)
	}
	file.CloseVoteProgress = 1
	fmt.Printf("fileName:%s close vote progress:%d\n", file.Name, file.CloseVoteProgress)
	fmt.Println("请确认参与投票")
	time.Sleep(time.Second * ensureWaitSecond)
	file.CloseVoteProgress = 2
	fmt.Printf("fileName:%s close vote progress:%d\n", file.Name, file.CloseVoteProgress)
	fmt.Println("请投票")
	time.Sleep(time.Second * caliWaitSecond)
	file.CloseVoteProgress = 3
	fmt.Printf("fileName:%s close vote progress:%d\n", file.Name, file.CloseVoteProgress)
	fmt.Println("开始计票")
	time.Sleep(time.Second * finishWaitSecond)
	file.CloseVoteProgress = 4
	fmt.Printf("fileName:%s close vote progress:%d\n", file.Name, file.CloseVoteProgress)
	fmt.Println("投票结束")

	// 如果通过
	file.CloseVoteProgress = 5
	fmt.Printf("fileName:%s close vote progress:%d\n", file.Name, file.CloseVoteProgress)
	// 开始开放式投票
	file.OpenCheckResult = &vote.VerifyResult{
		CaliResultDetail: make(map[string]string),
	}
}

func minerVerifyFiByCI(file *vote.File) {
	file.CiVoteTopic = &vote.Topic{
		TopicName: file.Name,
	}
	file.CiResult = &vote.VerifyResult{
		VoteResultDetail: make(map[string]*cryptography.RingSign),
		CaliResultDetail: make(map[string]string),
	}
	fmt.Println("开始CI验证" + file.Name)
	miners := randomChooseMiner()
	for _, miner := range miners {
		file.CiVoteTopic.AddAddressList(&miner.WalletSk.PublicKey)
		generateVoter(miner, file.CiVoteTopic)
	}
	voteWaiter := &sync.WaitGroup{}
	voteWaiter.Add(len(miners))
	caliVoteWaiter := &sync.WaitGroup{}
	caliVoteWaiter.Add(len(miners))
	caliResultWaiter := &sync.WaitGroup{}
	caliResultWaiter.Add(len(miners))
	mainWaiter := &sync.WaitGroup{}
	mainWaiter.Add(len(miners))
	for _, miner := range miners {
		go verifyFile(file, miner, voteWaiter, caliVoteWaiter, caliResultWaiter, mainWaiter)
	}
	mainWaiter.Wait()
	file.CiResult.FinalResult = "success"
	for _, v := range file.CiResult.CaliResultDetail {
		if strings.Contains(v, "00失败00") {
			file.CiResult.FinalResult = "failed"
			fmt.Println("自动化测试验证结果为失败！，终止验证流程")
			fmt.Println(v)
			break
		}
	}
	file.CloseVoteProgress = 5
	fmt.Println(file.Name + "验证结束,开启封闭式投票流程")
	go StartCloseVoteByUser(file)
}
