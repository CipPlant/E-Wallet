package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sync"
)

type Wallet struct {
	WalletID      uuid.UUID
	m             sync.RWMutex // for future balanceAmount surrender
	BalanceAmount decimal.Decimal
}

func NewWallet() *Wallet {
	var wallet Wallet

	wallet.m.Lock()
	defer wallet.m.Unlock()

	wallet.WalletID = uuid.New()
	wallet.BalanceAmount = decimal.NewFromFloat(100.00)

	return &wallet
}

// for future not only CRUD requests

/*
func (w *Wallet) Deposit(amount decimal.Decimal) {
	w.m.Lock()
	defer w.m.Unlock()
	w.BalanceAmount.Add(amount)
}

func (w *Wallet) WithDraw(amount decimal.Decimal) error {
	w.m.Lock()
	defer w.m.Unlock()

	if w.BalanceAmount.Sub(amount).IsNegative() {
		return outsideErrors.NotEnoughMoney
	} else {
		w.BalanceAmount.Sub(amount)
		return nil
	}
}
*/
