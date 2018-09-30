package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/tidyoux/chatbot/plugins"
)

func encryptMD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func encryptSHA1(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func encryptSHA256(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func crypto(command *bot.Cmd) (string, error) {
	if len(command.Args) < 2 {
		return plugins.InvalidAmountOfParams, nil
	}

	inputData := []byte(strings.Join(command.Args[1:], " "))
	switch strings.ToUpper(command.Args[0]) {
	case "MD5":
		return encryptMD5(inputData), nil
	case "SHA1", "SHA-1":
		return encryptSHA1(inputData), nil
	case "SHA256", "SHA-256":
		return encryptSHA256(inputData), nil
	default:
		return plugins.InvalidParams, nil
	}
}

func init() {
	bot.RegisterCommand(
		"crypto",
		"Encrypts the input data from its hash value",
		"md5|sha-1|sha-256 enter here text to encrypt",
		crypto)
}
