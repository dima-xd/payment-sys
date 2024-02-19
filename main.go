package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dima-xd/payment-sys/internal/config"
	"github.com/dima-xd/payment-sys/internal/iban"

	"github.com/dima-xd/payment-sys/entity"
	"github.com/dima-xd/payment-sys/pkg/storage"
)

func main() {
	inMemoryStorage := storage.NewStorage()

	cfg := config.LoadConfig()

	applicationLoop(cfg, inMemoryStorage)
}

func applicationLoop(cfg *config.Config, storage storage.Storage) {
	var choice int
	ctx, cancel := context.WithCancel(context.Background())

	stateAccount := entity.Account{
		IBAN:    cfg.StateAccountIBAN,
		Balance: 0.0,
		Status:  entity.Active,
	}

	destructionAccount := entity.Account{
		IBAN:    cfg.DestructionAccountIBAN,
		Balance: 0.0,
		Status:  entity.Active,
	}

	if err := storage.AddAccount(ctx, stateAccount); err != nil {
		fmt.Println(err)
	}

	if err := storage.AddAccount(ctx, destructionAccount); err != nil {
		fmt.Println(err)
	}

	for {
		listMenu()

		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println(err)
			continue
		}

		switch choice {
		case 1:
			acc, err := storage.GetAccount(ctx, cfg.StateAccountIBAN)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(acc)
		case 2:
			acc, err := storage.GetAccount(ctx, cfg.DestructionAccountIBAN)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(acc)
		case 3:
			fmt.Println("Введите сумму для перевода: ")
			var foundsChoice float64
			if _, err := fmt.Scan(&foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}

			if err := storage.DepositFunds(ctx, cfg.StateAccountIBAN, foundsChoice); err != nil {
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
			var IBANChoice iban.IBAN
			if _, err := fmt.Scan(&IBANChoice); err != nil {
				fmt.Println(err)
				continue
			}

			if err := storage.TransferFunds(ctx, IBANChoice, cfg.DestructionAccountIBAN, foundsChoice); err != nil {
				fmt.Println(err)
				continue
			}
		case 5:
			IBAN := iban.GenerateIBAN(cfg.CountryCode)
			newAccount := entity.Account{
				IBAN:    IBAN,
				Balance: 0.0,
				Status:  entity.Active,
			}

			if err := storage.AddAccount(ctx, newAccount); err != nil {
				fmt.Println(err)
				continue
			}

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
				var senderIBANChoice iban.IBAN
				if _, err := fmt.Scan(&senderIBANChoice); err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println("Введите IBAN получателя: ")
				var receiverIBANChoice iban.IBAN
				if _, err := fmt.Scan(&receiverIBANChoice); err != nil {
					fmt.Println(err)
					continue
				}

				if err := storage.TransferFunds(ctx, senderIBANChoice, receiverIBANChoice, foundsChoice); err != nil {
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

				if err := storage.TransferFundsJSON(ctx, transferChoice); err != nil {
					fmt.Println(err)
					continue
				}
			case 0:
				continue
			}
		case 7:
			accounts, err := storage.GetAllAccounts(ctx)
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
			cancel()
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
