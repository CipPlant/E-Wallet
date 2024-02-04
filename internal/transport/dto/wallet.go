package dto

import "EWallet/internal/domain/model"

type WalletInput struct {
	WalletID string `json:"walletID"`
	To       string `json:"to"`
	Amount   string `json:"amount"`
}

type WalletOutput struct {
	ID      string `json:"id"`
	Balance string `json:"balance"`
	Time    string `json:"time,omitempty"`
}

func ConvertWalletToDTO(wallet *model.Wallet) WalletOutput {
	return WalletOutput{
		ID:      wallet.WalletID.String(),
		Balance: wallet.BalanceAmount.String(),
	}
}
