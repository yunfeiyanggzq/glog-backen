package vote

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
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

func GenerateVoter(name string,walletSk *ecdsa.PrivateKey, topic *Topic) (voter *Voter, err error) {
	voter = &Voter{
		Name: name,
		WalletSk: walletSk,
	}
	voterSk, err := voter.GenKey()
	if err != nil {
		fmt.Printf("产生一次性投票密钥对错误,错误原因，err:%v\n", err)
		return nil, err
	}
	voter.VoterSk = voterSk
	topic.AddAddressList(&voter.WalletSk.PublicKey)
	votePublicKeyBytes := crypto.FromECDSAPub(&voter.VoterSk.PublicKey)
	signature, err := cryptography.EcdsaSign(rand.Reader, voter.WalletSk, votePublicKeyBytes)
	if err != nil {
		fmt.Printf("发送投票公钥错误,错误原因，err:%v\n", err)
	}

	index, err := topic.Add2VoterPublicKeyList(&voter.WalletSk.PublicKey, signature, votePublicKeyBytes)
	if err != nil {
		fmt.Printf("加入环失败,错误原因，err:%v\n", err)
	}
	voter.Index = index

	return voter, nil
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
