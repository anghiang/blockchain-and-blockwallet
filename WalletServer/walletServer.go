package main

import (
	"BlockWallet/Transaction"
	utils "BlockWallet/Utils"
	"BlockWallet/Wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

type WalletServer struct {
	port    uint16
	gateWay string
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}
func (ws *WalletServer) Gateway() string {
	return ws.gateWay
}

func NewWalletServer(port uint16, gateWay string) *WalletServer {
	return &WalletServer{port: port, gateWay: gateWay}
}

func (ws *WalletServer) GetPort() uint16 {
	return ws.port
}

func (ws *WalletServer) GetGateWay() string {
	return ws.gateWay
}

func (ws *WalletServer) GetWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//设置允许的方法
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	switch r.Method {
	case http.MethodPost:
		wallet := Wallet.NewWallet()
		walletJson, _ := wallet.MarshalJSON()
		io.WriteString(w, string(walletJson[:]))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Call WalletAmount  METHOD:%s\n", req.Method)
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
		color.Blue("请求查询账户%s的余额", blockchainAddress)

		// 构建请求数据
		requestData := struct {
			BlockchainAddress string `json:"blockchain_address"`
		}{
			BlockchainAddress: blockchainAddress,
		}

		// 将请求数据编码为JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			fmt.Printf("编码JSON时发生错误:%v", err)
			return
		}
		fmt.Println("string(jsonData): ", string(jsonData))
		bcsResp, _ := http.Post(ws.Gateway()+"/amount", "application/json", bytes.NewBuffer(jsonData))

		//返回给客户端
		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar Transaction.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			resp_message := struct {
				Message string  `json:"message"`
				Amount  float64 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			}
			m, _ := json.Marshal(resp_message)
			fmt.Println("string(m)", string(m))
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) LoadWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//设置允许的方法
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	switch r.Method {
	case http.MethodPost:
		privateKey := r.FormValue("privateKey")
		wallet := Wallet.LoadWallet(privateKey)
		walletJson, _ := wallet.MarshalJSON()
		io.WriteString(w, string(walletJson[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}

}

func (ws *WalletServer) CreateTransaction(
	w http.ResponseWriter,
	req *http.Request) {
	defer req.Body.Close()
	switch req.Method {
	case http.MethodPost:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//设置允许的方法
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		var t Transaction.TransactionFrontSign
		log.Println("req.Body==", req.Body)
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&t)
		log.Printf("\n\n\n")
		log.Println("发送人公钥SenderPublicKey ==", *t.SenderPublicKey)
		log.Println("发送人私钥SenderPrivateKey ==", *t.SenderPrivateKey)
		log.Println("发送人地址SenderBlockchainAddress ==", *t.SenderBlockchainAddress)
		log.Println("接收人地址RecipientBlockchainAddress ==", *t.RecipientBlockchainAddress)
		log.Println("金额Value ==", *t.Value)
		value, _ := new(big.Int).SetString(*t.Value, 10)
		tx := Transaction.NewTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value)

		privKey, err := crypto.HexToECDSA(*t.SenderPrivateKey)
		if err != nil {
			fmt.Println("crypto.HexToECDSA Error: ", err)
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("Validate fail")))
			return
		}

		ethTx := tx.ToEthTx()

		signTx, err := Transaction.SignTransaction(ethTx, privKey)
		if err != nil {
			fmt.Println("Transaction.SignTransaction error: ", err)
		}
		w.Header().Add("Content-Type", "application/json")
		bt := &Transaction.TransactionBankSign{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      value,
			SignTx:                     signTx,
		}

		btJson, err := json.Marshal(bt)
		if err != nil {
			fmt.Println("json.Marshal(bt): ", err)
		}
		buf := bytes.NewBuffer(btJson)
		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		fmt.Println(resp.StatusCode)
		if resp.StatusCode == 201 {
			// 201是哪里来的？请参见blockserver  Transactions方法的  w.WriteHeader(http.StatusCreated)语句
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: 非法的HTTP请求方式")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/wallet", ws.GetWallet)
	http.HandleFunc("/loadWallet", ws.LoadWallet)
	http.HandleFunc("/createTransaction", ws.CreateTransaction)
	http.HandleFunc("/amount", ws.WalletAmount)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.GetPort())), nil))
}
