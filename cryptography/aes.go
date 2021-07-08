package cryptography

import (
	"crypto/ecdsa"
	"github.com/wumansgy/goEncrypt"
)

func Encrypt(sk *ecdsa.PrivateKey, msg []byte) (cryptText []byte, err error) {
	cryptText, err = goEncrypt.AesCbcEncrypt(msg, []byte(sk.D.String()[0:16]))
	if err != nil {
		return nil, err
	}
	return cryptText, nil
}

func Decry(sk *ecdsa.PrivateKey, cryptText []byte) (newplaintext []byte, err error) {
	newplaintext, err = goEncrypt.AesCbcDecrypt(cryptText, []byte(sk.D.String()[0:16]))
	if err != nil {
		return nil, err
	}
	return newplaintext, nil
}
