package BlockChain

import (
	"BlockWallet/Block"
	"BlockWallet/Transaction"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"
)

type Blockchain struct {
	transactionPool   []*Transaction.Transaction `json:"transaction_pool"`
	chain             []*Block.Block             `json:"chain"`
	blockchainAddress string                     `json:"blockchain_address"`
	coinbase          string                     `json:"coinbase"`
	port              uint16                     `json:"port"`
	mux               sync.Mutex                 `json:"mux"`
	neighbors         []string                   `json:"neighbors"`
	muxNeighbors      sync.Mutex                 `json:"mux_neighbors"`
	MiningReward      *big.Int                   `json:"mining_reward"`
	difficult         int                        `json:"difficult"`
}

// 新建一条链的第一个区块
// NewBlockchain(blockchainAddress string) *Blockchain
// 函数定义了一个创建区块链的方法，它接收一个字符串类型的参数 blockchainAddress，
// 它返回一个区块链类型的指针。在函数内部，它创建一个区块链对象并为其设置地址，
// 然后创建一个创世块并将其添加到区块链中，最后返回区块链对象。
func NewBlockchain(blockchainAddress string, coinbase string, miningReward *big.Int, port uint16) *Blockchain {
	bc := new(Blockchain)
	b := &Block.Block{}
	bc.CreateBlock(big.NewInt(0), b.CalHash()) //创世纪块
	bc.blockchainAddress = blockchainAddress
	bc.coinbase = coinbase
	bc.MiningReward = miningReward
	bc.port = port
	bc.difficult = 0x80000
	err := SaveBlockChainToLevelDB(bc)
	if err != nil {
		log.Println("Failed to save blockchain to LevelDB:", err)
		// 处理错误
	}
	return bc
}

func (bc *Blockchain) Chain() []*Block.Block {
	return bc.chain
}

func (bc *Blockchain) SetChain(chain []*Block.Block) {
	bc.chain = chain
}
func (bc *Blockchain) TransactionPool() []*Transaction.Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) CreateBlock(nonce *big.Int, previousHash [32]byte) *Block.Block {
	var blockNumber *big.Int
	if bc.chain == nil {
		blockNumber = big.NewInt(0)
	} else {
		previousBlock := bc.LastBlock()
		blockNumber = new(big.Int).Add(previousBlock.Number(), big.NewInt(1))
	}
	b := Block.NewBlock(nonce, previousHash, bc.transactionPool)
	b.SetNumber(blockNumber)
	b.SetDifficulty(bc.difficult)
	b.SetHash(b.CalHash())
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction.Transaction{}

	err := SaveBlockToLevelDB(b)
	if err != nil {
		log.Println("Failed to save block to LevelDB:", err)
	}
	color.Green("Successful commit new Block")
	log.Printf("Commit new Block,number=%d  hash=%30x", b.Number(), b.Hash())
	return b
}

func (bc *Blockchain) AddTransaction(
	sender string,
	recipient string,
	value *big.Int,
	senderPublicKey *ecdsa.PublicKey,
	s *types.Transaction) bool {
	t := Transaction.NewTransaction(sender, recipient, value)

	if sender == bc.blockchainAddress {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.CalculateTotalAmount(sender).Cmp(value) == -1 {
		log.Printf("ERROR : %s 账户中没有足够的钱", sender)
		return false
	}
	if Transaction.VerifyTransactionSignature(s, senderPublicKey) {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: 验证交易")
	}
	return false
}

func (bc *Blockchain) CalculateTotalAmount(accountAddress string) *big.Int {
	var totalAmount = big.NewInt(0)
	for _, _chain := range bc.chain {
		for _, _tx := range _chain.Transactions() {
			if accountAddress == _tx.To() {
				totalAmount = totalAmount.Add(totalAmount, _tx.Value())
			}
			if accountAddress == _tx.From() {
				totalAmount = totalAmount.Sub(totalAmount, _tx.Value())
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction.Transaction {
	transactions := make([]*Transaction.Transaction, 0)
	for _, v := range bc.transactionPool {
		transactions = append(transactions,
			Transaction.NewTransaction(
				v.From(),
				v.To(),
				v.Value()))
	}
	return transactions
}

func (bc *Blockchain) LastBlock() *Block.Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) ProofOfWork() *big.Int {
	transactions := bc.CopyTransactionPool() //选择交易？控制交易数量？
	previousHash := bc.LastBlock().Hash()
	nonce := big.NewInt(0)
	begin := time.Now()
	if bc.getTime(len(bc.chain)-1) < 3e+9 {
		bc.difficult += 32
	} else if bc.getTime(len(bc.chain)-1) > 2e+9 {
		bc.difficult -= 32
	} else {
		bc.difficult = 2
	}
	for !bc.ValidProof(nonce, previousHash, transactions, bc.difficult) {
		nonce = nonce.Add(nonce, big.NewInt(1))
	}
	end := time.Now()
	//log.Printf("POW spend Time:%f Second", end.Sub(begin).Seconds())
	log.Printf("POW spend Time:%s Diff: %d", end.Sub(begin), bc.difficult)
	return nonce
}

func (bc *Blockchain) getTime(i int) uint64 {
	if i == 0 {
		return 0
	}
	return bc.chain[i].Timestamp() - bc.chain[i-1].Timestamp()
}

func (bc *Blockchain) ValidProof(nonce *big.Int,
	previousHash [32]byte,
	transactions []*Transaction.Transaction,
	difficulty int,
) bool {
	big_2 := big.NewInt(2)
	big_256 := big.NewInt(256)
	big_diff := big.NewInt(int64(difficulty))
	target := new(big.Int).Exp(big_2, big_256, nil)
	target = new(big.Int).Div(target, big_diff)
	tmpBlock := Block.Block{}
	tmpBlock.SetNonce(nonce)
	tmpBlock.SetPreviousHash(previousHash)
	tmpBlock.SetTransactions(transactions)
	tmpBlock.SetTimestamp(0)
	//log.Println("guessHashStr", tmpHashStr)
	result := bytesToBigInt(tmpBlock.CalHash())
	return target.Cmp(result) > 0
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		color.Green("%s BLOCK %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	color.Yellow("%s\n\n\n", strings.Repeat("*", 50))
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.AddTransaction(bc.blockchainAddress, bc.coinbase, bc.MiningReward, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()

	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	//for _, n := range bc.neighbors {
	//	endpoint := fmt.Sprintf("http://%s/consensus", n)
	//	client := &http.Client{}
	//	req, _ := http.NewRequest("PUT", endpoint, nil)
	//	resp, _ := client.Do(req)
	//	log.Printf("%v", resp)
	//}
	return true
}
