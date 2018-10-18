package rand

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/go-chat-bot/bot"
)

func getSpan(args []string) (from, to int) {
	from = 0
	to = 100
	if len(args) < 2 {
		return
	}

	l, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	r, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	if l == r {
		return
	}

	if l < r {
		return l, r
	}
	return r, l
}

func random(command *bot.Cmd) (msg string, err error) {
	from, to := getSpan(command.Args)
	return fmt.Sprint(from + rand.Intn(to-from)), nil
}

func init() {
	bot.RegisterCommand(
		"rand",
		"Returns a pseudo-random number in [a, b)",
		"0 100",
		random)
}
