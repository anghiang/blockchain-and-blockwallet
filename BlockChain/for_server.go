package BlockChain

import (
	"BlockWallet/Block"
	"BlockWallet/Transaction"
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"math/big"
	"time"
)

func (bc *Blockchain) Run() {
	bc.StartMining()
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	// 使用time.AfterFunc函数创建了一个定时器，它在指定的时间间隔后执行bc.StartMining函数（自己调用自己）。
	_ = time.AfterFunc(time.Second*5, bc.StartMining)
	color.Yellow("minetime: %v\n", time.Now())
	fmt.Println("")

}

func (bc *Blockchain) GetTransactions() []*Transaction.Transaction {
	var tmpTransactions []*Transaction.Transaction
	for _, block := range bc.chain {
		for _, trans := range block.Transactions() {
			if trans.From() == bc.blockchainAddress {
				continue
			}

			tmpTransactions = append(tmpTransactions, trans)
		}
	}
	return tmpTransactions
}

func (bc *Blockchain) GetTransactionRecord() []*Transaction.TransRecord {
	var tmpTransactions []*Transaction.TransRecord
	for _, block := range bc.chain {
		for _, trans := range block.Transactions() {
			if trans.From() == bc.blockchainAddress {
				continue
			}
			rdTrans := trans.ToResponseData()
			tmpTransactions = append(tmpTransactions, rdTrans)
		}
	}
	return tmpTransactions
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value *big.Int,
	senderPublicKey *ecdsa.PublicKey, s *types.Transaction) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	//if isTransacted {
	//	for _, n := range bc.neighbors {
	//		publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(),
	//			senderPublicKey.Y.Bytes())
	//		signatureStr := s.String()
	//		bt := &Block.TransactionRequest{
	//			&sender, &recipient, &publicKeyStr, &value, &signatureStr}
	//		m, _ := json.Marshal(bt)
	//		buf := bytes.NewBuffer(m)
	//		endpoint := fmt.Sprintf("http://%s/transactions", n)
	//		client := &http.Client{}
	//		req, _ := http.NewRequest("PUT", endpoint, buf)
	//		resp, _ := client.Do(req)
	//		log.Printf("   **  **  **  CreateTransaction : %v", resp)
	//	}
	//}

	return isTransacted
}

func (bc *Blockchain) GetBlockByHash(bHash [32]byte) *Block.Block {
	for _, block := range bc.chain {
		tmpHash := block.Hash()
		if bytes.Equal(tmpHash[:], bHash[:]) {
			return block
		}
	}
	return nil
}

func (bc *Blockchain) GetBlockByNumber(number *big.Int) *Block.Block {
	for _, block := range bc.chain {
		tmpNum := block.Number()
		if tmpNum.Cmp(number) == 0 {
			return block
		}
	}
	return nil
}

func (bc *Blockchain) GetTransactionByHash(tHash [32]byte) *Transaction.Transaction {
	for _, block := range bc.chain {
		for _, trans := range block.Transactions() {
			if trans.From() == bc.blockchainAddress {
				continue
			}
			tmpHash := trans.Hash()
			if bytes.Equal(tmpHash[:], tHash[:]) {
				return trans
			}

		}

	}
	return nil
}
