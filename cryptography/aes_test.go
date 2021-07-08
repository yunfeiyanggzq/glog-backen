package cryptography

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestAES(t *testing.T) {
	input := "ewe23efhwjfhjwefjcwecwheijhjehwjchjw"
	priv1, _ := crypto.GenerateKey()
	encoded, _ := Encrypt(priv1, []byte(input))
	decoded, _ := Decry(priv1, encoded)
	if input != string(decoded) {
		t.Error("aes加解密算法错误")
	}
}

func TestAESError(t *testing.T) {
	input := "ewe23efhwjfhjwefjcwecwheijhjehwjchjw"
	priv1, _ := crypto.GenerateKey()
	encoded, _ := Encrypt(priv1, []byte(input))
	priv2, _ := crypto.GenerateKey()
	decoded, err := Decry(priv2, encoded)
	if err == nil {
		t.Error("aes加解密算法错误")
	} else {
		fmt.Println("私钥错误，解密失败,符合预期")
	}
	fmt.Println(string(decoded))
}
