package cryptography

import (
	"crypto/ecdsa"
	"io"
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func EcdsaSign(rand io.Reader, priv *ecdsa.PrivateKey, hash []byte) (signature *Signature, err error) {
	r, s, err := ecdsa.Sign(rand, priv, hash)
	if err != nil {
		return nil, err
	}
	signature = &Signature{
		r: r,
		s: s,
	}
	return signature, nil
}

func EcdsaVerify(pub *ecdsa.PublicKey, hash []byte, signature *Signature) bool {
	return ecdsa.Verify(pub, hash, signature.r, signature.s)
}
