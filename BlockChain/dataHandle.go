package BlockChain

import (
	"BlockWallet/Block"
	"BlockWallet/Transaction"
	"encoding/json"
	"math/big"
	"sync"
)

func bytesToBigInt(b [32]byte) *big.Int {
	bytes := b[:]
	result := new(big.Int).SetBytes(bytes)
	return result
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TransactionPool   []*Transaction.Transaction `json:"transaction_pool"`
		BlockchainAddress string                     `json:"blockchain_address"`
		Chain             []*Block.Block             `json:"chain"`
		Coinbase          string                     `json:"coinbase"`
		Port              uint16                     `json:"port"`
		Mux               sync.Mutex                 `json:"mux"`
		Neighbors         []string                   `json:"neighbors"`
		MuxNeighbors      sync.Mutex                 `json:"mux_neighbors"`
		MiningReward      *big.Int                   `json:"mining_reward"`
		Difficult         int                        `json:"difficult"`
	}{
		TransactionPool:   bc.transactionPool,
		BlockchainAddress: bc.blockchainAddress,
		Chain:             nil,
		Coinbase:          bc.coinbase,
		Port:              bc.port,
		Mux:               bc.mux,
		Neighbors:         bc.neighbors,
		MuxNeighbors:      bc.muxNeighbors,
		MiningReward:      bc.MiningReward,
		Difficult:         bc.difficult,
	})
}

func (bc *Blockchain) UnmarshalJSON(blockchainByte []byte) error {
	type TmpBlockChain struct {
		TransactionPool   []*Transaction.Transaction `json:"transaction_pool"`
		BlockchainAddress string                     `json:"blockchain_address"`
		Chain             []*Block.Block             `json:"chain"`
		Coinbase          string                     `json:"coinbase"`
		Port              uint16                     `json:"port"`
		Mux               sync.Mutex                 `json:"mux"`
		Neighbors         []string                   `json:"neighbors"`
		MuxNeighbors      sync.Mutex                 `json:"mux_neighbors"`
		MiningReward      *big.Int                   `json:"mining_reward"`
		Difficult         int                        `json:"difficult"`
	}

	var tmpbc TmpBlockChain
	err := json.Unmarshal(blockchainByte, &tmpbc)
	if err != nil {
		return err
	}
	bc.transactionPool = tmpbc.TransactionPool
	bc.chain = tmpbc.Chain
	bc.blockchainAddress = tmpbc.BlockchainAddress
	bc.coinbase = tmpbc.Coinbase
	bc.port = tmpbc.Port
	bc.mux = tmpbc.Mux
	bc.neighbors = tmpbc.Neighbors
	bc.muxNeighbors = tmpbc.MuxNeighbors
	bc.MiningReward = tmpbc.MiningReward
	bc.difficult = tmpbc.Difficult

	return nil
}
