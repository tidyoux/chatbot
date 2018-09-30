package lisp

import (
	"strings"

	"github.com/glycerine/zygomys/zygo"
	"github.com/go-chat-bot/bot"
)

func lisp(command *bot.Cmd) (string, error) {
	code := strings.Join(command.Args, " ")
	env := zygo.NewZlispSandbox()
	err := env.LoadString(code)
	if err != nil {
		return "", err
	}

	expr, err := env.Run()
	if err != nil {
		return "", err
	}
	return expr.SexpString(nil), nil
}

func init() {
	bot.RegisterCommand(
		"lisp",
		"Runs a lisp code.",
		"'(+ 1 1)'",
		lisp)
}
