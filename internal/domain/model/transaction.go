package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type WalletTransaction struct {
	Output []Transaction
	Input  []Transaction
}

type Transaction struct {
	ID              int
	FromWalletID    uuid.UUID
	ToWalletID      uuid.UUID
	Amount          decimal.Decimal
	TransactionTime time.Time
}
