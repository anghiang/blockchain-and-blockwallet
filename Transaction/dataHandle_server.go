package Transaction

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

type AmountResponse struct {
	Amount float64 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float64 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}

type TransactionFrontSign struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      *string `json:"value"`
}

func (tr *TransactionFrontSign) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil || strings.TrimSpace(*tr.RecipientBlockchainAddress) == "" ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil || len(*tr.Value) == 0 {
		return false
	}
	return true
}

type TransactionBankSign struct {
	SenderBlockchainAddress    *string            `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string            `json:"recipient_blockchain_address"`
	SenderPublicKey            *string            `json:"sender_public_key"`
	Value                      *big.Int           `json:"value"`
	SignTx                     *types.Transaction `json:"sign_tx"`
}

func (tr *TransactionBankSign) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.SignTx == nil {
		return false
	}
	return true
}

type TransRecord struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Value *big.Int `json:"value"`
	Hash  string   `json:"hash"`
}
