package tracker

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/acdb"
)

var ch chan string

func StartUp() <-chan string {
	go readRegistrationDatabase()
	ch = make(chan string)
	return ch
}

func readRegistrationDatabase() {
	acdb.Setup(config.Conf)
	ch <- "db ready"
}
