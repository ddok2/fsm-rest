package bot

import (
	"fmt"

	"github.com/looplab/fsm"
)

const (
	// BotStateClosed fsm/event, closed state
	BotStateClosed = "closed"

	// BotStateNew fsm/event, new state
	BotStateNew = "new"

	// BotStateIssue fsm/event, issue balance to admin state
	BotStateIssue = "charge"

	// BotStateTrade fsm/event, trade bot to bot state
	BotStateTrade = "trade"

	// BotStateEnd fsm/event, end state
	BotStateEnd = "end"
)

// Bot represents bot for test
type Bot struct {
	Id      string
	Balance float64
	TokenId string
	FSM     *fsm.FSM
}

type initFn func(bot *Bot)

// NewBot return an initialized Bot
func NewBot(
	id string,
	balance float64,
	tokenId string,
	initFn initFn,
) *Bot {
	bot := &Bot{
		Id:      id,
		Balance: balance,
		TokenId: tokenId,
	}

	bot.FSM = fsm.NewFSM(
		BotStateClosed,
		fsm.Events{
			{Name: BotStateNew, Src: []string{BotStateClosed}, Dst: BotStateNew},
			{Name: BotStateIssue, Src: []string{BotStateTrade}, Dst: BotStateIssue},
			{Name: BotStateTrade, Src: []string{BotStateTrade}, Dst: BotStateIssue},
			{Name: BotStateTrade, Src: []string{BotStateNew}, Dst: BotStateTrade},
			{Name: BotStateTrade, Src: []string{BotStateIssue}, Dst: BotStateTrade},
			{Name: BotStateTrade, Src: []string{BotStateTrade}, Dst: BotStateTrade},
			{Name: BotStateEnd, Src: []string{BotStateTrade}, Dst: BotStateEnd},
		},
		fsm.Callbacks{
			BotStateNew: func(e *fsm.Event) {
				initFn(bot)
			},
			"enter_state": func(e *fsm.Event) { bot.onEvent(e) },
		},
	)

	return bot
}

func (b *Bot) onEvent(e *fsm.Event) {
	fmt.Printf("Bot: enter state - %s, -id: %s", e.Dst, b.Id)
}

// GetState return Bot current state
func (b *Bot) GetState() string {
	return b.FSM.Current()
}

// SetState set Bot state
func (b *Bot) SetState(state string) error {
	return b.FSM.Event(state)
}
