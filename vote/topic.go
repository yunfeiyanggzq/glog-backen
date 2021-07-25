package vote

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/noot/ring-go/cryptography"
	util "github.com/noot/ring-go/utils"
	"golang.org/x/crypto/sha3"
)

type Topic struct {
	TopicName         string
	Ring              []*ecdsa.PublicKey
	ImageMap          map[string]string
	VoterAddressMap   map[*ecdsa.PublicKey]*ecdsa.PublicKey
	CollectorSigMap   map[*cryptography.RingSign][]byte
	VoteStartTime     int64 `json:"voteStartTime"`
	CaliVoteStartTime int64 `json:"caliVoteStartTime"`
}

var TopicMap = make(map[string]*Topic)

func (topic *Topic) AddAddressList(publicKey *ecdsa.PublicKey) {
	if topic.VoterAddressMap == nil {
		topic.VoterAddressMap = make(map[*ecdsa.PublicKey]*ecdsa.PublicKey, 0)
	}
	topic.VoterAddressMap[publicKey] = nil
}

func (topic *Topic) GetVoterAddressList() []*ecdsa.PublicKey {
	s := make([]*ecdsa.PublicKey, 0, len(topic.VoterAddressMap))
	for _, v := range topic.VoterAddressMap {
		s = append(s, v)
	}
	return s
}

func (topic *Topic) GetRing() []*ecdsa.PublicKey {
	return topic.Ring
}

func (topic *Topic) Add2VoterPublicKeyList(walletAddress *ecdsa.PublicKey, walletSignature *cryptography.Signature, voterPublicKeyBytes []byte) (index int, err error) {
	v, ok := topic.VoterAddressMap[walletAddress]
	if !ok {
		return -1, errors.New("该投票者不在投票人列表之中")
	}
	if v != nil {
		return -1, errors.New("该投票人已经添加过投票公钥")
	}
	res := cryptography.EcdsaVerify(walletAddress, voterPublicKeyBytes, walletSignature)
	if !res {
		return -1, errors.New("投票人签名验证不通过，无法添加到环签名列表")
	}
	voterPublicKey, err := crypto.UnmarshalPubkey(voterPublicKeyBytes)
	if err != nil {
		return -1, errors.New("解析投票公钥失败")
	}
	if topic.Ring == nil {
		topic.Ring, err = cryptography.GenNewKeyRingWithPublicKey(1, voterPublicKey, 0)
		if err != nil {
			return -1, errors.New("添加环失败")
		}
	} else {
		topic.Ring, err = cryptography.GenKeyRingWithPublicKey(topic.Ring, voterPublicKey, len(topic.Ring))
		if err != nil {
			return -1, errors.New("添加环失败")
		}
	}

	topic.VoterAddressMap[walletAddress] = voterPublicKey
	return len(topic.Ring) - 1, nil
}

func (topic *Topic) VerifyRingSignature(ringSignature *cryptography.RingSign) (success bool, err error) {
	if topic.ImageMap == nil {
		topic.ImageMap = make(map[string]string, 0)
		topic.ImageMap[ringSignature.I.X.String()] = ringSignature.I.Y.String()
	} else {
		if v, ok := topic.ImageMap[ringSignature.I.X.String()]; ok && v == ringSignature.I.Y.String() {
			return false, errors.New("多次参与投票")
		}
		topic.ImageMap[ringSignature.I.X.String()] = ringSignature.I.Y.String()
	}
	if topic.CollectorSigMap == nil {
		topic.CollectorSigMap = make(map[*cryptography.RingSign][]byte)
	}
	topic.CollectorSigMap[ringSignature] = nil
	return cryptography.Verify(ringSignature), nil
}

func (topic *Topic) CaliVoterSignature(walletPubKey *ecdsa.PublicKey, walletSignature *cryptography.Signature, voterPrivateKeyBytes []byte, cryptText []byte) (voterAddress string, voterContent string, err error) {
	res := cryptography.EcdsaVerify(walletPubKey, voterPrivateKeyBytes, walletSignature)
	if !res {
		return "", "", errors.New("投票人签名验证不通过，无法添加到环签名列表")
	}
	votePrivateKey, _ := crypto.ToECDSA(voterPrivateKeyBytes)
	x, y := crypto.S256().ScalarBaseMult(votePrivateKey.D.Bytes())
	flag := false
	for k, v := range topic.VoterAddressMap {
		if fmt.Sprintf("%v", k.X) == fmt.Sprintf("%v", walletPubKey.X) &&
			fmt.Sprintf("%v", k.Y) == fmt.Sprintf("%v", walletPubKey.Y) {
			if fmt.Sprintf("%v", v.X) == fmt.Sprintf("%v", x) &&
				fmt.Sprintf("%v", v.Y) == fmt.Sprintf("%v", y) {
				flag = true
				voterAddress = util.GetMD5(x.Bytes())
			}
		}
	}

	if !flag {
		return "", "", errors.New("该用户不在投票人列表中")
	}

	index := -1
	for i := 0; i < len(topic.Ring); i++ {
		if fmt.Sprintf("%v", x) == fmt.Sprintf("%v", topic.Ring[i].X) {
			index = i
			break
		}
	}

	if index == -1 {
		return "", "", errors.New("该用户公钥不再环中")
	}

	msgHash := sha3.Sum256(cryptText)
	ringSignature, err := cryptography.Sign(msgHash, topic.Ring, votePrivateKey, index)
	if err != nil {
		return "", "", errors.New("重构签名过程失败")
	}

	linkFlag := false
	for k, _ := range topic.CollectorSigMap {
		if cryptography.Link(ringSignature, k) {
			linkFlag = true
			if fmt.Sprintf("%v", msgHash) == fmt.Sprintf("%v", k.M) {
				decry, err := cryptography.Decry(votePrivateKey, cryptText)
				if err != nil {
					return "", "", errors.New("解密投票内容失败")
				}
				voterContent = string(decry)
				return voterAddress, voterContent, nil
			} else {
				return "", "", errors.New("投票密文被更换")
			}
		}
	}
	if !linkFlag {
		return "", "", errors.New("用户在投票阶段未投票")
	}

	return "", "", errors.New("计算投票结果失败，原因由用户引起")
}
