package db

import (
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/tidyoux/chatbot/db"
	"github.com/tidyoux/chatbot/plugins"
)

const (
	namespace = "db"

	success = "status:ok"
)

func dbop(command *bot.Cmd) (string, error) {
	if len(command.Args) < 2 {
		return plugins.InvalidAmountOfParams, nil
	}

	db := db.New(namespace)

	op := command.Args[0]
	key := []byte(command.Args[1])
	switch op {
	case "set":
		var data []byte
		if len(command.Args) > 2 {
			data = []byte(strings.Join(command.Args[2:], " "))
		}

		err := db.Set(key, data)
		if err != nil {
			return "", err
		}
		return success, nil
	case "get":
		v, err := db.Get(key)
		if err != nil {
			return "", err
		}
		return success + fmt.Sprintf("\ndata:%s", string(v)), nil
	default:
		return plugins.InvalidParams, nil
	}
}

func init() {
	bot.RegisterCommand(
		"db",
		"Stores key-value into db.",
		"set key value (or, get key)",
		dbop)
}
