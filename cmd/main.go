package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	token     string
	debug     bool
	proxyAddr string

	rootCmd = &cobra.Command{
		Use: "chatbot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return start()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "the telegram api token")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "set debug mode")
	rootCmd.PersistentFlags().StringVarP(&proxyAddr, "proxy", "p", "", "the socks5 proxy address, e.g.: 127.0.0.1:1080")
}

func main() {
	rootCmd.Execute()
}

func start() error {
	client := &http.Client{}
	if len(proxyAddr) > 0 {
		dialer, _ := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		client.Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}
	return run(token, debug, client)
}

func run(token string, debug bool, client *http.Client) error {
	apiClient, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return err
	}

	apiClient.Debug = debug

	log.Printf("Authorized on account %s", apiClient.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := apiClient.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	b := bot.New(&bot.Handlers{
		Response: func(target string, message string, sender *bot.User) {
			id, err := strconv.ParseInt(target, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}

			if len(message) == 0 {
				return
			}

			apiClient.Send(tgbotapi.NewChatAction(id, tgbotapi.ChatTyping))

			msg := tgbotapi.NewMessage(id, message)
			apiClient.Send(msg)
		},
	})
	b.Disable([]string{"cmd", "url"})

	for update := range updates {
		target := &bot.ChannelData{
			Protocol:  "telegram",
			Server:    "telegram",
			Channel:   strconv.FormatInt(update.Message.Chat.ID, 10),
			IsPrivate: update.Message.Chat.IsPrivate()}
		name := []string{update.Message.From.FirstName, update.Message.From.LastName}
		message := &bot.Message{
			Text: update.Message.Text,
		}

		go b.MessageReceived(target, message, &bot.User{
			ID:       strconv.Itoa(update.Message.From.ID),
			Nick:     update.Message.From.UserName,
			RealName: strings.Join(name, " ")})
	}
	return nil
}

func newOneTimeReplyKeyboard(labels [][]string) tgbotapi.ReplyKeyboardMarkup {
	var menu [][]tgbotapi.KeyboardButton
	for _, line := range labels {
		var row []tgbotapi.KeyboardButton
		for _, label := range line {
			row = append(row, tgbotapi.NewKeyboardButton(label))
		}
		menu = append(menu, row)
	}

	keyboard := tgbotapi.NewReplyKeyboard(menu...)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
