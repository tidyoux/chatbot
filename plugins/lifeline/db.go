package lifeline

import "github.com/tidyoux/chatbot/db"

const (
	namespace = "lifeline"

	sectionKey = "section"
	statusKey  = "status"
)

var (
	dbInstance = db.New(namespace)
)

func getDBData(key string) (string, error) {
	sec, err := dbInstance.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(sec), nil
}

func setDBData(key string, value string) error {
	return dbInstance.Set([]byte(key), []byte(value))
}

func getSection(channel string) (string, error) {
	return getDBData(channel + sectionKey)
}

func setSection(channel string, section string) error {
	return setDBData(channel+sectionKey, section)
}

func getStatus(channel string, key string) (string, error) {
	return getDBData(channel + statusKey + key)
}

func setStatus(channel string, key string, value string) error {
	return setDBData(channel+statusKey+key, value)
}
