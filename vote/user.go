package vote

import (
	"crypto/ecdsa"
)

type User struct {
	Name          string            `json:"name"`
	Email         string            `json:"email"`
	Phone         string            `json:"phone"`
	Bobby         string            `json:"hobby"`
	Reputation    int               `json:"reputation"`
	WalletSk      *ecdsa.PrivateKey `json:"walletSk"`
	Voter         *Voter            `json:"voter"`
	Address       string            `json:"address"`
	LoginPassword string            `json:"password"`
	Image         interface{}       `json:"image"`
	Introduction  string            `json:"introduction"`
	FileMap       map[string]*File  `json:"fileMap"`
}

var UserMap = make(map[string]*User)
