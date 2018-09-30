package lifeline

import (
	"sync"
)

var (
	channelLocks sync.Map
)

func getLock(channel string) bool {
	if _, ok := channelLocks.Load(channel); ok {
		return false
	}
	channelLocks.Store(channel, struct{}{})
	return true
}

func releaseLock(channel string) {
	channelLocks.Delete(channel)
}
