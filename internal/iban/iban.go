package iban

import (
	"fmt"
	"math/rand"
)

type IBAN string

// GenerateIBAN generates new IBAN with specific country code
func GenerateIBAN(countryCode string) IBAN {
	checkNumber := fmt.Sprintf("%02d", rand.Intn(100))

	bankCode := ""
	for i := 0; i < 4; i++ {
		bankCode += string('A' + rune(rand.Intn(26))) // in ASCII Table: 'A' - 65, 'Z' - 90
	}

	var accountNumber string
	for i := 0; i < 20; i++ {
		accountNumber += fmt.Sprintf("%d", rand.Intn(10))
	}

	return IBAN(countryCode + checkNumber + bankCode + accountNumber)
}

func (i IBAN) String() string {
	return string(i)
}
