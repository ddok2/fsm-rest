package config

import (
	"path/filepath"

	"github.com/olebedev/config"
)

type Config struct {
	ServerPort           string
	OperationMode        string
	TxInterval           string
	RemittanceAmount     string
	RemittanceFee        string
	ChargeAmount         string
	AdminChargeAmount    string
	TransactionSendCount string
	BoosterAddr          string
	BoosterPort          string
	ExchangeAddr         string
	ExchangePort         string
	MemberCount          int
	MemberPrefix         string
	AdminId              string
	TrashId              string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(cfgAbsPath string) error {
	filename, err := filepath.Abs(cfgAbsPath)
	if err != nil {
		return err
	}

	cfg, err := config.ParseYamlFile(filename)
	if err != nil {
		return err
	}

	if err = parseAutomationConfig(c, cfg); err != nil {
		return err
	}

	if err = parseBoosterConfig(c, cfg); err != nil {
		return err
	}

	if err = parseDexConfig(c, cfg); err != nil {
		return err
	}

	if err = parseBotConfig(c, cfg); err != nil {
		return err
	}

	return nil
}

func parseAutomationConfig(c *Config, cfg *config.Config) error {
	var err error

	if c.ServerPort, err = cfg.String("tester.port"); err != nil {
		return err
	}
	if c.OperationMode, err = cfg.String("tester.operationMode"); err != nil {
		return err
	}
	if c.TxInterval, err = cfg.String("tester.txInterval"); err != nil {
		return err
	}
	if c.RemittanceAmount, err = cfg.String("tester.remittanceAmount"); err != nil {
		return err
	}
	if c.RemittanceFee, err = cfg.String("tester.remittanceFee"); err != nil {
		return err
	}
	if c.ChargeAmount, err = cfg.String("tester.chargeAmount"); err != nil {
		return err
	}
	if c.AdminChargeAmount, err = cfg.String("tester.adminChargeAmount"); err != nil {
		return err
	}
	if c.TransactionSendCount, err = cfg.String("tester.transactionSendCount"); err != nil {
		return err
	}

	return err
}

func parseBoosterConfig(c *Config, cfg *config.Config) error {
	var err error

	if c.BoosterAddr, err = cfg.String("booster.addr"); err != nil {
		return err
	}
	if c.BoosterPort, err = cfg.String("booster.port"); err != nil {
		return err
	}

	return err
}

func parseDexConfig(c *Config, cfg *config.Config) error {
	var err error

	if c.ExchangeAddr, err = cfg.String("exchange.addr"); err != nil {
		return err
	}
	if c.ExchangePort, err = cfg.String("exchange.port"); err != nil {
		return err
	}

	return err
}

func parseBotConfig(c *Config, cfg *config.Config) error {
	var err error

	if c.MemberCount, err = cfg.Int("members.count"); err != nil {
		return err
	}
	if c.MemberPrefix, err = cfg.String("members.prefix"); err != nil {
		return err
	}
	if c.AdminId, err = cfg.String("members.adminId"); err != nil {
		return err
	}
	if c.TrashId, err = cfg.String("members.trashId"); err != nil {
		return err
	}
	return err
}
