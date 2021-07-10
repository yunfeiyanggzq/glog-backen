package server

import (
	"crypto/ecdsa"
	"github.com/noot/ring-go/vote"
)

type User struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Phone      string            `json:"phone"`
	Bobby      string            `json:"hobby"`
	Reputation int               `json:"reputation"`
	Motto      string            `json:"motto"`
	WalletSk   *ecdsa.PrivateKey `json:"walletSk"`
	Voter      *vote.Voter       `json:"voter"`
	Address    string            `json:"address"`
}
