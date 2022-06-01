package goethix

import (
	"github.com/adamavixio/logger"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type AccountHandler struct {
	ks *keystore.KeyStore
	am *accounts.Manager
	sn *accounts.Account
}

func NewAccountHandler(keystorePath string) *AccountHandler {
	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)

	return &AccountHandler{
		ks: ks,
		am: am,
	}
}

func (h *AccountHandler) GetAccount(address string) *accounts.Account {
	for _, acc := range h.ks.Accounts() {
		if acc.Address.Hex() == address {
			return &acc
		}
	}

	return nil
}

func (h *AccountHandler) CreateAccount(passphrase string) *accounts.Account {
	account, err := h.ks.NewAccount(passphrase)
	logger.Error(err, "unable to create wallet")
	return &account
}

func (h *AccountHandler) UpdateAccount(account accounts.Account, old, new string) {
	err := h.ks.Update(account, old, new)
	logger.Error(err, "unable to update wallet")
}

func (h *AccountHandler) DeleteAccount(account accounts.Account, passphrase string) {
	err := h.ks.Delete(account, passphrase)
	logger.Error(err, "unable to delete wallet")
}

func (h *AccountHandler) ImportAccount(key []byte, old, new string) *accounts.Account {
	account, err := h.ks.Import(key, old, new)
	logger.Error(err, "unable to import wallet")
	return &account
}

func (h *AccountHandler) ExportAccount(account accounts.Account, old, new string) []byte {
	key, err := h.ks.Export(account, old, new)
	logger.Error(err, "unable to export wallet")
	return key
}

func (h *AccountHandler) Signer(passphrase string, hash string) error {
	signer, err := h.ks.NewAccount(passphrase)
	if err != nil {
		return err
	}

	h.sn = &signer
	return nil
}

func (h *AccountHandler) Sign(passphrase string) []byte {
	hash := common.HexToHash("").Bytes()
	sig, err := h.ks.SignHashWithPassphrase(*h.sn, passphrase, hash)
	logger.Error(err, "unable to sign")
	return sig
}
