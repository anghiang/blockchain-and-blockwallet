package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}

func PublicKeyFromString(publicKeyStr string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}

	curve := elliptic.P256() // 或者使用 elliptic.P256k1()，具体取决于您的需求

	x := new(big.Int).SetBytes(publicKeyBytes[:32])
	y := new(big.Int).SetBytes(publicKeyBytes[32:])
	publicKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return &publicKey, nil
}
