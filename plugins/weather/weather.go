package weather

import (
	"strings"

	"github.com/go-chat-bot/bot"
)

func weather(command *bot.Cmd) (string, error) {
	location := "beijing"
	if len(command.Args) >= 1 && len(command.Args[0]) > 0 {
		location = strings.Join(command.Args, " ")
	}

	weatherData, err := getWeatherData(location)
	if err != nil {
		return "", err
	}

	return weatherData.format(), nil
}

func init() {
	bot.RegisterCommand(
		"weather",
		"Searchs weather information",
		"chaoyang,beijing",
		weather)
}
