package Block

import (
	"BlockWallet/Transaction"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"
	"time"
)

type Block struct {
	nonce        *big.Int                   `json:"nonce"`
	previousHash [32]byte                   `json:"previous_hash"`
	timestamp    uint64                     `json:"timestamp"`
	number       *big.Int                   `json:"number"`
	difficulty   int                        `json:"difficulty"`
	hash         [32]byte                   `json:"hash"`
	txSize       uint16                     `json:"tx_size"`
	transactions []*Transaction.Transaction `json:"transactions"`
}

func (b *Block) Hash() [32]byte {
	return b.hash
}

func (b *Block) CalHash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Nonce() *big.Int {
	return b.nonce
}
func (b *Block) Number() *big.Int {
	return b.number
}

func (b *Block) Transactions() []*Transaction.Transaction {
	return b.transactions
}

func (b *Block) Timestamp() uint64 {
	return b.timestamp
}

func (b *Block) SetNonce(nonce *big.Int) {
	b.nonce = nonce
}

func (b *Block) SetPreviousHash(previousHash [32]byte) {
	b.previousHash = previousHash
}

func (b *Block) SetTransactions(transactions []*Transaction.Transaction) {
	b.transactions = transactions
}

func (b *Block) SetTimestamp(timestamp uint64) {
	b.timestamp = timestamp
}

func (b *Block) SetNumber(number *big.Int) {
	b.number = number
}
func (b *Block) SetDifficulty(difficulty int) {
	b.difficulty = difficulty
}

func (b *Block) SetHash(hash [32]byte) {
	b.hash = hash
}
func NewBlock(nonce *big.Int, previousHash [32]byte, txs []*Transaction.Transaction) *Block {
	b := new(Block)
	b.nonce = nonce
	b.transactions = txs
	b.previousHash = previousHash
	b.timestamp = uint64(time.Now().Unix())
	b.txSize = uint16(len(b.transactions))
	return b
}

func (b *Block) Print() {
	log.Printf("%-15v:%30d\n", "timestamp", b.timestamp)
	//fmt.Printf("timestamp       %d\n", b.timestamp)
	log.Printf("%-15v:%30d\n", "nonce", b.nonce)
	log.Printf("%-15v:%30d\n", "number", b.number)
	log.Printf("%-15v:%30x\n", "previous_hash", b.previousHash)
	log.Printf("%-15v:%30x\n", "hash", b.hash)
	//log.Printf("%-15v:%30s\n", "transactions", b.transactions)

}
