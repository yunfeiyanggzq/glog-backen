package vote

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noot/ring-go/cryptography"
	"golang.org/x/crypto/sha3"
)

type Voter struct {
	Index     int
	VoterSk   *ecdsa.PrivateKey `json:"voterSk"`
	WalletSk  *ecdsa.PrivateKey `json:"walletSk"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	ExtraInfo string            `json:"extraInfo"`
}

func (voter *Voter) GenKey() (sk *ecdsa.PrivateKey, err error) {
	return crypto.GenerateKey()
}

func (voter *Voter) Vote(res string, ring []*ecdsa.PublicKey) (cryptText []byte, ringSignature *cryptography.RingSign, err error) {
	cryptText, err = cryptography.Encrypt(voter.VoterSk, []byte(res))
	if err != nil {
		return nil, nil, err
	}
	msgHash := sha3.Sum256(cryptText)
	ringSignature, err = cryptography.Sign(msgHash, ring, voter.VoterSk, voter.Index)
	return cryptText, ringSignature, err
}

func (voter *Voter) PublishResult() (*cryptography.Signature, []byte, error) {
	votePrivateKeyBytes := crypto.FromECDSA(voter.VoterSk)
	signature, err := cryptography.EcdsaSign(rand.Reader, voter.WalletSk, votePrivateKeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return signature, votePrivateKeyBytes, nil
}
