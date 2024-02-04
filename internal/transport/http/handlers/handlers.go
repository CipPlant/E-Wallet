package handlers

import (
	"EWallet/internal/domain/model"
	"EWallet/internal/transport/dto"
)

type DatabaseUseCase interface {
	CreateWallet(input *model.Wallet) (dto.WalletOutput, error)
	SendMoney(input dto.WalletInput) error
	GetWalletHistory(input dto.WalletInput) (dto.TransactionResponse, error)
	GetWalletBalance(wallet dto.WalletInput) (dto.WalletOutput, error)
}
