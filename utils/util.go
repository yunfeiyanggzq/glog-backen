package util

import (
	"bytes"
	"code.byted.org/gopkg/logs"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"github.com/tidwall/sjson"
)

func SetJsonBytes(bytes []byte, path string, value interface{}) (result []byte) {
	var err error
	result, err = sjson.SetBytes(bytes, path, value)
	if err != nil {
		logs.Error("set json bytes error:%s", err.Error())
		return bytes
	}
	return result
}

func DeleteJsonBytes(bytes []byte, path string) []byte {
	bytes, err := sjson.DeleteBytes(bytes, path)
	if err != nil {
		logs.Error("delete json bytes error:%s", err.Error())
		return bytes
	}
	return bytes
}

func GetMD5(b[]byte) (result string) {
	//1、创建Hash接口
	myHash:=md5.New()  //返回 Hash interface
	//2、添加数据
	myHash.Write(b)  //写入数据
	//3、计算结果
	/*
	  执行原理为：myHash.Write(b1)写入的数据进行hash运算  +  myHash.Sum(b2)写入的数据进行hash运算
	              结果为：两个hash运算结果的拼接。若myHash.Write()省略或myHash.Write(nil) ，则默认为写入的数据为“”。
	              根据以上原理，一般不采用两个hash运算的拼接，所以参数为nil
	*/
	res:=myHash.Sum(nil)  //进行运算
	//4、数据格式化
	result=hex.EncodeToString(res) //转换为string
	return
}

func DeepCopyByGob(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&buffer).Decode(dst)
}