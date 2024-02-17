package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type Storage interface {
	GetAccountByNumber(IBAN string) (*Account, error)
	GetAllAccounts() ([]*Account, error)
	AddAccount(account Account)
	DepositFunds(IBAN string, amount float64) error
	WithdrawFunds(IBAN string, amount float64) error
	TransferFunds(senderIBAN, receiverIBAN string, amount float64) error
	TransferFundsJSON(transferJSON []byte) error
}

type InMemoryStorage struct {
	items map[string]Account
}

type Account struct {
	// Bank account number
	IBAN string `json:"IBAN"`

	// Current account balance
	Balance float64 `json:"balance"`

	// Current account status, can be Active or Blocked
	Status int `json:"status"`
}

type Transfer struct {
	SenderIBAN   string  `json:"sender_IBAN"`
	ReceiverIBAN string  `json:"receiver_IBAN"`
	Amount       float64 `json:"amount"`
}

func (a *Account) String() string {
	return fmt.Sprintf(a.IBAN)
}

// Status
const (
	Active = iota
	Blocked
)

const (
	BelarusCountryCode = "BY"

	StateAccountIBAN       = "BY20OLMP31350000001000000933"
	DestructionAccountIBAN = "BY87MTBK38190000000000353409"
)

// NewStorage creates new storage
func NewStorage() *InMemoryStorage {
	items := make(map[string]Account)

	storage := InMemoryStorage{
		items: items,
	}

	return &storage
}

// InitDefaultAccounts initializes state and destruction accounts
func (i *InMemoryStorage) InitDefaultAccounts() {
	i.items[StateAccountIBAN] = Account{
		IBAN:    StateAccountIBAN,
		Balance: 0.0,
	}

	i.items[DestructionAccountIBAN] = Account{
		IBAN:    DestructionAccountIBAN,
		Balance: 0.0,
	}
}

// GetAccount gets specific account by IBAN
func (i *InMemoryStorage) GetAccount(IBAN string) (*Account, error) {
	account, found := i.items[IBAN]

	if !found {
		return nil, errors.New("аккаунт не найден")
	}

	return &account, nil
}

// GetAllAccounts gets all registered accounts
func (i *InMemoryStorage) GetAllAccounts() ([]*Account, error) {
	accounts := make([]*Account, 0)
	for _, v := range i.items {
		accounts = append(accounts, &v)
	}

	return accounts, nil
}

// AddAccount registers new account
func (i *InMemoryStorage) AddAccount(account Account) {
	i.items[account.IBAN] = account
}

// DepositFunds adds funds into account with specific IBAN
func (i *InMemoryStorage) DepositFunds(IBAN string, amount float64) error {
	account, found := i.items[IBAN]
	if !found {
		return errors.New("аккаунт не найден")
	}

	if account.Status == Blocked {
		return errors.New("аккаунт заблокирован")
	}

	account.Balance += amount

	i.items[IBAN] = account

	return nil
}

// WithdrawFunds removed funds from account with specific IBAN
func (i *InMemoryStorage) WithdrawFunds(IBAN string, amount float64) error {
	account, found := i.items[IBAN]
	if !found {
		return errors.New("аккаунт не найден")
	}

	if account.Status == Blocked {
		return errors.New("аккаунт заблокирован")
	}

	if account.Balance < amount {
		return errors.New("недостаточно средств")
	}

	account.Balance -= amount
	i.items[IBAN] = account

	return nil
}

// TransferFunds transfers funds between two accounts
func (i *InMemoryStorage) TransferFunds(senderIBAN, receiverIBAN string, amount float64) error {
	if err := i.WithdrawFunds(senderIBAN, amount); err != nil {
		return err
	}

	if err := i.DepositFunds(receiverIBAN, amount); err != nil {
		return err
	}

	return nil
}

// TransferFundsJSON transfers funds between two accounts using JSON
func (i *InMemoryStorage) TransferFundsJSON(transferJSON []byte) error {
	var transfer Transfer
	if err := json.Unmarshal(transferJSON, &transfer); err != nil {
		return err
	}

	if err := i.TransferFunds(transfer.SenderIBAN, transfer.ReceiverIBAN, transfer.Amount); err != nil {
		return err
	}

	return nil
}

// GenerateIBAN generates new IBAN with specific country code
func GenerateIBAN(countryCode string) string {
	checkNumber := fmt.Sprintf("%02d", rand.Intn(100))

	bankCode := ""
	for i := 0; i < 4; i++ {
		bankCode += string('A' + rune(rand.Intn(26))) // in ASCII Table: 'A' - 65, 'Z' - 90
	}

	var accountNumber string
	for i := 0; i < 20; i++ {
		accountNumber += fmt.Sprintf("%d", rand.Intn(10))
	}

	return countryCode + checkNumber + bankCode + accountNumber
}

func main() {
	storage := NewStorage()

	storage.InitDefaultAccounts()

	applicationLoop(storage)
}

func applicationLoop(storage *InMemoryStorage) {
	var choice int

	for {
		listMenu()

		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println(err)
			continue
		}

		switch choice {
		case 1:
			account, err := storage.GetAccount(StateAccountIBAN)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(account)
		case 2:
			account, err := storage.GetAccount(DestructionAccountIBAN)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(account)
		case 3:
			fmt.Println("Введите сумму для перевода: ")
			var foundsChoice float64
			if _, err := fmt.Scan(&foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}

			if err := storage.DepositFunds(StateAccountIBAN, foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}
		case 4:
			fmt.Println("Введите сумму для перевода: ")
			var foundsChoice float64
			if _, err := fmt.Scan(&foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Введите IBAN отправителя: ")
			var IBANChoice string
			if _, err := fmt.Scan(&IBANChoice); err != nil {
				fmt.Println(err)
				continue
			}

			if err := storage.TransferFunds(IBANChoice, DestructionAccountIBAN, foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}
		case 5:
			IBAN := GenerateIBAN(BelarusCountryCode)
			newAccount := Account{
				IBAN:    IBAN,
				Balance: 0.0,
				Status:  Active,
			}

			storage.AddAccount(newAccount)

			fmt.Println("Зарегестрирован новый аккаунт c IBAN: ", IBAN)
		case 6:
			listTypeChoices()

			fmt.Println("Выберите тип для перевода: ")
			var typeChoice float64
			if _, err := fmt.Scan(&typeChoice); err != nil {
				fmt.Println(err)
				continue
			}

			switch typeChoice {
			case 1:
				fmt.Println("Введите сумму для перевода: ")
				var foundsChoice float64
				if _, err := fmt.Scan(&foundsChoice); err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println("Введите IBAN отправителя: ")
				var senderIBANChoice string
				if _, err := fmt.Scan(&senderIBANChoice); err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println("Введите IBAN получателя: ")
				var receiverIBANChoice string
				if _, err := fmt.Scan(&receiverIBANChoice); err != nil {
					fmt.Println(err)
					continue
				}

				if err := storage.TransferFunds(senderIBANChoice, receiverIBANChoice, foundsChoice); err != nil {
					fmt.Println(err)
					continue
				}
			case 2:
				fmt.Println("Введите JSON транзакции: ")
				var transferChoice []byte
				if _, err := fmt.Scan(&transferChoice); err != nil {
					fmt.Println(err)
					continue
				}

				if err := storage.TransferFundsJSON(transferChoice); err != nil {
					fmt.Println(err)
					continue
				}
			case 0:
				continue
			}
		case 7:
			accounts, err := storage.GetAllAccounts()
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, acc := range accounts {
				accJSON, err := json.Marshal(acc)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(string(accJSON))
			}
		case 0:
			os.Exit(0)
		default:
			fmt.Println("Повторите свой выбор")
		}
		fmt.Println()
	}
}

func listMenu() {
	fmt.Println("1. Вывести номер специального счета для “эмиссии”")
	fmt.Println("2. Выводить номер специального счета для “уничтожения”")
	fmt.Println("3. Осуществить эмиссию, по добавлению на счет “эмиссии” указанной суммы")
	fmt.Println("4. Осуществить отправку определенной суммы денег c указанного счета на счет “уничтожения”")
	fmt.Println("5. Открыть новый счет")
	fmt.Println("6. Осуществить перевод заданной суммы денег")
	fmt.Println("7. Вывести список всех счетов")
	fmt.Println("0. Выход")
	fmt.Print("Введите ваш выбор: ")
}

func listTypeChoices() {
	fmt.Println("1. C несколькими параметрами")
	fmt.Println("2. B виде JSON")
	fmt.Println("0. Отмена")
}
