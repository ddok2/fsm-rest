package bot

import (
	"strconv"

	cfg "blockchain.automation/common/config"
	"blockchain.automation/core/queue"
)

// Jobs make action job for test
type Jobs struct {
	BotList  []*Bot
	AdminBot *Bot
	TrashBot *Bot
	Config   *cfg.Config
	Queue    *queue.MessageQueue
	ch       chan *Bot
}

// NewJobs initialize queue & jobs return Jobs
func NewJobs(config *cfg.Config) *Jobs {
	j := &Jobs{}
	j.initialize(config)
	j.Queue = queue.NewMessageQueue()
	j.Queue.Initialize()
	return j
}

func (j *Jobs) initialize(config *cfg.Config) {
	j.Config = config

	j.BotList = make([]*Bot, j.Config.MemberCount)
	j.AdminBot = NewBot(config.AdminId, 0, "general_token_id", func(bot *Bot) {})
	j.TrashBot = NewBot(config.TrashId, 0, "general_token_id", func(bot *Bot) {})
	j.ch = make(chan *Bot)

	j.generateBots()
}

func (j *Jobs) generateBots() {
	botCount := j.Config.MemberCount
	botPrefixId := j.Config.MemberPrefix

	for i := 0; i < botCount; i++ {
		j.BotList[i] = NewBot(
			botPrefixId+strconv.Itoa(i),
			0,
			"general_token_id",
			func(bot *Bot) {
				j.Queue.Put(bot)
			},
		)
	}
}
