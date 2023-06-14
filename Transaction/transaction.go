package Transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Transaction struct {
	from     string   `json:"from"`
	nonce    uint64   `json:"nonce"`
	to       string   `json:"to"`
	value    *big.Int `json:"value"`
	gasLimit uint64   `json:"gas_limit"`
	gasPrice *big.Int `json:"gas_price"`
	chainID  *big.Int `json:"chain_id"`
}

func NewTransaction(sender string, receiver string, value *big.Int) *Transaction {
	t := Transaction{}
	t.from = sender
	t.to = receiver
	t.value = value
	t.chainID = big.NewInt(1)
	t.gasLimit = uint64(21000)
	t.gasPrice = big.NewInt(20000000000)
	return &t
}

func (t *Transaction) ToEthTx() *types.Transaction {
	recipientAddress := common.HexToAddress(t.to)
	return types.NewTransaction(t.nonce, recipientAddress, t.value, t.gasLimit, t.gasPrice, nil)
}

func (t *Transaction) Value() *big.Int {
	return t.value
}

func (t *Transaction) From() string {
	return t.from
}

func (t *Transaction) To() string {
	return t.to
}

func (t *Transaction) Hash() [32]byte {
	m, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(m)
}

func SignTransaction(tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	//chainID := new(big.Int)
	signer := types.NewEIP155Signer(tx.ChainId())
	// 获取签名哈希
	hash := signer.Hash(tx).Bytes()

	// 对哈希进行签名
	signature, err := crypto.Sign(hash, privateKey)
	//fmt.Println("signature: ", signature)
	if err != nil {
		return nil, err
	}

	// 设置交易的签名
	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		fmt.Println(err)
	}

	return signedTx, nil
}

// 验证交易签名
func VerifyTransactionSignature(tx *types.Transaction, fromPubkey *ecdsa.PublicKey) bool {
	// 获取交易哈希值
	signer := types.NewEIP155Signer(tx.ChainId())
	hash := signer.Hash(tx).Bytes()

	// 获取签名信息
	v, r, s := tx.RawSignatureValues()

	var adjustedV byte
	chainIdInt64 := tx.ChainId().Int64()
	vInt64 := v.Int64()
	if vInt64 == chainIdInt64 || vInt64 == chainIdInt64+1 {
		adjustedV = byte(vInt64 - 2*chainIdInt64 - 36)
	} else {
		adjustedV = byte(vInt64 - 2*chainIdInt64 - 35)
	}
	//// 构建签名
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	vBytes := []byte{adjustedV}

	// 创建一个具有合适大小的字节数组来存储 r 和 s
	sig := make([]byte, 64)

	// 将 r 的字节复制到 sig
	copy(sig[:32], rBytes)

	// 将 s 的字节复制到 sig
	copy(sig[32:], sBytes)
	sig = append(sig, vBytes...)
	//fmt.Println("signature2: ", sig)
	// 使用 Ecrecover 进行签名验证
	recoveredPublicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		fmt.Println(err)
	}
	// 将公钥结构转换为以太坊地址
	signerAddress := crypto.PubkeyToAddress(*recoveredPublicKey)

	fromAddress := crypto.PubkeyToAddress(*fromPubkey)

	return bytes.Equal(signerAddress.Bytes(), fromAddress.Bytes())
}
