package bot

import (
	"testing"

	"blockchain.automation/common/config"

	"github.com/stretchr/testify/assert"
)

func TestNewJobs(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.Load("../../common/config/config.yml")
	assert.Nil(t, err)

	jobs := NewJobs(cfg)
	assert.Equal(t, len(jobs.BotList), cfg.MemberCount)
}

func TestJobs(t *testing.T) {
	t.Run("generate bots", func(t *testing.T) {
		t.Skip("work fine")
		cfg := config.NewConfig()
		err := cfg.Load("../../common/config/config.yml")
		assert.Nil(t, err)

		jobs := NewJobs(cfg)
		jobs.generateBots()

		// check bot list
		assert.Equal(t, jobs.Config.MemberCount, len(jobs.BotList))

		// check queue bot
		for _, bot := range jobs.BotList {
			err = bot.SetState(BotStateNew)
			assert.Nil(t, err)
		}

		queueBot, err := jobs.Queue.Get()
		assert.Nil(t, err)
		assert.Equal(
			t,
			jobs.BotList[0],
			queueBot.(*Bot),
		)
	})
}
