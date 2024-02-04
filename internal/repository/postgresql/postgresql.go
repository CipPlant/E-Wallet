package postgresql

import (
	"EWallet/internal/domain/model"
	"EWallet/internal/repository/postgresql/outsideErrors"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type PGSQLRepository struct {
	db *pgxpool.Pool
}

func NewPGSQLRepo(name, password, host, port, dbname string, pingTime time.Duration) (*PGSQLRepository, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		name,
		password,
		host,
		port,
		dbname,
	)

	conn, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTime)
	defer cancel()

	err = conn.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("Unable to ping database: %v\n", err)
	}

	return &PGSQLRepository{db: conn}, nil
}

func (p *PGSQLRepository) CreateWallet(id, amount string) (*model.Wallet, error) {

	var resWallet model.Wallet

	tx, err := p.db.Begin(context.Background())

	if err != nil {
		return &model.Wallet{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	err = tx.QueryRow(context.Background(),
		"INSERT INTO Wallet (custom_id, amount) "+
			"VALUES ($1, $2) "+
			"RETURNING custom_id, amount", id, amount).
		Scan(&resWallet.WalletID,
			&resWallet.BalanceAmount)

	if err != nil {
		return &model.Wallet{}, err
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return &model.Wallet{}, err
	}

	return &resWallet, nil
}

func (p *PGSQLRepository) SendMoney(from, to, amount string) error {

	tx, err := p.db.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	var senderBalance decimal.Decimal
	senderQuery := "SELECT amount FROM Wallet WHERE custom_id = $1 FOR UPDATE"
	err = tx.QueryRow(context.Background(), senderQuery, from).Scan(&senderBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return outsideErrors.NoSuchOutgoingWallet
		}
		return err
	}

	if senderBalance.LessThan(decimal.RequireFromString(amount)) {
		return outsideErrors.NotEnoughMoney
	}

	var receiverBalance decimal.Decimal

	receiverQuery := "SELECT amount FROM Wallet WHERE custom_id = $1 FOR UPDATE"

	err = tx.QueryRow(context.Background(), receiverQuery, to).Scan(&receiverBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return outsideErrors.NoSuchDestinationWallet
		}
		return err
	}

	updateSenderQuery := "UPDATE Wallet SET amount = amount - $1 WHERE custom_id = $2"
	_, err = tx.Exec(context.Background(), updateSenderQuery, amount, from)
	if err != nil {
		return err
	}

	updateReceiverQuery := "UPDATE Wallet SET amount = amount + $1 WHERE custom_id = $2"
	_, err = tx.Exec(context.Background(), updateReceiverQuery, amount, to)

	if err != nil {
		return err
	}

	updateTransactionQuery := "INSERT INTO Transactions (from_wallet_id, to_wallet_id, amount, transaction_time) VALUES ($1, $2, $3, $4)"

	_, err = tx.Exec(context.Background(), updateTransactionQuery, from, to, amount, time.Now())

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return err
	}
	return nil
}

func (p *PGSQLRepository) GetWalletHistory(walletID string) (model.WalletTransaction, error) {

	var userCount int

	checkUserQuery := "SELECT id FROM Wallet WHERE custom_id = $1"
	err := p.db.QueryRow(context.Background(), checkUserQuery, walletID).Scan(&userCount)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.WalletTransaction{}, outsideErrors.NoSuchWallet
		}
		return model.WalletTransaction{}, err
	}

	selectQuery := "SELECT id, from_wallet_id, to_wallet_id, amount, transaction_time" +
		" FROM transactions" +
		" WHERE from_wallet_id = $1" +
		" OR to_wallet_id = $1"

	rows, err := p.db.Query(context.Background(), selectQuery, walletID)

	defer rows.Close()

	if err != nil {
		// TODO: what error can be there?
		return model.WalletTransaction{}, err
	}

	response := model.WalletTransaction{
		Input:  make([]model.Transaction, 0),
		Output: make([]model.Transaction, 0),
	}

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.FromWalletID, &transaction.ToWalletID, &transaction.Amount, &transaction.TransactionTime)
		if err != nil {
			// TODO: what error can be there?
			return model.WalletTransaction{}, err
		}
		transactionResponse := model.Transaction{
			TransactionTime: transaction.TransactionTime,
			FromWalletID:    transaction.FromWalletID,
			ToWalletID:      transaction.ToWalletID,
			Amount:          transaction.Amount,
		}

		if transaction.FromWalletID.String() == walletID {
			response.Output = append(response.Output, transactionResponse)
		} else {
			response.Input = append(response.Input, transactionResponse)
		}
	}

	if err := rows.Err(); err != nil {
		// TODO: what error can be there?
		return model.WalletTransaction{}, err
	}

	return response, nil
}

func (p *PGSQLRepository) GetWalletBalance(walletId string) (*model.Wallet, error) {

	resWallet := &model.Wallet{WalletID: uuid.MustParse(walletId)}

	sqlStr := "SELECT amount FROM Wallet WHERE custom_id = $1"

	err := p.db.QueryRow(context.Background(), sqlStr, walletId).Scan(&resWallet.BalanceAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, outsideErrors.NoSuchWallet
		}
	}
	return resWallet, nil
}
