package storage

import (
	"context"
	"github.com/dima-xd/payment-sys/entity"
	"github.com/dima-xd/payment-sys/internal/iban"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testAccountActive = entity.Account{
		IBAN:    "BY20OLMP31350000001000000933",
		Balance: 100.0,
		Status:  entity.Active,
	}

	testAccountDisabled = entity.Account{
		IBAN:    "BY87MTBK38190000000000353409",
		Balance: 100.0,
		Status:  entity.Disabled,
	}

	testAccountWithZeroBalance = entity.Account{
		IBAN:    "BY87MTBK98220010020000353411",
		Balance: 0.0,
		Status:  entity.Active,
	}
)

const (
	TestJSON                      = "{\"sender_IBAN\":\"BY20OLMP31350000001000000933\",\"receiver_IBAN\":\"BY87MTBK98220010020000353411\",\"amount\":100}"
	TestJSONWithWrongSenderIBAN   = "{\"sender_IBAN\":\"BY20OLMP31350000001030000933\",\"receiver_IBAN\":\"BY87MTBK98220010020000353411\",\"amount\":100}"
	TestJSONWithWrongReceiverIBAN = "{\"sender_IBAN\":\"BY20OLMP31350000001000000933\",\"receiver_IBAN\":\"BY87MTBK98220010340000353411\",\"amount\":100}"
)

func TestInMemoryStorage_AddAccount(t *testing.T) {
	storage := NewStorage()

	ctx := context.Background()

	err := storage.AddAccount(ctx, testAccountActive)

	assert.NoError(t, err)
}

func TestInMemoryStorage_AddAccount_WithCancel(t *testing.T) {
	storage := NewStorage()

	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	err := storage.AddAccount(ctx, testAccountActive)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestInMemoryStorage_DepositFunds(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.DepositFunds(ctx, testAccountActive.IBAN, 100)
	acc, _ := storage.GetAccount(ctx, testAccountActive.IBAN)

	assert.NoError(t, err)
	assert.Equal(t, testAccountActive.Balance+100, acc.Balance)
}

func TestInMemoryStorage_DepositFunds_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	err := storage.DepositFunds(ctx, testAccountActive.IBAN, 100)

	assert.ErrorIs(t, context.Canceled, err)
}

func TestInMemoryStorage_DepositFunds_OnBlockedAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.DepositFunds(ctx, testAccountDisabled.IBAN, 100)

	assert.ErrorIs(t, err, entity.AccountBlocked)
}

func TestInMemoryStorage_DepositFunds_OnNotFoundAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.DepositFunds(ctx, "testiban", 100)

	assert.ErrorIs(t, err, entity.AccountNotFound)
}

func TestInMemoryStorage_GetAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	acc, err := storage.GetAccount(ctx, testAccountActive.IBAN)

	assert.NoError(t, err)
	assert.Equal(t, &testAccountActive, acc)
}

func TestInMemoryStorage_GetAccount_OnNotFoundAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	acc, err := storage.GetAccount(ctx, "TESTIBAN")

	assert.ErrorIs(t, err, entity.AccountNotFound)
	assert.Nil(t, acc)
}

func TestInMemoryStorage_GetAccount_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	acc, err := storage.GetAccount(ctx, testAccountActive.IBAN)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, acc)
}

func TestInMemoryStorage_GetAllAccounts(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	_, err := storage.GetAllAccounts(ctx)

	assert.NoError(t, err)
}

func TestInMemoryStorage_GetAllAccounts_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	accounts, err := storage.GetAllAccounts(ctx)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, accounts)
}

func TestInMemoryStorage_TransferFunds(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.TransferFunds(ctx, testAccountActive.IBAN, testAccountWithZeroBalance.IBAN, 100)
	senderAcc, _ := storage.GetAccount(ctx, testAccountActive.IBAN)
	receiverAcc, _ := storage.GetAccount(ctx, testAccountWithZeroBalance.IBAN)

	assert.NoError(t, err)
	assert.Equal(t, testAccountActive.Balance-100, senderAcc.Balance)
	assert.Equal(t, testAccountWithZeroBalance.Balance+100, receiverAcc.Balance)
}

func TestInMemoryStorage_TransferFunds_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	err := storage.TransferFunds(ctx, testAccountActive.IBAN, testAccountWithZeroBalance.IBAN, 100)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestInMemoryStorage_TransferFundsJSON(t *testing.T) {

	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.TransferFundsJSON(ctx, []byte(TestJSON))
	senderAcc, _ := storage.GetAccount(ctx, testAccountActive.IBAN)
	receiverAcc, _ := storage.GetAccount(ctx, testAccountWithZeroBalance.IBAN)

	assert.NoError(t, err)
	assert.Equal(t, testAccountActive.Balance-100, senderAcc.Balance)
	assert.Equal(t, testAccountWithZeroBalance.Balance+100, receiverAcc.Balance)
}

func TestInMemoryStorage_TransferFundsJSON_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	err := storage.TransferFundsJSON(ctx, []byte(TestJSON))

	assert.ErrorIs(t, err, context.Canceled)
}

func TestInMemoryStorage_TransferFundsJSON_OnNotFoundSenderAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.TransferFundsJSON(ctx, []byte(TestJSONWithWrongSenderIBAN))

	assert.ErrorIs(t, err, entity.AccountNotFound)
}

func TestInMemoryStorage_TransferFundsJSON_WithWrongJSON(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.TransferFundsJSON(ctx, []byte("testjson"))

	assert.Error(t, err)
}

func TestInMemoryStorage_TransferFundsJSON_OnNotFoundReceiverAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.TransferFundsJSON(ctx, []byte(TestJSONWithWrongReceiverIBAN))

	assert.ErrorIs(t, err, entity.AccountNotFound)
}

func TestInMemoryStorage_WithdrawFunds_OnBlockedAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.WithdrawFunds(ctx, testAccountDisabled.IBAN, 100)

	assert.ErrorIs(t, err, entity.AccountBlocked)
}

func TestInMemoryStorage_WithdrawFunds_OnNotFoundAccount(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.WithdrawFunds(ctx, "TESTIBAN", 100)

	assert.ErrorIs(t, err, entity.AccountNotFound)
}

func TestInMemoryStorage_WithdrawFunds_InsufficientFunds(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx := context.Background()

	err := storage.WithdrawFunds(ctx, testAccountWithZeroBalance.IBAN, 100)

	assert.ErrorIs(t, err, entity.InsufficientFunds)
}

func TestInMemoryStorage_WithdrawFunds_WithCancel(t *testing.T) {
	storage := createStorageWithTestAccounts()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	err := storage.WithdrawFunds(ctx, testAccountWithZeroBalance.IBAN, 100)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestNewStorage(t *testing.T) {
	storage := NewStorage()

	assert.Equal(t, storage.items, make(map[iban.IBAN]entity.Account))
}

func createStorageWithTestAccounts() *InMemoryStorage {
	ctx := context.Background()

	storage := NewStorage()
	storage.AddAccount(ctx, testAccountActive)
	storage.AddAccount(ctx, testAccountDisabled)
	storage.AddAccount(ctx, testAccountWithZeroBalance)

	return storage
}
