package Block

import (
	"BlockWallet/Transaction"
	"encoding/json"
	"math/big"
)

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    uint64                     `json:"timestamp"`
		Nonce        *big.Int                   `json:"nonce"`
		PreviousHash [32]byte                   `json:"previous_hash"`
		Transactions []*Transaction.Transaction `json:"transactions"`
		Number       *big.Int                   `json:"number"`
		Difficulty   int                        `json:"difficulty"`
		Hash         [32]byte                   `json:"hash"`
		TxSize       uint16                     `json:"tx_size"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
		Number:       b.number,
		Difficulty:   b.difficulty,
		Hash:         b.hash,
		TxSize:       b.txSize,
	})
}

func (b *Block) UnmarshalJSON(blockByte []byte) error {
	type TmpBlock struct {
		Nonce        *big.Int                   `json:"nonce"`
		PreviousHash [32]byte                   `json:"previous_hash"`
		Timestamp    uint64                     `json:"timestamp"`
		Number       *big.Int                   `json:"number"`
		Difficulty   int                        `json:"difficulty"`
		Hash         [32]byte                   `json:"hash"`
		TxSize       uint16                     `json:"tx_size"`
		Transactions []*Transaction.Transaction `json:"transactions"`
	}
	var tmpb TmpBlock
	err := json.Unmarshal(blockByte, &tmpb)
	if err != nil {
		return err
	}
	b.nonce = tmpb.Nonce
	b.previousHash = tmpb.PreviousHash
	b.timestamp = tmpb.Timestamp
	b.number = tmpb.Number
	b.difficulty = tmpb.Difficulty
	b.hash = tmpb.Hash
	b.txSize = tmpb.TxSize
	b.transactions = tmpb.Transactions
	return nil
}
