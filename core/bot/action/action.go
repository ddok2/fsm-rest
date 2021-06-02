package action

import (
	"blockchain.automation/common/config"
	"blockchain.automation/core/bot"
)

// Action defines bots actions for test
type Action interface {
	NewWallet(bot *bot.Bot) (bool, error)
	Issue(bot *bot.Bot) (bool, error)
	Trade(origin *bot.Bot, target *bot.Bot) (bool, error)
	GetAdminBot() (bot *bot.Bot)
	Handler()
}

// NewAction return Action
// (bc -> NewBcBotAction, others -> DexBotAction)
func NewAction(cfg *config.Config) Action {
	switch cfg.OperationMode {
	case "blockchain":
		return NewBcBotAction(
			bot.NewJobs(
				cfg,
			),
		)
	default:
		return &DexBotAction{
			jobs: bot.NewJobs(cfg),
		}
	}
}
