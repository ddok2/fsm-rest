package common

import (
	"io/ioutil"
	"path/filepath"

	"blockchain.automation/configloader"
	"gopkg.in/yaml.v2"
)

var config *Config

type ConfigYaml struct {
	Tester struct {
		Port          string `yaml:"port"`
		OperationMode string `yaml:"operationMode"`
	}
	Booster struct {
		Addr string `yaml:"addr"`
	}
	Exchange struct {
		Addr string `yaml:"addr"`
	}
	Members struct {
		Count  int    `yaml:"count"`
		Prefix string `yaml:"prefix"`
	}
}

type Config struct {
	conf *configloader.Config

	ServerPort    string
	OperationMode string
	BoosterAddr   string
	ExchangeAddr  string
	MemberCount   int
	MemberPrefix  string
}

func NewConfig() *Config {
	if config == nil {
		config = new(Config)
	}

	return config
}

func (c *Config) LoadConfigYaml() error {

	filename, _ := filepath.Abs("./configs/automation_test.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var conf ConfigYaml

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return err
	}

	c.ServerPort = conf.Tester.Port
	c.OperationMode = conf.Tester.OperationMode
	c.BoosterAddr = conf.Booster.Addr
	c.ExchangeAddr = conf.Exchange.Addr
	c.MemberCount = conf.Members.Count
	c.MemberPrefix = conf.Members.Prefix

	return nil
}

func (c *Config) LoadConfigJson() error {

	confType := "json"
	c.conf = configloader.NewConfig()

	if err := c.conf.LoadConfigs("./configs", confType, "AutomationTester.json"); err != nil {
		return err
	}

	c.ServerPort = c.conf.ValueString("tester.port")
	c.OperationMode = c.conf.ValueString("tester.operationMode")
	c.BoosterAddr = c.conf.ValueString("booster.Addr")
	c.MemberCount = c.conf.ValueInt("members.Count")

	return nil
}
