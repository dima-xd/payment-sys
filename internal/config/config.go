package config

import (
	"fmt"
	"os"

	"github.com/dima-xd/payment-sys/internal/iban"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	CountryCode            string    `yaml:"country_code" env-default:"BY"`
	StateAccountIBAN       iban.IBAN `yaml:"state_account_IBAN" env-default:"BY20OLMP31350000001000000933"`
	DestructionAccountIBAN iban.IBAN `yaml:"destruction_account_IBAN" env-default:"BY87MTBK38190000000000353409"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		fmt.Printf("failed to read config: %s\n", err)
	}

	return &cfg
}
