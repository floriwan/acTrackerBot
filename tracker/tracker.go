package tracker

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/acdb"
	"fmt"
	"log"
	"strings"
	"time"
)

var ch chan string
var stopChannels map[string]chan int

func StartUp() <-chan string {
	stopChannels = make(map[string]chan int)
	go readRegistrationDatabase()
	ch = make(chan string)
	return ch
}

func readRegistrationDatabase() {
	acdb.Setup(config.Conf)
	ch <- "db ready"
}

func GetRegListSize() int {
	return len(stopChannels)
}

func GetRegList() string {
	keys := make([]string, 0, len(stopChannels))
	for k := range stopChannels {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func RemoveReg(reg string) error {
	c, ok := stopChannels[reg]
	if !ok {
		return fmt.Errorf("no update process for reg '%v' found", reg)
	}

	// send something to stop the go routine
	log.Printf("stopping update process for '%v'\n", reg)
	c <- 0
	delete(stopChannels, reg)
	return nil
}

func AddNewReg(reg string) error {
	if !acdb.IsValidReg(reg) {
		return fmt.Errorf("registration '%v' is not a valid registration", reg)
	}
	go startTracker(config.Conf.UpdateIntervall, reg)
	return nil
}

func startTracker(interval int, reg string) {
	log.Printf("starting new tracker for '%v'\n", reg)

	sc := make(chan int)
	stopChannels[reg] = sc

	ticker := time.NewTicker(time.Duration(interval * int(time.Minute)))

	for {
		select {
		case <-ticker.C:
			log.Printf("update reg '%v'\n", reg)
		case <-ch:
			log.Printf("stop update for registration '%v'\n", reg)
			return
		}
	}

}
