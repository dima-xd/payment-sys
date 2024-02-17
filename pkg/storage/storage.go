package storage

import (
	"context"
	"github.com/dima-xd/payment-sys/entity"
	"github.com/dima-xd/payment-sys/internal/iban"
)

type Storage interface {
	GetAccount(ctx context.Context, IBAN iban.IBAN) (*entity.Account, error)
	GetAllAccounts(ctx context.Context) ([]*entity.Account, error)
	AddAccount(ctx context.Context, account entity.Account) error
	DepositFunds(ctx context.Context, IBAN iban.IBAN, amount float64) error
	WithdrawFunds(ctx context.Context, IBAN iban.IBAN, amount float64) error
	TransferFunds(ctx context.Context, senderIBAN, receiverIBAN iban.IBAN, amount float64) error
	TransferFundsJSON(ctx context.Context, transferJSON []byte) error
}

type TransactionalStorage interface {
	BeginTransaction(ctx context.Context)
	Rollback(ctx context.Context)
	Commit(ctx context.Context)
}
