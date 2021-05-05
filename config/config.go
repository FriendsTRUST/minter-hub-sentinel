package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	Telegram  Telegram  `yaml:"telegram"`
	MinterHub MinterHub `yaml:"minter_hub"`
}

type Telegram struct {
	Token  string `yaml:"token"`
	Admins []int  `yaml:"admins"`
}

type MinterHub struct {
	Api              []string `yaml:"api"`
	ValidatorAddress string   `yaml:"validator_address"`
	Sleep            int      `yaml:"sleep"`
}

func New(path string) (*Config, error) {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipDefaults: true,
		SkipFiles:    false,
		SkipEnv:      true,
		SkipFlags:    true,
		Files:        []string{path},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
