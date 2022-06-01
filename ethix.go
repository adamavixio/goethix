package goethix

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	WeiToEth = big.NewFloat(math.Pow10(18))

	ErrStringToBigFloat = errors.New("invalid conversion from string to big float")
)

type Ethix struct {
	client *ethclient.Client
}

func NewEthix() *Ethix {
	return &Ethix{}
}

//
// Client
//

func (e *Ethix) Client() *ethclient.Client {
	return e.client
}

func (e *Ethix) Dial(url string) error {
	client, err := ethclient.Dial(url)

	if err != nil {
		return err
	}

	e.client = client
	return nil
}

func (e *Ethix) Balance(address string) (string, error) {
	account := common.HexToAddress(address)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, err := e.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return "", err
	}

	float := &big.Float{}
	_, success := float.SetString(balance.String())
	if !success {
		return "", ErrStringToBigFloat
	}

	eth := big.Float{}
	eth.Quo(float, big.NewFloat(math.Pow10(18)))

	return eth.String(), nil
}

func (e *Ethix) Transfer(amount, key, address string) (*types.Transaction, error) {
	toAddress := common.HexToAddress(address)

	edcsaPrivateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	eth, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return nil, ErrStringToBigFloat
	}

	wei := &big.Int{}
	eth.Mul(eth, WeiToEth).Int(wei)

	fromAddress, err := e.address(key)
	if err != nil {
		return nil, err
	}

	nonce, err := e.nonce(fromAddress)
	if err != nil {
		return nil, err
	}

	gas, err := e.gasPrice()
	if err != nil {
		return nil, err
	}

	chainID, err := e.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	limit := uint64(21000)
	tx := types.NewTransaction(nonce, toAddress, wei, limit, gas, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), edcsaPrivateKey)
	if err != nil {
		return nil, err
	}

	err = e.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (e *Ethix) Authorize(privateKey string) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	address, err := e.address(privateKey)
	if err != nil {
		return nil, err
	}

	nonce, err := e.nonce(address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.gasPrice()
	if err != nil {
		return nil, err
	}

	chainID, err := e.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}

func (e *Ethix) Subscribe(contract common.Address) (ethereum.Subscription, chan types.Log, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
	}

	logs := make(chan types.Log)

	sub, err := e.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, nil, err
	}

	return sub, logs, err
}
