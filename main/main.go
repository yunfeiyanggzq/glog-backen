package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	jsonIterator "github.com/json-iterator/go"
	util "github.com/noot/ring-go/utils"
	"github.com/noot/ring-go/vote"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	jsonProcessor = jsonIterator.ConfigFastest
)

type httpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")                    //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                        //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                  //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

func Login(writer http.ResponseWriter, request *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(writer).Encode(result)
	}()
	user := &vote.User{}
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}

	userInfo, ok := vote.UserMap[user.Name]
	if !ok {
		result.Code = 400
		result.Message = "user name is not exist"
		return
	}
	if userInfo.LoginPassword == user.LoginPassword {
		result.Code = 200
		result.Message = "success"
	} else {
		result.Code = 400
		result.Message = "password or user name is wrong"
	}
}

func Register(writer http.ResponseWriter, request *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(writer).Encode(result)
	}()
	user := &vote.User{}
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}

	if _, ok := vote.UserMap[user.Name]; ok {
		result.Code = 400
		result.Message = "user name is exist,please change user name"
		return
	}

	sk, err := crypto.GenerateKey()
	if err != nil {
		result.Code = 500
		result.Message = "generate user wallet private key failed"
		return
	}
	user.WalletSk = sk
	user.Address = util.GetMD5(sk.X.Bytes())
	vote.UserMap[user.Name] = user

	entity := []byte("{}")
	entity = util.SetJsonBytes(entity, "address", user.Address)
	result.Code = 200
	result.Message = "success"
	result.Object = string(entity)
}

func UploadImage(writer http.ResponseWriter, request *http.Request) {
	SaveFile(writer, request, "png")
}

func UploadFile(writer http.ResponseWriter, request *http.Request) {
	SaveFile(writer, request, "zip")
}
func UploadFileInfo(writer http.ResponseWriter, request *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(writer).Encode(result)
	}()
	file := &vote.File{
		CreateTime: time.Now().Format(time.RFC1123),
	}
	err := json.NewDecoder(request.Body).Decode(file)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}

	if _, ok := vote.FileMap[file.Name]; ok {
		result.Code = 400
		result.Message = "file name is exist,please change file name"
		return
	}

	vote.FileMap[file.Name] = file
	result.Code = 200
	result.Message = "success"

	go minerVerifyFiByCI(file)
}

func GetFileInfoList(writer http.ResponseWriter, request *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(writer).Encode(result)
	}()

	var fileInfoList []vote.FileExternal
	for _, fileInfo := range vote.FileMap {
		fileInfoCopy := vote.FileExternal{
			Name:              fileInfo.Name,
			Introduction:      fileInfo.Introduction,
			Usage:             fileInfo.Usage,
			Extra:             fileInfo.Extra,
			Install:           fileInfo.Install,
			CiProgress:        fileInfo.CiProgress,
			CloseVoteProgress: fileInfo.CloseVoteProgress,
			CreateTime:        fileInfo.CreateTime,
			ViewCount:         fileInfo.ViewCount,
		}
		fileInfoList = append(fileInfoList, fileInfoCopy)
	}

	result.Code = 200
	result.Message = "success"

	fileInfoBytes, err := json.Marshal(fileInfoList)
	if err != nil {
		result.Code = 400
		result.Message = "server go wrong"
		fmt.Println(err)
		return
	}
	fmt.Println(string(fileInfoBytes))
	result.Object = string(fileInfoBytes)
}

func GetFileInfo(writer http.ResponseWriter, request *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(writer).Encode(result)
	}()
	file := &vote.File{}
	err := json.NewDecoder(request.Body).Decode(file)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}

	fileInfo, ok := vote.FileMap[file.Name]
	if !ok {
		result.Code = 400
		result.Message = "file name is not exist,please ensure file name is right"
		return
	}
	fileInfo.ViewCount = fileInfo.ViewCount + 1

	result.Code = 200
	result.Message = "success"
	fileInfoCopy := &vote.FileExternal{
		Name:              fileInfo.Name,
		Introduction:      fileInfo.Introduction,
		Usage:             fileInfo.Usage,
		Extra:             fileInfo.Extra,
		Install:           fileInfo.Install,
		CiProgress:        fileInfo.CiProgress,
		CloseVoteProgress: fileInfo.CloseVoteProgress,
		CreateTime:        fileInfo.CreateTime,
		ViewCount:         fileInfo.ViewCount,
	}

	fileInfoBytes, err := json.Marshal(fileInfoCopy)
	if err != nil {
		result.Code = 400
		result.Message = "server go wrong"
		fmt.Println(err)
		return
	}
	fmt.Println(string(fileInfoBytes))
	result.Object = string(fileInfoBytes)
}

func SaveFile(w http.ResponseWriter, r *http.Request, fileType string) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	r.ParseMultipartForm(32 << 20)
	for fileName, _ := range r.MultipartForm.File {
		file, fileinfo_w, err := r.FormFile(fileName)
		if err != nil {
			fmt.Println("接收文件异常: ", err)
			return
		}
		filepath := fmt.Sprintf("../saveFilePlace/%s.%s", fileName, fileType)
		if fileinfo_w != nil {
			defer file.Close()
			os.Remove(filepath)
			f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Printf("%s 接收失败，创建文件失败\n", filepath)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			ResponseWithOrigin(w, r, http.StatusOK)
			ResponseWithOrigin(w, r, http.StatusOK)
			fmt.Printf("%s 接收完成\n", filepath)
		}
	}
}
func CloseVoteEnsure(w http.ResponseWriter, r *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(w).Encode(result)
	}()
	ensureVoteInfo := &vote.EnsureVote{}
	err := json.NewDecoder(r.Body).Decode(ensureVoteInfo)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}
	userInfo, ok := vote.UserMap[ensureVoteInfo.UserName]
	if !ok {
		result.Code = 400
		result.Message = "用户不存在，请注册"
		return
	}

	fileInfo, ok := vote.FileMap[ensureVoteInfo.FileName]
	if !ok {
		result.Code = 400
		result.Message = "文件不存在，请确认后文件名是否正确"
		return
	}
	if err = generateVoter(userInfo, fileInfo.CloseVoteTopic); err != nil {
		result.Code = 400
		result.Message = fmt.Sprintf("生成voter失败 err:%v", err)
		return
	}
	result.Code = 200
	result.Message = "确认投票成功"
	fmt.Println("成功添加" + userInfo.Name + "到" + fileInfo.Name + "的封闭式投票环中")
	return
}

func OpenVote(w http.ResponseWriter, r *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(w).Encode(result)
	}()
	voteInfo := &vote.VoteInfo{}
	err := json.NewDecoder(r.Body).Decode(voteInfo)
	if err != nil {
		result.Code = 400
		result.Message = fmt.Sprintf("解码错误 %v", err)
		return
	}
	userInfo, ok := vote.UserMap[voteInfo.UserName]
	if !ok {
		result.Code = 400
		result.Message = "用户不存在，请注册"
		return
	}

	fileInfo, ok := vote.FileMap[voteInfo.FileName]
	if !ok {
		result.Code = 400
		result.Message = "文件不存在，请确认后文件名是否正确"
		return
	}
	contentBytes, err := jsonProcessor.Marshal(voteInfo.VoteContent)
	if err != nil {
		result.Code = 400
		result.Message = "投票内容压缩错误"
		return
	}
	if fileInfo.CloseVoteProgress != 5 {
		result.Code = 400
		result.Message = "封闭式投票尚未结束，请继续等待其结束后进行投票"
		return
	}
	fileInfo.OpenCheckResult.CaliResultDetail[userInfo.Address] = string(contentBytes)
	fmt.Println("结果：" + userInfo.Address + "：" + string(contentBytes))
	result.Code = 200
	result.Message = "投票成功"

}
func CloseVote(w http.ResponseWriter, r *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(w).Encode(result)
	}()
	voteInfo := &vote.VoteInfo{}
	err := json.NewDecoder(r.Body).Decode(voteInfo)
	if err != nil {
		result.Code = 400
		result.Message = "decode request failed"
		return
	}
	userInfo, ok := vote.UserMap[voteInfo.UserName]
	if !ok {
		result.Code = 400
		result.Message = "用户不存在，请注册"
		return
	}

	fileInfo, ok := vote.FileMap[voteInfo.FileName]
	if !ok {
		result.Code = 400
		result.Message = "文件不存在，请确认后文件名是否正确"
		return
	}
	if time.Now().Unix() < fileInfo.CloseVoteTopic.VoteStartTime {
		result.Code = 400
		result.Message = "投票尚未开始，请稍等"
		return
	}

	if time.Now().Unix() > fileInfo.CloseVoteTopic.CaliVoteStartTime {
		result.Code = 400
		result.Message = "投票已结束"
		return
	}

	contentBytes, err := jsonProcessor.Marshal(voteInfo.VoteContent)
	if err != nil {
		result.Code = 400
		result.Message = "投票内容压缩错误"
		return
	}

	// 投票阶段
	fmt.Println(userInfo.Voter)
	cryptText, oneTimeSignature, err := userInfo.Voter.Vote(string(contentBytes), fileInfo.CloseVoteTopic.GetRing())
	if err != nil {
		result.Code = 400
		result.Message = fmt.Sprintf("一次性环签名签名错误 %v\n", err)
		return
	}
	vote.VoteContentCryptMap[userInfo.Name] = cryptText
	success, err := fileInfo.CloseVoteTopic.VerifyRingSignature(oneTimeSignature)
	if err != nil || !success {
		result.Code = 400
		result.Message = fmt.Sprintf("验证一次性环签名签名错误%v\n", err)
		return
	} else {
		fileInfo.CloseCheckResult.VoteResultDetail[string(cryptText)] = oneTimeSignature
		result.Code = 200
		result.Message = "投票成功，请耐心等待计票"
		return
	}
}

func PublishVoteSk(w http.ResponseWriter, r *http.Request) {
	var result httpResponse
	defer func() {
		_ = json.NewEncoder(w).Encode(result)
	}()
	voteInfo := &vote.VoteInfo{}
	err := json.NewDecoder(r.Body).Decode(voteInfo)
	if err != nil {
		result.Code = 400
		fmt.Println(err)
		return
	}
	userInfo, ok := vote.UserMap[voteInfo.UserName]
	if !ok {
		result.Code = 400
		result.Message = "用户不存在，请注册"
		return
	}

	fileInfo, ok := vote.FileMap[voteInfo.FileName]
	if !ok {
		result.Code = 400
		result.Message = "文件不存在，请确认后文件名是否正确"
		return
	}

	//if time.Now().Unix() < fileInfo.CloseVoteTopic.CaliVoteStartTime {
	//	result.Code = 400
	//	result.Message = "未到计票阶段，请耐心等待"
	//	return
	//}

	publishSignature, votePrivateKeyBytes, _ := userInfo.Voter.PublishResult()
	voterAddress, voterContent, err := fileInfo.CloseVoteTopic.CaliVoterSignature(&userInfo.Voter.WalletSk.PublicKey, publishSignature, votePrivateKeyBytes, vote.VoteContentCryptMap[userInfo.Name])
	if err != nil {
		result.Code = 400
		result.Message = fmt.Sprintf("计票错误。err:%v", err)
		return
	}
	fileInfo.CloseCheckResult.CaliResultDetail[voterAddress] = voterContent
	fmt.Println("结果：" + voterAddress + "：" + voterContent)
	result.Code = 200
	result.Message = "成功"
}

func ResponseWithOrigin(w http.ResponseWriter, r *http.Request, code int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)
	//w.Write(json)
}

func main() {
	err := generateMinerMock()
	if err != nil {
		log.Fatal("failed to mock miners", err)
	}
	err = generateUserMock()
	if err != nil {
		log.Fatal("failed to mock users", err)
	}
	http.HandleFunc("/publishSk", cors(PublishVoteSk))
	http.HandleFunc("/closeVote", cors(CloseVote))
	http.HandleFunc("/openVote", cors(OpenVote))
	http.HandleFunc("/ensureVote", cors(CloseVoteEnsure))
	http.HandleFunc("/register", cors(Register))
	http.HandleFunc("/getFileInfo", cors(GetFileInfo))
	http.HandleFunc("/getFileInfoList", cors(GetFileInfoList))
	http.HandleFunc("/uploadImage", cors(UploadImage))
	http.HandleFunc("/uploadFile", cors(UploadFile))
	http.HandleFunc("/login", cors(Login))
	http.HandleFunc("/uploadFileInfo", cors(UploadFileInfo))
	fmt.Println("初始化handler成功")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("err:", err)
	}
}
