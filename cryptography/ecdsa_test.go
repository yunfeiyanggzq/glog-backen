package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestEcdsa(t *testing.T) {
	p256 := elliptic.P256()
	priv, _ := ecdsa.GenerateKey(p256, rand.Reader)
	hashed := []byte("testing")
	signature, err := EcdsaSign(rand.Reader, priv, hashed)
	if err != nil {
		t.Error("签名错误")
	}
	res := EcdsaVerify(&priv.PublicKey, hashed, signature)
	if !res {
		t.Error("验证错误")
	}
}
