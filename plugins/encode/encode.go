package encoding

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/tidyoux/chatbot/plugins"
)

func encode(command *bot.Cmd) (string, error) {
	if len(command.Args) < 2 {
		return plugins.InvalidAmountOfParams, nil
	}

	var str string
	var err error
	switch command.Args[0] {
	case "base64":
		s := strings.Join(command.Args[1:], " ")
		str, err = encodeBase64(s)
	default:
		return plugins.InvalidParams, nil
	}

	if err != nil {
		return fmt.Sprintf("Error: %s", err), nil
	}

	return str, nil
}

func encodeBase64(str string) (string, error) {
	data := []byte(str)
	return base64.StdEncoding.EncodeToString(data), nil
}

func init() {
	bot.RegisterCommand(
		"encode",
		"Allows you encoding a value",
		"base64 enter here text to encode",
		encode)
}
