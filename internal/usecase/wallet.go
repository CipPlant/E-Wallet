package usecase

import (
	"EWallet/internal/domain/model"
	"EWallet/internal/repository/postgresql/outsideErrors"
	"EWallet/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Database interface {
	CreateWallet(id, amount string) (*model.Wallet, error)
	SendMoney(from, to, amount string) error
	GetWalletHistory(walletId string) (model.WalletTransaction, error)
	GetWalletBalance(walletId string) (*model.Wallet, error)
}

type DatabaseUseCase struct {
	repo Database
}

func NewRepository(sqlRepository Database) *DatabaseUseCase {
	return &DatabaseUseCase{repo: sqlRepository}
}

func (du *DatabaseUseCase) CreateWallet(input *model.Wallet) (dto.WalletOutput, error) {

	wallet, err := du.repo.CreateWallet(input.WalletID.String(), input.BalanceAmount.String())

	if err != nil {
		return dto.WalletOutput{}, err
	}

	var result dto.WalletOutput

	result.ID = wallet.WalletID.String()
	result.Balance = wallet.BalanceAmount.String()

	return result, nil
}

func (du *DatabaseUseCase) SendMoney(input dto.WalletInput) error {
	var receiverWallet model.Wallet
	var err error

	receiverWallet.BalanceAmount, err = decimal.NewFromString(input.Amount)
	if err != nil {
		return outsideErrors.InvalidMoneyFormat
	}

	receiverWallet.WalletID, err = uuid.Parse(input.To)
	if err != nil {
		return outsideErrors.InvalidWalletFormat
	}

	var SendWallet model.Wallet

	SendWallet.BalanceAmount, err = decimal.NewFromString(input.Amount)
	if err != nil {
		return outsideErrors.InvalidMoneyFormat
	}

	SendWallet.WalletID, err = uuid.Parse(input.WalletID)
	if err != nil {
		return outsideErrors.InvalidWalletFormat
	}

	err = du.repo.SendMoney(SendWallet.WalletID.String(), receiverWallet.WalletID.String(), receiverWallet.BalanceAmount.String())
	if err != nil {
		return err
	}

	return nil
}

func (du *DatabaseUseCase) GetWalletHistory(wallet dto.WalletInput) (dto.TransactionResponse, error) {
	walletID, err := uuid.Parse(wallet.WalletID)
	if err != nil {
		return dto.TransactionResponse{}, outsideErrors.InvalidWalletFormat
	}

	modelWallet := model.Wallet{
		WalletID: walletID,
	}

	walletHistory, err := du.repo.GetWalletHistory(modelWallet.WalletID.String())
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	toDTO := dto.СonvertTransactionToDTO(walletHistory)

	return toDTO, nil
}

func (du *DatabaseUseCase) GetWalletBalance(wallet dto.WalletInput) (dto.WalletOutput, error) {
	walletID, err := uuid.Parse(wallet.WalletID)
	if err != nil {
		return dto.WalletOutput{}, outsideErrors.InvalidWalletFormat
	}

	modelWallet := model.Wallet{
		WalletID: walletID,
	}

	dtoOutput, err := du.repo.GetWalletBalance(modelWallet.WalletID.String())
	if err != nil {
		return dto.WalletOutput{}, err
	}

	toDTO := dto.ConvertWalletToDTO(dtoOutput)

	return toDTO, nil
}
