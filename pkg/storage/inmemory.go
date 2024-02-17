package storage

import (
	"context"
	"encoding/json"
	"github.com/dima-xd/payment-sys/internal/iban"
	"github.com/dima-xd/payment-sys/internal/transfer"
	"sync"

	"github.com/dima-xd/payment-sys/entity"
)

type InMemoryStorage struct {
	sync.RWMutex
	items map[iban.IBAN]entity.Account
}

// NewStorage creates new storage
func NewStorage() *InMemoryStorage {
	items := make(map[iban.IBAN]entity.Account)

	storage := InMemoryStorage{
		items: items,
	}

	return &storage
}

// GetAccount gets specific account by IBAN
func (i *InMemoryStorage) GetAccount(ctx context.Context, IBAN iban.IBAN) (*entity.Account, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	acc, found := i.items[IBAN]

	if !found {
		return nil, entity.AccountNotFound
	}

	return &acc, nil
}

// GetAllAccounts gets all registered accounts
func (i *InMemoryStorage) GetAllAccounts(ctx context.Context) ([]*entity.Account, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	accounts := make([]*entity.Account, 0)
	for _, v := range i.items {
		accounts = append(accounts, &v)
	}

	return accounts, nil
}

// AddAccount registers new account
func (i *InMemoryStorage) AddAccount(ctx context.Context, account entity.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.Lock()
	defer i.Unlock()

	i.items[account.IBAN] = account
	return nil
}

// DepositFunds adds funds into account with specific IBAN
func (i *InMemoryStorage) DepositFunds(ctx context.Context, IBAN iban.IBAN, amount float64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	acc, found := i.items[IBAN]
	if !found {
		return entity.AccountNotFound
	}

	if acc.Status != entity.Active {
		return entity.AccountBlocked
	}

	i.Lock()
	defer i.Unlock()

	acc.Balance += amount

	i.items[IBAN] = acc

	return nil
}

// WithdrawFunds removed funds from account with specific IBAN
func (i *InMemoryStorage) WithdrawFunds(ctx context.Context, IBAN iban.IBAN, amount float64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	acc, found := i.items[IBAN]
	if !found {
		return entity.AccountNotFound
	}

	if acc.Status != entity.Active {
		return entity.AccountBlocked
	}

	if acc.Balance < amount {
		return entity.InsufficientFunds
	}

	i.Lock()
	defer i.Unlock()

	acc.Balance -= amount
	i.items[IBAN] = acc

	return nil
}

// TransferFunds transfers funds between two accounts
func (i *InMemoryStorage) TransferFunds(ctx context.Context, senderIBAN, receiverIBAN iban.IBAN, amount float64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.BeginTransaction(ctx)
	defer i.Rollback(ctx)

	if err := i.WithdrawFunds(ctx, senderIBAN, amount); err != nil {
		return err
	}

	if err := i.DepositFunds(ctx, receiverIBAN, amount); err != nil {
		return err
	}

	i.Commit(ctx)
	return nil
}

// TransferFundsJSON transfers funds between two accounts using JSON
func (i *InMemoryStorage) TransferFundsJSON(ctx context.Context, transferJSON []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var transf transfer.Transfer
	if err := json.Unmarshal(transferJSON, &transf); err != nil {
		return err
	}

	if err := i.TransferFunds(ctx, transf.SenderIBAN, transf.ReceiverIBAN, transf.Amount); err != nil {
		return err
	}

	return nil
}

func (i *InMemoryStorage) BeginTransaction(ctx context.Context) {
	// TODO: needs BeginTransaction implementation
}

func (i *InMemoryStorage) Rollback(ctx context.Context) {
	// TODO: needs Rollback implementation
}

func (i *InMemoryStorage) Commit(ctx context.Context) {
	// TODO: needs Commit implementation
}
