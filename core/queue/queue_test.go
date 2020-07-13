package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"blockchain.automation/core/bot"
)

func TestMessageQueue_Put(t *testing.T) {
	q := NewMessageQueue()

	b1 := bot.NewBot(
		"bot1",
		0,
		"token1",
		func(bot *bot.Bot) {},
	)

	q.Put(b1)
	assert.Equal(t, 1, q.queue.Len())

	got, err := q.Get()
	assert.Nil(t, err)
	assert.Equal(t, b1, got)

	got, err = q.Get()
	assert.Error(t, err, "empty queue")
	assert.Nil(t, got)
}

func TestMessageQueue_Get(t *testing.T) {
	q := NewMessageQueue()

	b1 := bot.NewBot("bot1", 0, "token1", func(bot *bot.Bot) {})
	b2 := bot.NewBot("bot2", 0, "token1", func(bot *bot.Bot) {})

	q.Put(b1)
	q.Put(b2)

	got1, err := q.Get()
	assert.Nil(t, err)
	assert.Equal(t, b1, got1)

	got2, err := q.Get()
	assert.Nil(t, err)
	assert.Equal(t, b2, got2)

	g1 := got1.(*bot.Bot)
	assert.Equal(t, g1.GetState(), "closed")

	err = g1.SetState(bot.BotStateNew)
	assert.Nil(t, err)
	assert.Equal(t, bot.BotStateNew, g1.GetState())
}
