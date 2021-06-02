package config

import (
	"path/filepath"
	"testing"

	"github.com/olebedev/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {

	t.Run("get Filepath", func(t *testing.T) {
		filename, err := filepath.Abs("../../common/config/config.yml")
		assert.Nil(t, err)
		assert.Equal(
			t,
			filename,
			"/Users/sung/Development/10.nuritelecom/03.go/blockchain.automation/common/config/config.yml",
		)
	})

	t.Run("parse yml file", func(t *testing.T) {
		filename, err := filepath.Abs("./config.yml")
		assert.Nil(t, err)

		cfg, err := config.ParseYamlFile(filename)
		assert.Nil(t, err)
		assert.NotNil(t, cfg)

		c := NewConfig()

		err = parseAutomationConfig(c, cfg)
		assert.Nil(t, err)
		assert.Equal(t, c.ServerPort, "8089")
		assert.Equal(t, c.TxInterval, "100")
		assert.Equal(t, c.RemittanceFee, "0")

		err = parseBoosterConfig(c, cfg)
		assert.Nil(t, err)
		assert.Equal(t, c.BoosterAddr, "txbooster.nuriflex.com")
		assert.Equal(t, c.BoosterPort, "8080")

		err = parseDexConfig(c, cfg)
		assert.Nil(t, err)
		assert.Equal(t, c.ExchangeAddr, "dex")
	})

	t.Run("load", func(t *testing.T) {
		c := NewConfig()
		err := c.Load("../../common/config/config.yml")
		assert.Nil(t, err)

		assert.Equal(t, c.ServerPort, "8089")
		assert.Equal(t, c.TxInterval, "100")
		assert.Equal(t, c.RemittanceFee, "0")
		assert.Equal(t, c.BoosterAddr, "txbooster.nuriflex.com")
		assert.Equal(t, c.BoosterPort, "8080")
		assert.Equal(t, c.ExchangeAddr, "dex")
	})
}
