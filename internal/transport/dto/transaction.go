package dto

import (
	"EWallet/internal/domain/model"
	"time"
)

type TransactionOutput struct {
	Time   string `json:"time"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type TransactionResponse struct {
	IncomingTransaction []TransactionOutput `json:"incomingTransaction,omitempty"`
	OutgoingTransaction []TransactionOutput `json:"outgoingTransaction,omitempty"`
}

func convertToTransactionOutput(transaction model.Transaction) TransactionOutput {
	return TransactionOutput{
		Time:   transaction.TransactionTime.Format(time.RFC3339),
		From:   transaction.FromWalletID.String(),
		To:     transaction.ToWalletID.String(),
		Amount: transaction.Amount.String(),
	}
}

func СonvertTransactionToDTO(walletTransaction model.WalletTransaction) TransactionResponse {
	var incomingTransactions []TransactionOutput
	var outgoingTransactions []TransactionOutput

	for _, incomingTransaction := range walletTransaction.Input {
		incomingTransactions = append(incomingTransactions, convertToTransactionOutput(incomingTransaction))
	}

	for _, outgoingTransaction := range walletTransaction.Output {
		outgoingTransactions = append(outgoingTransactions, convertToTransactionOutput(outgoingTransaction))
	}

	return TransactionResponse{
		IncomingTransaction: incomingTransactions,
		OutgoingTransaction: outgoingTransactions,
	}
}
