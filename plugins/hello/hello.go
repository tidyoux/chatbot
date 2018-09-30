package hello

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

func hello(command *bot.Cmd) (string, error) {
	msg := fmt.Sprintf("Hello %s!", command.User.RealName)
	return msg, nil
}

func init() {
	bot.RegisterCommand(
		"hello",
		"Sends a 'Hello' message to you on the channel.",
		"",
		hello)
}
