package action

import (
	"strconv"
	"testing"

	"blockchain.automation/common/config"
	"blockchain.automation/core/bot"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBcBotAction(t *testing.T) {

	t.Run("NewWallet", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.Load("../../../common/config/config.yml")
		assert.Nil(t, err)

		// action := NewBcBotAction(bot.NewJobs(cfg))
		action := NewAction(cfg)
		isDone, err := action.NewWallet(
			bot.NewBot(
				uuid.Must(uuid.NewUUID()).String(),
				0,
				"test_token",
				func(bot *bot.Bot) {}))
		assert.Nil(t, err)
		assert.True(t, isDone)
	})

	t.Run("Issue", func(t *testing.T) {
		t.Skip("skipping this.")
		cfg := config.NewConfig()
		err := cfg.Load("../../../common/config/config.yml")
		assert.Nil(t, err)

		action := NewAction(cfg)
		adminBot := action.GetAdminBot()
		assert.NotNil(t, adminBot)

		adminBot.Balance = 10000000000

		isDone, err := action.Issue(nil)
		assert.Nil(t, err)
		assert.True(t, isDone)
	})

	t.Run("Trade-Issue to wallet", func(t *testing.T) {
		t.Skip("skipping this. pass.")
		cfg := config.NewConfig()
		err := cfg.Load("../../../common/config/config.yml")
		assert.Nil(t, err)

		// update admin current balance
		action := NewAction(cfg)
		adminBot := action.GetAdminBot()
		assert.NotNil(t, adminBot)
		adminBot.Balance = 20000000000

		// new wallet
		newId, err := uuid.NewUUID()
		assert.Nil(t, err)
		isDone, err := action.NewWallet(
			bot.NewBot(
				newId.String(),
				0,
				"test_token",
				func(bot *bot.Bot) {}))
		assert.Nil(t, err)
		assert.True(t, isDone)

		// issue to new wallet
		isDone, err = action.Trade(
			adminBot,
			bot.NewBot(
				newId.String(),
				0,
				"test_token",
				func(bot *bot.Bot) {}))
		assert.Nil(t, err)
		assert.True(t, isDone)

		// verify balance
		chargeAmount, err := strconv.ParseFloat(cfg.ChargeAmount, 64)
		assert.Nil(t, err)
		assert.Equal(t, adminBot.Balance, 20000000000-chargeAmount)
	})

	t.Run("Trade - user to user", func(t *testing.T) {
		t.Skip("skipping this. pass.")
		var isDone bool

		cfg := config.NewConfig()
		err := cfg.Load("../../../common/config/config.yml")
		assert.Nil(t, err)

		// update admin current balance
		action := NewAction(cfg)
		adminBot := action.GetAdminBot()
		assert.NotNil(t, adminBot)
		adminBot.Balance = 19999970000

		// new wallet1
		newId1, err := uuid.NewUUID()
		assert.Nil(t, err)

		bot1 := bot.NewBot(
			newId1.String(),
			0,
			"test_token",
			func(bot *bot.Bot) {})

		isDone, err = action.NewWallet(bot1)
		assert.Nil(t, err)
		assert.True(t, isDone)
		isDone = false

		// new wallet2
		newId2, err := uuid.NewUUID()
		assert.Nil(t, err)

		bot2 := bot.NewBot(
			newId2.String(),
			// "8725977a-c1ec-11eb-b696-ce3c5a30693e",
			0,
			"test_token",
			func(bot *bot.Bot) {})

		isDone, err = action.NewWallet(bot2)
		assert.Nil(t, err)
		assert.True(t, isDone)
		isDone = false

		// issue to new wallet1
		isDone, err = action.Trade(
			adminBot,
			bot1,
		)
		assert.Nil(t, err)
		assert.True(t, isDone)
		isDone = false

		// trade wallet to wallet
		isDone, err = action.Trade(
			bot1,
			bot2,
		)
		assert.Nil(t, err)
		assert.True(t, isDone)

		// verify balance
		should, err := strconv.ParseFloat(cfg.RemittanceAmount, 64)
		fee, err := strconv.ParseFloat(cfg.RemittanceFee, 64)
		assert.Nil(t, err)
		assert.Equal(t, bot2.Balance, should-fee)
	})

	t.Run("handler", func(t *testing.T) {
		t.Skip("skipping this. pass.")
		cfg := config.NewConfig()
		err := cfg.Load("../../../common/config/config.yml")
		assert.Nil(t, err)

		// update admin current balance
		action := NewAction(cfg)
		adminBot := action.GetAdminBot()
		assert.NotNil(t, adminBot)
		adminBot.Balance = 209997270000

		// start handler
		action.Handler()
	})
}
