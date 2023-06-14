package Transaction

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
)

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		From     string   `json:"from"`
		Nonce    uint64   `json:"nonce"`
		To       string   `json:"to"`
		Value    *big.Int `json:"value"`
		GasLimit uint64   `json:"gas_limit"`
		GasPrice *big.Int `json:"gas_price"`
		ChainID  *big.Int `json:"chain_id"`
	}{
		From:     t.from,
		Nonce:    t.nonce,
		To:       t.to,
		Value:    t.value,
		GasLimit: t.gasLimit,
		GasPrice: t.gasPrice,
		ChainID:  t.chainID,
	})
}

func (t *Transaction) UnmarshalJSON(transactionByte []byte) error {
	type TmpTransaction struct {
		From     string   `json:"from"`
		Nonce    uint64   `json:"nonce"`
		To       string   `json:"to"`
		Value    *big.Int `json:"value"`
		GasLimit uint64   `json:"gas_limit"`
		GasPrice *big.Int `json:"gas_price"`
		ChainID  *big.Int `json:"chain_id"`
	}
	var tmptrans TmpTransaction
	err := json.Unmarshal(transactionByte, &tmptrans)
	if err != nil {
		return err
	}
	t.from = tmptrans.From
	t.nonce = tmptrans.Nonce
	t.to = tmptrans.To
	t.value = tmptrans.Value
	t.gasLimit = tmptrans.GasLimit
	t.gasPrice = tmptrans.GasPrice
	t.chainID = tmptrans.ChainID
	return nil
}

func (t *Transaction) ToResponseData() *TransRecord {
	var tmptrans TransRecord
	tmptrans.From = t.from
	tmptrans.To = t.to
	tmptrans.Value = t.value
	thash := t.Hash()
	tstr := hex.EncodeToString(thash[:])
	tmptrans.Hash = tstr
	return &tmptrans
}
