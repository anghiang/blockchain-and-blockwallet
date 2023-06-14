package main

import (
	"BlockWallet/BlockChain"
	"BlockWallet/Transaction"
	utils "BlockWallet/Utils"
	"BlockWallet/Wallet"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

var cache = make(map[string]*BlockChain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *BlockChain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		tmpBlockChain, _ := BlockChain.LoadBlockChainFromLevelDB()
		//fmt.Println("tmpBlockChain:", tmpBlockChain)
		if tmpBlockChain != nil {
			bc, err := BlockChain.LoadBlockFromLevelDB(tmpBlockChain)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			cache["blockchain"] = bc
			return bc
		}

		minersWallet := Wallet.NewWallet()
		blockchainAddress := "0x0B6e9A8CD3901Cf83a1898748D43FF88c06b98ff"
		// NewBlockchain与以前的方法不一样,增加了地址和端口2个参数,是为了区别不同的节点
		bc = BlockChain.NewBlockchain(blockchainAddress, minersWallet.BlockchainAddress(), big.NewInt(2000000000000000000), bcs.Port())
		cache["blockchain"] = bc
		color.Magenta("===矿工帐号信息====\n")
		color.Magenta("矿工private_key\n %v\n", minersWallet.PrivateKeyStr())
		color.Magenta("矿工publick_key\n %v\n", minersWallet.PublicKeyStr())
		color.Magenta("矿工blockchain_address\n %s\n", minersWallet.BlockchainAddress())
		color.Magenta("===============\n")
	}
	return bc
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			// Get:显示交易池的内容，Mine成功后清空交易池
			w.Header().Add("Content-Type", "application/json")
			bc := cache["blockchain"]

			transactions := bc.GetTransactions()
			m, _ := json.Marshal(struct {
				Transactions []*Transaction.Transaction `json:"transactions"`
				Length       int                        `json:"length"`
			}{
				Transactions: transactions,
				Length:       len(transactions),
			})
			io.WriteString(w, string(m[:]))
		}
	case http.MethodPost:
		{
			log.Printf("\n\n\n")
			log.Println("接受到wallet发送的交易")
			decoder := json.NewDecoder(req.Body)
			var t Transaction.TransactionBankSign
			err := decoder.Decode(&t)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("Decode Transaction失败")))
				return
			}

			log.Println("发送人公钥SenderPublicKey:", *t.SenderPublicKey)
			log.Println("发送人私钥SenderPrivateKey:", *t.SenderBlockchainAddress)
			log.Println("接收人地址RecipientBlockchainAddress:", *t.RecipientBlockchainAddress)
			log.Println("金额Value:", *t.Value)
			log.Println("交易Signature:", *t.SignTx)

			if !t.Validate() {
				log.Println("ERROR: missing field(s)")
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			bc := cache["blockchain"]
			senderPrivateKey, err := utils.PublicKeyFromString(*t.SenderPublicKey)
			if err != nil {
				fmt.Println("Failed to decode public key string:", err)
				log.Fatal(err)
			}

			isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress,
				*t.RecipientBlockchainAddress, t.Value, senderPrivateKey, t.SignTx)
			//
			w.Header().Add("Content-Type", "application/json")
			var m []byte
			if !isCreated {
				w.WriteHeader(http.StatusBadRequest)
				m = utils.JsonStatus("fail[from:blockchainServer]")
			} else {
				w.WriteHeader(http.StatusCreated)
				m = utils.JsonStatus("success[from:blockchainServer]")
			}
			io.WriteString(w, string(m))

		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:

		var data map[string]interface{}
		// 解析JSON数据

		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			http.Error(w, "无法解析JSON数据", http.StatusBadRequest)
			return
		}
		// 获取JSON字段的值
		blockchainAddress := data["blockchain_address"].(string)

		color.Green("查询账户: %s 余额请求", blockchainAddress)

		amount := bcs.GetBlockchain().CalculateTotalAmount(blockchainAddress)
		amountFloat := new(big.Float).SetInt(amount)

		// 创建10^18的big.Float类型的除数
		divisorFloat := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

		// 将amountFloat除以divisorFloat
		resultFloat := new(big.Float).Quo(amountFloat, divisorFloat)

		// 将resultFloat转换为float64类型
		tmpAmount, _ := resultFloat.Float64()
		if err != nil {
			fmt.Println("Failed to convert string to float64:", err)
			return
		}
		if err != nil {
			fmt.Println(err)
		}
		ar := &Transaction.AmountResponse{Amount: tmpAmount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) TransactionRecords(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := cache["blockchain"]

		transactions := bc.GetTransactionRecord()

		m, _ := json.Marshal(struct {
			Transactions []*Transaction.TransRecord `json:"transactions"`
			Length       int                        `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) GetBlockByHash(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		bc := cache["blockchain"]
		blockHash := req.FormValue("blockHash")
		// 将十六进制字符串解码为字节数组
		blockHashBytes, err := hex.DecodeString(blockHash)
		if err != nil {
			fmt.Println("无法解码十六进制字符串:", err)
			return
		}
		// 将字节数组转换为[32]byte切片
		var blockHashSlice [32]byte
		copy(blockHashSlice[:], blockHashBytes)
		block := bc.GetBlockByHash(blockHashSlice)
		if block != nil {
			blockBytes, err := block.MarshalJSON()
			if err != nil {
				fmt.Println(err)
			}
			io.WriteString(w, string(blockBytes[:]))
		} else {
			io.WriteString(w, "没有找到该区块")
		}

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) GetBlockByNumber(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		bc := cache["blockchain"]
		blockNumber := req.FormValue("blockNumber")
		// 将十六进制字符串解码为字节数组

		num, err := strconv.Atoi(blockNumber)
		if err != nil {
			fmt.Println(err)
		}

		block := bc.GetBlockByNumber(new(big.Int).SetInt64(int64(num)))
		if block != nil {
			blockBytes, err := block.MarshalJSON()
			if err != nil {
				fmt.Println(err)
			}
			io.WriteString(w, string(blockBytes[:]))
		} else {
			io.WriteString(w, "没有找到该区块")
		}

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) GetTransactionByHash(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		bc := cache["blockchain"]
		transHash := req.FormValue("transHash")
		// 将十六进制字符串解码为字节数组
		blockHashBytes, err := hex.DecodeString(transHash)
		if err != nil {
			fmt.Println("无法解码十六进制字符串:", err)
			return
		}

		// 将字节数组转换为[32]byte切片
		var blockHashSlice [32]byte
		copy(blockHashSlice[:], blockHashBytes)
		trans := bc.GetTransactionByHash(blockHashSlice)
		if trans != nil {
			blockBytes, err := trans.MarshalJSON()
			if err != nil {
				fmt.Println(err)
			}
			io.WriteString(w, string(blockBytes[:]))
		} else {
			io.WriteString(w, "没有找到该交易")
		}

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	bc := bcs.GetBlockchain()
	if bc != nil {
		bc.Run()
	} else {
		fmt.Println("空链")
	}
	http.HandleFunc("/transactions", bcs.Transactions)
	http.HandleFunc("/amount", bcs.Amount)
	http.HandleFunc("/transactionRecords", bcs.TransactionRecords)
	http.HandleFunc("/getBlockByHash", bcs.GetBlockByHash)
	http.HandleFunc("/getBlockByNumber", bcs.GetBlockByNumber)
	http.HandleFunc("/getTransactionByHash", bcs.GetTransactionByHash)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(bcs.Port())), nil))
}
