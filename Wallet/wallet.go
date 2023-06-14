package Wallet

import (
	utils "BlockWallet/Utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Wallet struct {
	privateKey   *ecdsa.PrivateKey
	publicKey    *ecdsa.PublicKey
	blockAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := crypto.GenerateKey()

	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey

	w.blockAddress = crypto.PubkeyToAddress(*w.publicKey).Hex()
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

func LoadWallet(privkey string) *Wallet {

	theWallet := new(Wallet)
	thepriKey, err := crypto.HexToECDSA(privkey)
	if err != nil {
		panic(err)
	}
	theWallet.privateKey = thepriKey
	theWallet.publicKey = &thepriKey.PublicKey
	theWallet.blockAddress = crypto.PubkeyToAddress(*theWallet.publicKey).Hex()

	return theWallet
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      *big.Int
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string   `json:"sender_blockchain_address"`
		Recipient string   `json:"recipient_blockchain_address"`
		Value     *big.Int `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value *big.Int) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}
