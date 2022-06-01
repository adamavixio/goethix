package goethix

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (e *Ethix) publicKey(privateKey string) (*ecdsa.PublicKey, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &key.PublicKey, nil
}

func (e *Ethix) address(privateKey string) (string, error) {
	publicKey, err := e.publicKey(privateKey)
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	return fromAddress.Hex(), nil
}

func (e *Ethix) gasPrice() (*big.Int, error) {
	return e.client.SuggestGasPrice(context.Background())
}

func (e *Ethix) nonce(address string) (uint64, error) {
	hex := common.HexToAddress(address)

	nonce, err := e.client.PendingNonceAt(context.Background(), hex)
	if err != nil {
		return math.MaxUint64, err
	}

	return nonce, nil
}
