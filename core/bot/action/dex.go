package action

import "blockchain.automation/core/bot"

// DexBotAction represents Action for Dex
type DexBotAction struct {
	jobs *bot.Jobs
}

// Handler handled bot 's actions
func (d *DexBotAction) Handler() {
}

// NewWallet returns an rest call of dex endpoint
func (d *DexBotAction) NewWallet(bot *bot.Bot) (bool, error) {
	return false, nil
}

// Issue returns an rest call of dex endpoint
func (d *DexBotAction) Issue(bot *bot.Bot) (bool, error) {
	return false, nil
}

// Trade returns an rest call of dex endpoint
func (d *DexBotAction) Trade(origin *bot.Bot, target *bot.Bot) (bool, error) {
	return false, nil
}

// GetAdminBot return admin bot from Jobs
func (b *DexBotAction) GetAdminBot() *bot.Bot {
	return b.jobs.AdminBot
}
