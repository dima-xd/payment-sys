package entity

import (
	"errors"
	"github.com/dima-xd/payment-sys/internal/iban"
)

type Account struct {
	// Bank account number
	IBAN iban.IBAN `json:"IBAN"`

	// Current account balance
	Balance float64 `json:"balance"`

	// Current account status, can be Active or Blocked
	Status `json:"status"`
}

type Status string

var (
	AccountNotFound   = errors.New("аккаунт не найден")
	AccountBlocked    = errors.New("аккаунт заблокирован")
	InsufficientFunds = errors.New("недостаточно средств")
)

// Status
const (
	Active   Status = "ACTIVE"
	Disabled Status = "DISABLED"
)

func (a *Account) String() string {
	return a.IBAN.String()
}
