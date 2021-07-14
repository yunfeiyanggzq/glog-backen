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
	file := &vote.File{}
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

	result.Code = 200
	result.Message = "success"
	fileInfoBytes,err:= jsonProcessor.Marshal(fileInfo)
	if err != nil {
		result.Code = 400
		result.Message = "server go wrong"
		return
	}
	fmt.Println(string(fileInfoBytes))
	result.Object=string(fileInfoBytes)
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
		filepath := fmt.Sprintf("./%s.%s", fileName, fileType)
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

func ResponseWithOrigin(w http.ResponseWriter, r *http.Request, code int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)
	//w.Write(json)
}

func main() {
	var err error
	http.HandleFunc("/register", cors(Register))
	http.HandleFunc("/getFileInfo", cors(GetFileInfo))
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
