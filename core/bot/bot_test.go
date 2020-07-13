package bot

import (
	"testing"

	"blockchain.automation/core/queue"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	q := queue.NewMessageQueue()

	got := NewBot(
		"id",
		0,
		"token",
		func(b *Bot) {
			current := b.FSM.Current()
			assert.Equal(t, current, BotStateNew)
			q.Put(b)
		})

	if got.GetState() != BotStateClosed {
		t.Errorf("expected state to be %s", BotStateClosed)
	}
	err := got.SetState(BotStateNew)
	if err != nil {
		t.Error(err.Error())
	}

	i, err := q.Get()
	assert.Nil(t, err)
	assert.NotNil(t, i)

	g1 := i.(*Bot)
	assert.Equal(t, got, g1)
}
