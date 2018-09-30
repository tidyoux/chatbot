package lifeline

import (
	"log"

	"github.com/go-chat-bot/bot"
	"github.com/tidyoux/chatbot/plugins/lifeline/data"
)

const (
	startCmd = "start"
)

var (
	storyInstance *Story
)

func closeResultChan(result *bot.CmdResultV3) {
	close(result.Message)
	close(result.Done)
}

func lifeline(cmd *bot.Cmd) (bot.CmdResultV3, error) {
	result := bot.CmdResultV3{
		Message: make(chan string),
		Done:    make(chan bool),
	}

	var answer string
	if len(cmd.Args) > 0 {
		answer = cmd.Args[0]
	}

	go func() {
		defer closeResultChan(&result)

		if !getLock(cmd.Channel) {
			return
		}

		defer releaseLock(cmd.Channel)

		var currentSection string
		if answer == startCmd {
			answer = ""
			err := setSection(cmd.Channel, startSection)
			if err != nil {
				log.Println(namespace, err)
				return
			}
			currentSection = startSection
		} else {
			sec, err := getSection(cmd.Channel)
			if err != nil {
				log.Println(namespace, err)
				return
			}

			if len(sec) == 0 {
				err := setSection(cmd.Channel, startSection)
				if err != nil {
					log.Println(namespace, err)
					return
				}
				currentSection = startSection
			} else {
				currentSection = sec
			}
		}

		var err error
		ctx := &Context{
			channel:        cmd.Channel,
			currentSection: currentSection,
			msgch:          result.Message,
			data:           answer,
		}
		if len(answer) > 0 {
			err = storyInstance.Reply(ctx)
		} else {
			err = storyInstance.Play(ctx)
		}
		if err != nil {
			log.Println(namespace, err)
			return
		}

		setSection(cmd.Channel, ctx.currentSection)
	}()

	return result, nil
}

func init() {
	bot.RegisterCommandV3(
		"lifeline",
		"Play a game named lifeline.",
		"start",
		lifeline)

	storyInstance = newStory(data.Story)
	err := storyInstance.Init()
	if err != nil {
		log.Printf("Error: init lifeline story failed, %v\n", err)
	}
}
