package rand

import (
	"fmt"
	"math/rand"

	"github.com/go-chat-bot/bot"
)

func random(command *bot.Cmd) (msg string, err error) {
	return fmt.Sprint(rand.Intn(100)), nil
}

func init() {
	bot.RegisterCommand(
		"rand",
		"Returns a pseudo-random number in [0,100)",
		"",
		random)
}
